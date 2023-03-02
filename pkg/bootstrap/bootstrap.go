package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/MartinZitterkopf/gocurse_domain/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func DBConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	instanceDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		instanceDB = instanceDB.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := instanceDB.AutoMigrate(&domain.Curse{}); err != nil {
			return nil, err
		}
	}

	return instanceDB, nil
}
