package initial

import (
	"log"

	"github.com/baimhons/stadiumhub.git/internal"
	"github.com/wisaitas/share-pkg/utils"
)

func init() {
	if err := utils.ReadConfig(&internal.ENV); err != nil {
		log.Fatalf("error reading config: %v\n", utils.Error(err))
	}
}

// func InitializeApp() *App {
// 	clientConfig := newClientConfig()

// 	ginEngine := gin.Default()

// 	setupMiddleware(ginEngine)

// 	sharePkg := newSharePkg(clientConfig)

// 	repository := newRepository(clientConfig)
// 	service := newService(repository, sharePkg)
// 	handler := newHandler(service)
// 	validate := newValidate(sharePkg)
// 	middleware := newMiddleware(sharePkg)

// 	newRoute(app, handler, validate, middleware)

// 	return &App{
// 		App:          app,
// 		ClientConfig: clientConfig,
// 	}
// }
