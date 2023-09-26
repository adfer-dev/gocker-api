package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"

	"github.com/joho/godotenv"
)

type UserRole int

const (
	Admin UserRole = iota + 1
	Standard
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  []byte `json:"password" validate:"required"`
	Role      UserRole
}

// Function that encodes user's password using AES encryption.
func (user *User) EncodePassword(password string) error {
	key, envErr := getPasswordKey()

	if envErr != nil {
		return envErr
	}

	cipherBlock, cipherErr := aes.NewCipher([]byte(key))

	if cipherErr != nil {
		return cipherErr
	}

	gcm, gcmErr := cipher.NewGCM(cipherBlock)

	if gcmErr != nil {
		return gcmErr
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, readErr := io.ReadFull(rand.Reader, nonce); readErr != nil {
		return readErr
	}

	user.Password = gcm.Seal(nonce, nonce, []byte(password), nil)

	return nil
}

// Function that decodes user's password using AES decryption and compares it to the input password.
func (user User) ComparePassword(password string) error {
	key, envErr := getPasswordKey()

	if envErr != nil {
		return envErr
	}

	cipherBlock, cipherErr := aes.NewCipher([]byte(key))

	if cipherErr != nil {
		return cipherErr
	}

	gcm, gcmErr := cipher.NewGCM(cipherBlock)
	nonceSize := gcm.NonceSize()

	if gcmErr != nil || (len(user.Password) < nonceSize) {
		return gcmErr
	}

	nonce, cipherText := user.Password[:nonceSize], user.Password[nonceSize:]

	passwordText, decryptErr := gcm.Open(nil, []byte(nonce), []byte(cipherText), nil)

	if decryptErr != nil {
		return decryptErr
	}

	if string(passwordText) != password {
		return errors.New("wrong password. Please, try again")
	}

	return nil
}

func getPasswordKey() (string, error) {
	envErr := godotenv.Load()

	return os.Getenv("USER_PASSWORD_KEY"), envErr
}
