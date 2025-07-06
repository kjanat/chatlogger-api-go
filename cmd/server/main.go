// Package main is the entry point for the ChatLogger API server.
// It initializes the configuration, sets up database connections,
// initializes all services and repositories, and starts the HTTP server.
package main

import (
	"log"

	"github.com/kjanat/chatlogger-api-go/internal/api"
	"github.com/kjanat/chatlogger-api-go/internal/config"
	"github.com/kjanat/chatlogger-api-go/internal/jobs"
	"github.com/kjanat/chatlogger-api-go/internal/repository"
	"github.com/kjanat/chatlogger-api-go/internal/service"
	"github.com/kjanat/chatlogger-api-go/internal/version"
)

//	@title			ChatLogger API (Go)
//	@description	API for logging and managing chat sessions.
//	@termsOfService	https://github.com/kjanat/chatlogger-api-go#terms-of-service

//	@contact.name	ChatLogger
//	@contact.url	https://github.com/kjanat/chatlogger-api-go/issues
//	@contact.email	chatlogger-api-go+swagger@kjanat.com

//	@license.name	MIT License
//	@license.url	https://github.com/kjanat/chatlogger-api-go/blob/master/LICENSE

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						x-organization-api-key

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

//	@externalDocs.description	GitHub Wiki
//	@externalDocs.url			https://github.com/kjanat/chatlogger-api-go/wiki
//	@externalDocs.name			Wiki

func main() {
	// Log version information at startup
	log.Printf("Starting ChatLogger API v%s (built: %s, commit: %s)",
		version.Version, version.BuildTime, version.GitCommit)

	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Database with migrations enabled (server is responsible for migrations)
	dbOptions := repository.DefaultDatabaseOptions()
	dbOptions.RunMigrations = true // Server should run migrations

	db, err := repository.NewDatabaseWithOptions(cfg.DatabaseURL, dbOptions)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// 3. Initialize Repositories
	orgRepo := repository.NewOrganizationRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	exportRepo := repository.NewExportRepository(db.DB)

	// 4. Initialize Job Queue
	queue := jobs.NewQueue(cfg.RedisAddr)
	defer func() {
		if err := queue.Close(); err != nil {
			log.Printf("Error closing queue connection: %v", err)
		}
	}()

	// 5. Initialize Services
	orgService := service.NewOrganizationService(orgRepo)
	apiKeyService := service.NewAPIKeyService(apiKeyRepo)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	chatService := service.NewChatService(chatRepo)
	messageService := service.NewMessageService(messageRepo)
	exportService := service.NewExportService(exportRepo, queue)
	swaggerService := service.NewSwaggerService()

	// Configure Swagger documentation with API information
	swaggerService.SetSwaggerInfo(
		version.Version,
		cfg.ApiServer.Scheme,
		cfg.ApiServer.Host,
		cfg.ApiServer.Port,
	)

	// 6. Bundle services for dependency injection
	services := &api.AppServices{
		OrganizationService: orgService,
		APIKeyService:       apiKeyService,
		UserService:         userService,
		ChatService:         chatService,
		MessageService:      messageService,
		ExportService:       exportService,
		SwaggerService:      swaggerService,
		Config: &api.AppConfig{
			ExportDir: cfg.ExportDir,
			APIServer: struct {
				Host   string
				Port   string
				Scheme string
			}{
				Host:   cfg.ApiServer.Host,
				Port:   cfg.ApiServer.Port,
				Scheme: cfg.ApiServer.Scheme,
			},
		},
	}

	// 7. Set up Gin Router with routes and inject services
	router := api.NewRouter(services, cfg.JWTSecret)

	// 8. Start the Server
	port := cfg.ServerPort
	log.Printf("Server listening on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Printf("Failed to run server: %v", err)
		return
	}
}
