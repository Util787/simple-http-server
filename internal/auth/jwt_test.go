package auth

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secretkey"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if token == "" {
		t.Errorf("Expected token to be non-empty")
	}

	parsedToken, _ := jwt.Parse(token, nil)
	fmt.Println("MakeJWT Parsed token: ", parsedToken)
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secretkey"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Failed to generate token: %v", err)
	}

	parsedUUID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if parsedUUID != userID {
		t.Errorf("Expected parsed UUID to be %v, but got %v", userID, parsedUUID)
	}

	invalidToken := "invalid-token"
	_, err = ValidateJWT(invalidToken, tokenSecret)
	if err == nil {
		t.Error("Expected error for invalid token, but got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{"Authorization": {"Bearer TOKEN_STRING"}}

	tokenstring, err := GetBearerToken(header)
	log.Println(tokenstring)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	wrongheader:= http.Header{"A": {"Bearer TOKEN_STRING"}}
	wrongtokenstring, err := GetBearerToken(wrongheader)
	log.Println(wrongtokenstring)
	if err == nil {
		t.Errorf("Expected error for empty Authorization header, but got none")
	}
}
