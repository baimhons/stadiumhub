package initial

import (
	"fmt"
	"log"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/baimhons/stadiumhub/internal/zone"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type clientConfig struct {
	DB    *gorm.DB
	Redis utils.RedisClient
}

func newClientConfig() *clientConfig {
	db := ConnectMySQLDatabase(
		internal.ENV.Database.Host,
		internal.ENV.Database.Port,
		internal.ENV.Database.User,
		internal.ENV.Database.Password,
		internal.ENV.Database.Name,
	)

	redis := utils.ConnectRedis(
		internal.ENV.Redis.Host,
		internal.ENV.Redis.Port,
		internal.ENV.Redis.Password,
	)

	redisWrapper := utils.NewRedisClient(redis)

	return &clientConfig{
		DB:    db,
		Redis: redisWrapper,
	}
}

func ConnectMySQLDatabase(
	host string,
	port int,
	username string,
	password string,
	database string,
) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	errAutoMigrate := db.AutoMigrate(
		&user.User{},
		&team.Team{},
		&match.Match{},
		&zone.Zone{},
		&seat.Seat{},
		&booking.Booking{},
		&booking.BookingSeat{},
	)

	// seed.SeedTeam(db)
	// match.SeedMatches(db)
	// seed.SeedZones(db)
	// seed.SeedSeats(db)

	if errAutoMigrate != nil {
		log.Fatalf("failed to auto migrate database: %v", errAutoMigrate)
	}

	log.Println("MySQL database connected successfully")

	return db
}
