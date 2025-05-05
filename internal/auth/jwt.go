package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claims := jwt.RegisteredClaims{}
	claims.Issuer = "chirpy"
	claims.IssuedAt = &jwt.NumericDate{Time: time.Now()}
	claims.ExpiresAt = &jwt.NumericDate{Time: time.Now().Add(expiresIn)}
	claims.Subject = userID.String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedjwt, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedjwt, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unidentified sign method")
		}
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, keyfunc)
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	uuidString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	TokenString := headers.Get("Authorization")
	if TokenString == "" {
		return "", errors.New("empty Authorization header") // failed to get token string:
	}
	return TokenString[7:], nil
}
