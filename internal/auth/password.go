package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password []byte) ([]byte, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	return hashed_password, nil
}

func CheckPasswordHash(password, hashed_password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashed_password, password)
	if err != nil {
		return err
	}
	return nil
}
