package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"user-simple-crud/config"
	"user-simple-crud/internal/delivery/http"
	api "user-simple-crud/internal/delivery/http/middleware"
	"user-simple-crud/internal/delivery/http/route"
	"user-simple-crud/internal/repository"
	services "user-simple-crud/internal/services"
	"user-simple-crud/migration"
	"user-simple-crud/pkg/database"
	"user-simple-crud/pkg/httpclient"
	"user-simple-crud/pkg/logger"
	"user-simple-crud/pkg/server"
	"user-simple-crud/pkg/signature"
	"user-simple-crud/pkg/xvalidator"
)

var (
	sqlClientRepo *database.Database
)

// @title           user-simple-crud
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9004
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	validate, _ := xvalidator.NewValidator()
	conf := config.InitAppConfig(validate)
	logger.SetupLogger(&logger.Config{
		AppENV:  conf.AppEnvConfig.AppEnv,
		LogPath: conf.AppEnvConfig.LogFilePath,
		Debug:   conf.AppEnvConfig.AppDebug,
	})
	initInfrastructure(conf)
	ginServer := server.NewGinServer(&server.GinConfig{
		HttpPort:     conf.AppEnvConfig.HttpPort,
		AllowOrigins: conf.AppEnvConfig.AllowOrigins,
		AllowMethods: conf.AppEnvConfig.AllowMethods,
		AllowHeaders: conf.AppEnvConfig.AllowHeaders,
	})
	// external
	signaturer := signature.NewSignature(conf.AuthConfig.JwtSecretAccessToken)
	// repository
	userRepository := repository.NewUserSQLRepository()

	// service
	userService := services.NewUserService(sqlClientRepo.GetDB(), userRepository, signaturer, validate)
	// Handler
	authMiddleware := api.NewAuthMiddleware(signaturer)
	userHandler := http.NewUserHTTPHandler(userService)

	router := route.Router{
		App:            ginServer.App,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
	router.Setup()
	router.SwaggerRouter()
	echan := make(chan error)
	go func() {
		echan <- ginServer.Start()
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		slog.Info("signal terminated detected")
	case err := <-echan:
		slog.Error("Failed to start http server", err)
	}
}

func initInfrastructure(config *config.Config) {
	//initPostgreSQL()
	sqlClientRepo = initSQL(config)
}

func initSQL(conf *config.Config) *database.Database {
	db := database.NewDatabase(conf.DatabaseConfig.Dbservice, &database.Config{
		DbHost:   conf.DatabaseConfig.Dbhost,
		DbUser:   conf.DatabaseConfig.Dbuser,
		DbPass:   conf.DatabaseConfig.Dbpassword,
		DbName:   conf.DatabaseConfig.Dbname,
		DbPort:   strconv.Itoa(conf.DatabaseConfig.Dbport),
		DbPrefix: conf.DatabaseConfig.DbPrefix,
	})
	if conf.IsStaging() {
		migration.AutoMigration(db)
	}
	return db
}

func initHttpclient() httpclient.Client {
	httpClientFactory := httpclient.New()
	httpClient := httpClientFactory.CreateClient()
	return httpClient
}
