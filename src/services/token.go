package services

import (
	"errors"
	"gocker-api/database"
	"gocker-api/models"
)

func GetTokenByUserRefer(userRefer uint) (token models.Token, err error) {
	database := database.GetInstance().GetDB()
	result := database.Find(&token, "user_refer = ?", userRefer)

	if result.RowsAffected == 0 {
		err = errors.New("user not found")
	}

	return
}

func CreateToken(token *models.Token) {
	database := database.GetInstance().GetDB()
	database.Create(token)
}
