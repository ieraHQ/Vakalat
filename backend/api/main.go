package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/ieraHQ/Vakalat/backend/api/ai"
	"github.com/ieraHQ/Vakalat/backend/api/auth"
	"github.com/ieraHQ/Vakalat/backend/api/config"
	"github.com/ieraHQ/Vakalat/backend/api/database"
	"github.com/ieraHQ/Vakalat/backend/api/handlers"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/middleware"
	"github.com/ieraHQ/Vakalat/backend/api/ocr"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"github.com/ieraHQ/Vakalat/backend/api/storage"
	"github.com/ieraHQ/Vakalat/backend/api/worker"
	"github.com/ieraHQ/Vakalat/backend/api/websocket"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		logger.GetLogger().Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize logger
	logger.InitLogger()
	defer logger.GetLogger().Sync()

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		logger.GetLogger().Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.CloseDB()

	// Initialize storage
	storage := storage.NewLocalStorage("./storage")

	// Initialize OCR service
	ocrService := ocr.NewOCRService(true, true)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	clientRepo := repositories.NewClientRepository(database.DB)
	matterRepo := repositories.NewMatterRepository(database.DB)
	documentRepo := repositories.NewDocumentRepository(database.DB)
	roleRepo := repositories.NewRoleRepository(database.DB)
	hearingRepo := repositories.NewHearingRepository(database.DB)
	orderRepo := repositories.NewOrderRepository(database.DB)
	searchRepo := repositories.NewSearchRepository(database.DB)
	aiRepo := repositories.NewAIRepository(database.DB)
	refreshTokenRepo := auth.NewRefreshTokenRepository(database.DB)
	passwordResetRepo := repositories.NewPasswordResetRepository(database.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	searchService := services.NewSearchService(searchRepo)
	llmClient := ai.NewLLMClient(cfg.AI.BaseURL, cfg.AI.Model, cfg.AI.EmbeddingModel)
	clientService := services.NewClientService(clientRepo, searchService, llmClient)
	matterService := services.NewMatterService(matterRepo, searchService, llmClient)
	documentService := services.NewDocumentService(documentRepo, storage)
	hearingService := services.NewHearingService(hearingRepo)
	orderService := services.NewOrderService(orderRepo)
	timelineService := services.NewTimelineService(hearingRepo, orderRepo)
	aiService := services.NewAIService(aiRepo, llmClient, cfg.AI.Model)
	permissionService := auth.NewPermissionService(userRepo, roleRepo)
	sessionService := auth.NewSessionService(userRepo, passwordResetRepo)

	// Initialize OCR worker
	ocrWorker := worker.NewOCRWorker(documentRepo, documentService, ocrService, searchService, llmClient, 30*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go ocrWorker.Start(ctx)
	defer cancel()

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize Fiber app
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware())

	// API routes
	api := app.Group("/api")

	// Auth routes
	authGroup := api.Group("/auth")
	authGroup.Post("/login", handlers.LoginHandler(cfg, userService, refreshTokenRepo))
	authGroup.Post("/refresh", handlers.RefreshTokenHandler(cfg, refreshTokenRepo, userService))
	authGroup.Post("/forgot-password", handlers.ForgotPasswordHandler(sessionService))
	authGroup.Post("/reset-password", handlers.ResetPasswordHandler(sessionService))
	authGroup.Post("/mfa/enable", middleware.AuthMiddleware(cfg, userService), handlers.EnableMFAHandler(sessionService))
	authGroup.Post("/mfa/verify", middleware.AuthMiddleware(cfg, userService), handlers.VerifyMFAHandler(sessionService))

	// User routes
	users := api.Group("/users", middleware.AuthMiddleware(cfg, userService))
	users.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_users"), handlers.GetUserHandler(userService))
	users.Post("/", middleware.RBACMiddleware(permissionService, "manage_users"), handlers.CreateUserHandler(userService))
	users.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_users"), handlers.UpdateUserHandler(userService))
	users.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_users"), handlers.DeleteUserHandler(userService))

	// Client routes
	clients := api.Group("/clients", middleware.AuthMiddleware(cfg, userService))
	clients.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_clients"), handlers.GetClientHandler(clientService))
	clients.Post("/", middleware.RBACMiddleware(permissionService, "manage_clients"), handlers.CreateClientHandler(clientService))
	clients.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_clients"), handlers.UpdateClientHandler(clientService))
	clients.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_clients"), handlers.DeleteClientHandler(clientService))
	clients.Get("/", middleware.RBACMiddleware(permissionService, "manage_clients"), handlers.ListClientsHandler(clientService))

	// Matter routes
	matters := api.Group("/matters", middleware.AuthMiddleware(cfg, userService))
	matters.Get("/", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.ListMattersHandler(matterService))
	matters.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.GetMatterHandler(matterService))
	matters.Post("/", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.CreateMatterHandler(matterService))
	matters.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.UpdateMatterHandler(matterService))
	matters.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.DeleteMatterHandler(matterService))
	matters.Get("/client/:clientID", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.ListMattersByClientHandler(matterService))
	matters.Get("/advocate/:advocateID", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.ListMattersByAdvocateHandler(matterService))
	matters.Get("/:matterID/timeline", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.GetMatterTimelineHandler(timelineService))

	// Hearing routes
	hearings := api.Group("/hearings", middleware.AuthMiddleware(cfg, userService))
	hearings.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.GetHearingHandler(hearingService))
	hearings.Post("/", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.CreateHearingHandler(hearingService))
	hearings.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.UpdateHearingHandler(hearingService))
	hearings.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.DeleteHearingHandler(hearingService))
	hearings.Get("/matter/:matterID", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.ListHearingsByMatterHandler(hearingService))

	// Order routes
	orders := api.Group("/orders", middleware.AuthMiddleware(cfg, userService))
	orders.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.GetOrderHandler(orderService))
	orders.Post("/", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.CreateOrderHandler(orderService))
	orders.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.UpdateOrderHandler(orderService))
	orders.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.DeleteOrderHandler(orderService))
	orders.Get("/matter/:matterID", middleware.RBACMiddleware(permissionService, "manage_matters"), handlers.ListOrdersByMatterHandler(orderService))

	// Document routes
	documents := api.Group("/documents", middleware.AuthMiddleware(cfg, userService))
	documents.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_documents"), handlers.GetDocumentHandler(documentService))
	documents.Post("/", middleware.RBACMiddleware(permissionService, "manage_documents"), handlers.CreateDocumentHandler(documentService))
	documents.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_documents"), handlers.UpdateDocumentHandler(documentService))
	documents.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_documents"), handlers.DeleteDocumentHandler(documentService))
	documents.Get("/matter/:matterID", middleware.RBACMiddleware(permissionService, "manage_documents"), handlers.ListDocumentsByMatterHandler(documentService))

	// Document version routes
	versions := documents.Group("/:id/versions")
	versions.Post("/", middleware.RBACMiddleware(permissionService, "manage_documents"), handlers.CreateDocumentVersionHandler(documentService))
	versions.Get("/", middleware.RBACMiddleware(permissionService, "manage_documents"), handlers.ListDocumentVersionsHandler(documentService))
	versions.Get("/:versionID", middleware.RBACMiddleware(permissionService, "manage_documents"), handlers.GetDocumentVersionHandler(documentService))

	// Search route
	api.Get("/search", middleware.AuthMiddleware(cfg, userService), middleware.RBACMiddleware(permissionService, "manage_search"), handlers.SearchHandler(searchService, llmClient))

	// AI Workspace routes — rate-limited on top of RBAC since each call can hit a
	// real LLM and is otherwise an easy abuse/cost vector.
	aiGroup := api.Group("/ai",
		middleware.AuthMiddleware(cfg, userService),
		middleware.RBACMiddleware(permissionService, "manage_ai"),
		limiter.New(limiter.Config{
			Max:        20,
			Expiration: 1 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string {
				if user, ok := c.Locals("user").(*repositories.User); ok {
					return user.ID
				}
				return c.IP()
			},
		}),
	)
	aiGroup.Post("/summarize", handlers.SummarizeHandler(aiService))
	aiGroup.Post("/ask", handlers.AskHandler(aiService))
	aiGroup.Post("/draft", handlers.DraftHandler(aiService))

	// Health check — used by the container orchestrator to know when the
	// service is actually ready to serve traffic, not just that the process is up.
	app.Get("/healthz", func(c *fiber.Ctx) error {
		if err := database.DB.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable", "error": "database unreachable"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// WebSocket route
	app.Get("/ws", websocket.WebSocketHandler(hub))

	// Start server
	logger.GetLogger().Info("Server started", zap.String("port", cfg.App.Port))

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.GetLogger().Info("Shutting down server...")
		cancel()
		app.Shutdown()
	}()

	if err := app.Listen(":" + cfg.App.Port); err != nil {
		logger.GetLogger().Fatal("Failed to start server", zap.Error(err))
	}
}