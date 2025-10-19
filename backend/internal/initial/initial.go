package initial

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/wisaitas/share-pkg/utils"
)

func init() {
	if err := utils.ReadConfig(&internal.ENV); err != nil {
		log.Fatalf("error reading config: %v\n", utils.Error(err))
	}
}

type App struct {
	App          *gin.Engine
	ClientConfig *clientConfig
}

func InitializeApp() *App {
	clientConfig := newClientConfig()

	ginEngine := gin.Default()
	SetupMiddleware(ginEngine)

	repository := NewRepository(clientConfig)
	service := NewService(repository, clientConfig.Redis)
	handler := NewHandler(service)
	validate := NewValidate()

	middleware := NewMiddleware(clientConfig.Redis)

	NewRoute(
		ginEngine,
		*handler,
		*validate,
		*middleware,
	)

	go worker.NewMatchWorker(service.MatchService).Start()

	go worker.NewBookingWorker(service.BookingService).Start()

	return &App{
		App:          ginEngine,
		ClientConfig: clientConfig,
	}
}

func (a *App) Run() chan os.Signal {
	go func() {
		if err := a.App.Run(fmt.Sprintf(":%d", internal.ENV.Server.Port)); err != nil {
			log.Fatalf("error starting server: %v\n", utils.Error(err))
		}
	}()

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-gracefulShutdown

	return gracefulShutdown
}
