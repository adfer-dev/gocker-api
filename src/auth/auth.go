package auth

import (
	"errors"
	"gocker-api/models"
	"gocker-api/services"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// Returns a new token as string and an error (if there was one)
func GenerateToken(user models.User, kind models.TokenKind) (string, error) {
	secretKey, envErr := getSecretKey()

	if envErr != nil {
		return "", nil
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	var expiration int64

	if kind == models.Bearer {
		expiration = time.Now().Add(24 * time.Hour).Unix()
	} else {
		expiration = time.Now().Add(8766 * time.Hour).Unix()
	}

	claims["exp"] = expiration
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

	token, parseErr := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("signing method not valid")
		}
		return secretKey, nil
	})

	if !token.Valid {
		return parseErr
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

// Function that performs JWT authentication check
func CheckAuth(res http.ResponseWriter, req *http.Request) error {
	fullToken := req.Header.Get("Authorization")

	if fullToken == "" || !strings.HasPrefix(fullToken, "Bearer") {
		return errors.New("authorization header must be provided, starting with Bearer")
	}

	tokenString := fullToken[7:]

	//Validate token
	if err := ValidateToken(tokenString); err != nil {
		//If token is not valid beacause it is expired
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return errors.New("token expired. Please, get a new one on /api/v1/auth/refresh-token")
		} else {
			return errors.New("token is not valid")
		}
	}

	claims, claimsErr := GetClaims(tokenString)

	if claimsErr != nil {
		return claimsErr
	}

	user, _ := services.GetUserByEmail(claims["email"].(string))

	if (req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE") && user.Role != 1 {
		return errors.New("method not allowed")
	}

	return nil
}
