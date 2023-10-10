package storage

import (
	"errors"
	"gocker-api/database"
	"gocker-api/models"
)

type UserStorage struct{}

const userTypeMismatchErr = "must be type user"

func (userStorage *UserStorage) Get(id int) (interface{}, error) {
	var user *models.User
	database := database.GetInstance().GetDB()
	if result := database.First(&user, "id = ?", id); result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (userStorage *UserStorage) Create(item interface{}) error {
	user, ok := item.(*models.User)

	if !ok {
		return errors.New(userTypeMismatchErr)
	}

	database := database.GetInstance().GetDB()
	database.Create(user)

	return nil
}

func (userStorage *UserStorage) Update(item interface{}) error {
	user, ok := item.(*models.User)

	if !ok {
		return errors.New(userTypeMismatchErr)
	}

	database := database.GetInstance().GetDB()
	database.Save(user)
	return nil
}

func (userStorage *UserStorage) Delete(item interface{}) error {
	user, ok := item.(*models.User)

	if !ok {
		return errors.New(userTypeMismatchErr)
	}

	database := database.GetInstance().GetDB()
	database.Delete(user)

	return nil
}
