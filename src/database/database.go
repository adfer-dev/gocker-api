package database

import (
	"gocker-api/models"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

func (database Database) GetDB() *gorm.DB {
	return database.db
}

var databaseInstance *Database
var lock = &sync.Mutex{}

func GetInstance() *Database {
	lock.Lock()
	defer lock.Unlock()

	if databaseInstance == nil {
		db, dbErr := gorm.Open(postgres.Open(os.Getenv("DB_STRING")), &gorm.Config{})

		if dbErr != nil {
			log.Fatal(dbErr.Error())
		}

		db.Logger = logger.Default.LogMode(logger.Info)
		db.AutoMigrate(&models.User{}, &models.Token{})
		databaseInstance = &Database{db}
	}

	return databaseInstance
}
