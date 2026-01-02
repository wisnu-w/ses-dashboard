// Package main SES Monitoring API
//
//	@title			SES Monitoring API
//	@version		1.0
//	@description	API for monitoring AWS SES events
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"

	_ "ses-monitoring/docs"
	"ses-monitoring/internal/config"
	"ses-monitoring/internal/delivery/http"
	"ses-monitoring/internal/infrastructure/database"
	"ses-monitoring/internal/infrastructure/repository"
	"ses-monitoring/internal/services"
	"ses-monitoring/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db := database.NewPostgres(dsn)

	sesRepo := repository.NewSESEventRepository(db)
	userRepo := repository.NewUserRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)
	suppressionRepo := repository.NewSuppressionRepository(db)
	suppressionDBRepo := database.NewSuppressionRepository(db)

	// Initialize AWS client and sync service
	// Initialize services
	syncService := services.NewSyncService(
		settingsRepo,
		suppressionDBRepo,
	)
	cleanupService := services.NewCleanupService(settingsRepo, sesRepo)

	// Start background services
	go syncService.StartBackgroundSync(context.Background())
	go cleanupService.StartCleanupScheduler(context.Background())

	sesUC := usecase.NewSESUsecase(sesRepo)
	authUC := usecase.NewAuthUsecase(userRepo, cfg.App.JWTSecret)

	snsHandler := http.NewSNSHandler(sesUC, cfg)
	monitoringHandler := http.NewMonitoringHandler(sesUC, settingsRepo)
	authHandler := http.NewAuthHandler(authUC)
	userHandler := http.NewUserHandler(authUC)
	settingsHandler := http.NewSettingsHandler(settingsRepo)
	suppressionHandler := http.NewSuppressionHandler(settingsRepo, suppressionRepo, suppressionDBRepo, syncService)
	healthHandler := http.NewHealthHandler()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// ========================
	// CORS CONFIG (FIX NGROK + SWAGGER)
	// ========================
	corsConfig := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // allow all origins (swagger + ngrok)
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	// ========================
	// PUBLIC ROUTES
	// ========================
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)
	r.POST("/sns/ses", snsHandler.Handle)
	r.POST("/api/login", authHandler.Login)

	// ========================
	// SWAGGER
	// ========================
	if cfg.App.EnableSwagger {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
	}

	// ========================
	// PROTECTED ROUTES
	// ========================
	api := r.Group("/api")
	api.Use(http.JWTAuthMiddleware([]byte(cfg.App.JWTSecret)))
	api.Use(func(c *gin.Context) {
		c.Set("monitoring_handler", monitoringHandler)
		c.Next()
	})
	{
		api.GET("/events", monitoringHandler.GetEvents)
		api.GET("/metrics", monitoringHandler.GetMetrics)
		api.GET("/metrics/daily", monitoringHandler.GetDailyMetrics)
		api.GET("/metrics/monthly", monitoringHandler.GetMonthlyMetrics)
		api.GET("/metrics/hourly", monitoringHandler.GetHourlyMetrics)

		// User management routes (admin only)
		admin := api.Group("")
		admin.Use(http.AdminMiddleware())
		{
			admin.POST("/users", userHandler.CreateUser)
			admin.GET("/users", userHandler.GetUsers)
			admin.PUT("/users/:id/reset-password", userHandler.ResetPassword)
			admin.PUT("/users/:id/disable", userHandler.DisableUser)
			admin.PUT("/users/:id/enable", userHandler.EnableUser)
			admin.DELETE("/users/:id", userHandler.DeleteUser)

			// Settings routes (admin only)
			admin.GET("/settings/aws", settingsHandler.GetAWSSettings)
			admin.PUT("/settings/aws", settingsHandler.UpdateAWSSettings)
			admin.POST("/settings/aws/test", settingsHandler.TestAWSConnection)
			admin.GET("/settings/retention", settingsHandler.GetRetentionSettings)
			admin.PUT("/settings/retention", settingsHandler.UpdateRetentionSettings)
			admin.GET("/settings/timezone", settingsHandler.GetTimezoneSettings)
			admin.PUT("/settings/timezone", settingsHandler.UpdateTimezoneSettings)

			// AWS SES Suppression management routes (admin only)
			admin.GET("/suppression", suppressionHandler.GetSuppressions)
			admin.POST("/suppression", suppressionHandler.AddSuppression)
			admin.POST("/suppression/bulk", suppressionHandler.BulkAddSuppression)
			admin.DELETE("/suppression/bulk", suppressionHandler.BulkRemoveSuppression)
			admin.POST("/suppression/sync", suppressionHandler.SyncFromAWS)
			admin.GET("/suppression/sync/status", suppressionHandler.GetSyncStatus)
			admin.DELETE("/suppression/:email", suppressionHandler.RemoveSuppression)
			admin.GET("/suppression/:email/status", settingsHandler.CheckEmailSuppression)
		}

		// User routes (authenticated users)
		api.PUT("/change-password", userHandler.ChangePassword)
	}

	r.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
