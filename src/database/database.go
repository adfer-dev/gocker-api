package database

import (
	"gocker-api/models"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		envErr := godotenv.Load()

		if envErr != nil {
			log.Fatal(envErr.Error())
		}

		db, dbErr := gorm.Open(postgres.Open(os.Getenv("DB_STRING")), &gorm.Config{})

		if dbErr != nil {
			log.Fatal(dbErr.Error())
		}

		db.AutoMigrate(&models.User{}, &models.Token{})
		databaseInstance = &Database{db}
	}

	return databaseInstance
}
