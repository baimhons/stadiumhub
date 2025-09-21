package initial

import (
	"fmt"
	"log"

	"github.com/baimhons/stadiumhub.git/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type clientConfig struct {
	DB *gorm.DB
}

func newClientConfig() *clientConfig {
	db := connectPostgresDatabase(
		internal.ENV.Database.Host,
		internal.ENV.Database.Port,
		internal.ENV.Database.User,
		internal.ENV.Database.Password,
		internal.ENV.Database.Name,
	)

	return &clientConfig{
		DB: db,
	}
}

func connectPostgresDatabase(
	host string,
	port int,
	user string,
	password string,
	database string,
) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok",
		host,
		user,
		password,
		database,
		port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("database connected successfully")

	return db
}
