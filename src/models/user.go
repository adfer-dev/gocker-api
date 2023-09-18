package models

import (
	"golang.org/x/crypto/bcrypt"
)

type UserRole int

const (
	adminRole UserRole = iota
	userRole
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Role      UserRole
}

func (user User) GetRole() string {
	switch user.Role {
	case adminRole:
		return "ADMIN"
	case userRole:
		return "USER"
	default:
		return "USER"
	}
}

func (user *User) SetRole(role string) {
	switch role {
	case "ADMIN":
		user.Role = adminRole
	case "USER":
		user.Role = userRole
	}
}

func (user *User) EncodePassword(password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	user.Password = string(hashedPass)

	return err
}

func (user User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	return err
}
