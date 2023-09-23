package models

type TokenKind int

const (
	Bearer TokenKind = iota + 1
	Refresh
)

type Token struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	TokenValue string `json:"token" validate:"required"`
	UserRefer  uint   `json:"user_id" validate:"required"`
	User       User   `gorm:"foreignKey:UserRefer"`
	Kind       TokenKind
}
