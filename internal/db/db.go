package db

import (
	"log"

	loggerModels "github.com/caseapia/goproject-flush/internal/models/logger"
	user "github.com/caseapia/goproject-flush/internal/models/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "root:root@tcp(127.0.0.1:3306)/flushproject?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	DB = db

	if err := db.AutoMigrate(&user.User{}, &loggerModels.ActionLog{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
}
