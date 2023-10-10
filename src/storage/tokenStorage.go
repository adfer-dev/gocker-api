package storage

import (
	"errors"
	"gocker-api/database"
	"gocker-api/models"
)

const tokenTypeMismatchErr = "type must be token"

type TokenStorage struct{}

func (tokenStorage *TokenStorage) Get(id int) (interface{}, error) {
	var token *models.Token
	database := database.GetInstance().GetDB()

	if result := database.Find(&token, "id = ?", id); result.RowsAffected == 0 {
		return nil, errors.New("token not found")
	}

	return token, nil
}

func (tokenStorage *TokenStorage) Create(item interface{}) error {
	token, ok := item.(*models.Token)

	if !ok {
		return errors.New(tokenTypeMismatchErr)
	}

	database := database.GetInstance().GetDB()
	database.Create(&token)
	return nil
}

func (tokenStorage *TokenStorage) Update(item interface{}) error {
	token, ok := item.(*models.Token)

	if !ok {
		return errors.New(tokenTypeMismatchErr)
	}

	database := database.GetInstance().GetDB()
	database.Save(&token)
	return nil
}

func (tokenStorage *TokenStorage) Delete(item interface{}) error {
	token, ok := item.(*models.Token)

	if !ok {
		return errors.New(tokenTypeMismatchErr)
	}

	database := database.GetInstance().GetDB()
	database.Delete(token)
	return nil
}
