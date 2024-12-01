package helpers

import (
	"fmt"
	"log"

	"github.com/anandtiwari11/library-management/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnectDB() {
	DB, err = gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to Connect to the Database", err)
	}
	err = DB.AutoMigrate(&models.Author{}, &models.Book{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to Migrate into the table", err)
	}
	fmt.Println("Successfully connected to SQLite")
}