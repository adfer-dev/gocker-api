package models

type TokenKind int

const (
	Access TokenKind = iota + 1
	Refresh
)

type Token struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	TokenValue string `json:"token" validate:"required"`
	UserRefer  uint   `json:"user_id" validate:"required"`
	Kind       TokenKind
}
