package services

import (
	"errors"
	"gocker-api/database"
	"gocker-api/models"
)

func GetTokensByUserReferAndKind(userRefer uint, kind models.TokenKind) (token models.Token, err error) {
	database := database.GetInstance().GetDB()
	result := database.Find(&token, "user_refer = ? AND kind = ?", userRefer, kind)

	if result.RowsAffected == 0 {
		err = errors.New("user not found")
	}

	return
}

func TokenExists(tokenString string) bool {
	var token models.Token
	database := database.GetInstance().GetDB()

	result := database.Find(&token, "token_value LIKE ?", tokenString)

	return result.RowsAffected > 0
}

func CreateToken(token *models.Token) {
	database := database.GetInstance().GetDB()
	database.Create(token)
}

func UpdateTokenValue(newTokenValue string) (token models.Token) {
	database := database.GetInstance().GetDB()

	token.TokenValue = newTokenValue
	database.Save(&token)

	return
}

func DeleteToken(token models.Token) {
	database := database.GetInstance().GetDB()

	database.Delete(token)
}
