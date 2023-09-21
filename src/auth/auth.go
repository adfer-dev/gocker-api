package auth

import (
	"errors"
	"gocker-api/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// Returns a new token as string and an error (if there was one)
func GenerateToken(user models.User) (string, error) {
	secretKey, envErr := getSecretKey()

	if envErr != nil {
		return "", nil
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["email"] = user.Email

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validates the token string passed and returns an error if it's not valid
func ValidateToken(tokenString string) error {

	secretKey, envErr := getSecretKey()

	if envErr != nil {
		return envErr
	}

	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("signing method not valid")
		}
		return secretKey, nil
	})

	if !token.Valid {
		return errors.New("token invalid or expired")
	}

	return nil
}

func GetClaims(tokenString string) (jwt.MapClaims, error) {
	secretKey, envErr := getSecretKey()

	if envErr != nil {
		return nil, envErr
	}

	jwtToken, parseErr := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if parseErr != nil {
		return nil, parseErr
	}

	claims := jwtToken.Claims.(jwt.MapClaims)

	return claims, nil
}

func getSecretKey() ([]byte, error) {
	envErr := godotenv.Load()

	if envErr != nil {
		return []byte{}, envErr
	}

	return []byte(os.Getenv("SECRET_KEY")), nil
}
