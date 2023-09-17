package models

import (
	"crypto/sha256"
	"encoding/hex"
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

func (user *User) EncodePassword(password string) {
	hash := sha256.Sum256([]byte(password))
	user.Password = hex.EncodeToString(hash[:])
}
