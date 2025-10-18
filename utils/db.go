package utils

import (
	"log"
	"soccer-team-api/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(env *Env) error {

	db, err := gorm.Open(postgres.Open(env.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	DB = db
	log.Println("✅ Database connected")
	return nil
}

func MigrateDB() error {
	err := DB.AutoMigrate(
		&model.Match{},
		&model.Player{},
		&model.Team{},
		&model.Goal{},
		&model.User{},
	)
	if err != nil {
		return err
	}
	log.Println("✅ Database migrated")
	return nil
}
func GetDB() *gorm.DB {
	return DB
}
