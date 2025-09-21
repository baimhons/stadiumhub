package initial

import (
	"github.com/baimhons/stadiumhub.git/internal"
	"github.com/baimhons/stadiumhub.git/internal/utils"
	"gorm.io/gorm"
)

type clientConfig struct {
	DB *gorm.DB
}

func newClientConfig() *clientConfig {
	db := utils.ConnectPostgresDatabase(
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
