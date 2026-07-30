package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/auth"
	"github.com/ieraHQ/Vakalat/backend/api/config"
	"github.com/ieraHQ/Vakalat/backend/api/database"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/middleware"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"github.com/ieraHQ/Vakalat/backend/api/websocket"
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

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	clientRepo := repositories.NewClientRepository(database.DB)
	matterRepo := repositories.NewMatterRepository(database.DB)
	documentRepo := repositories.NewDocumentRepository(database.DB)
	refreshTokenRepo := auth.NewRefreshTokenRepository(database.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	clientService := services.NewClientService(clientRepo)
	matterService := services.NewMatterService(matterRepo)
	documentService := services.NewDocumentService(documentRepo)
	permissionService := auth.NewPermissionService(userRepo)
	sessionService := auth.NewSessionService(userRepo)

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
	authGroup.Post("/login", loginHandler(cfg, userService, refreshTokenRepo))
	authGroup.Post("/refresh", refreshTokenHandler(cfg, refreshTokenRepo, userService))
	authGroup.Post("/forgot-password", forgotPasswordHandler(sessionService))
	authGroup.Post("/reset-password", resetPasswordHandler(sessionService))

	// User routes
	users := api.Group("/users", middleware.AuthMiddleware(cfg, userService))
	users.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_users"), getUserHandler(userService))
	users.Post("/", middleware.RBACMiddleware(permissionService, "manage_users"), createUserHandler(userService))
	users.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_users"), updateUserHandler(userService))
	users.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_users"), deleteUserHandler(userService))

	// Client routes
	clients := api.Group("/clients", middleware.AuthMiddleware(cfg, userService))
	clients.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_clients"), getClientHandler(clientService))
	clients.Post("/", middleware.RBACMiddleware(permissionService, "manage_clients"), createClientHandler(clientService))
	clients.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_clients"), updateClientHandler(clientService))
	clients.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_clients"), deleteClientHandler(clientService))
	clients.Get("/", middleware.RBACMiddleware(permissionService, "manage_clients"), listClientsHandler(clientService))

	// Matter routes
	matters := api.Group("/matters", middleware.AuthMiddleware(cfg, userService))
	matters.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), getMatterHandler(matterService))
	matters.Post("/", middleware.RBACMiddleware(permissionService, "manage_matters"), createMatterHandler(matterService))
	matters.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), updateMatterHandler(matterService))
	matters.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_matters"), deleteMatterHandler(matterService))
	matters.Get("/client/:clientID", middleware.RBACMiddleware(permissionService, "manage_matters"), listMattersByClientHandler(matterService))
	matters.Get("/advocate/:advocateID", middleware.RBACMiddleware(permissionService, "manage_matters"), listMattersByAdvocateHandler(matterService))

	// Document routes
	documents := api.Group("/documents", middleware.AuthMiddleware(cfg, userService))
	documents.Get("/:id", middleware.RBACMiddleware(permissionService, "manage_documents"), getDocumentHandler(documentService))
	documents.Post("/", middleware.RBACMiddleware(permissionService, "manage_documents"), createDocumentHandler(documentService))
	documents.Put("/:id", middleware.RBACMiddleware(permissionService, "manage_documents"), updateDocumentHandler(documentService))
	documents.Delete("/:id", middleware.RBACMiddleware(permissionService, "manage_documents"), deleteDocumentHandler(documentService))
	documents.Get("/matter/:matterID", middleware.RBACMiddleware(permissionService, "manage_documents"), listDocumentsByMatterHandler(documentService))

	// WebSocket route
	app.Get("/ws", websocket.WebSocketHandler(hub))

	// Start server
	logger.GetLogger().Info("Server started", zap.String("port", cfg.App.Port))
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		logger.GetLogger().Fatal("Failed to start server", zap.Error(err))
	}
}

// Handler functions
func loginHandler(cfg *config.Config, userService services.UserService, refreshTokenRepo auth.RefreshTokenRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		user, err := userService.GetUserByEmail(c.Context(), req.Email)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		if !auth.VerifyPassword(req.Password, user.PasswordHash, cfg) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		token, err := auth.GenerateToken(user.ID, user.RoleID, cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		refreshToken, err := auth.GenerateRefreshToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
		}

		if err := refreshTokenRepo.Create(c.Context(), refreshToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save refresh token"})
		}

		return c.JSON(fiber.Map{
			"token":         token,
			"refresh_token": refreshToken.Token,
		})
	}
}

func refreshTokenHandler(cfg *config.Config, refreshTokenRepo auth.RefreshTokenRepository, userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type RefreshRequest struct {
			RefreshToken string `json:"refresh_token"`
		}

		var req RefreshRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		refreshToken, err := refreshTokenRepo.FindByToken(c.Context(), req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
		}

		if time.Now().After(refreshToken.ExpiresAt) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token expired"})
		}

		user, err := userService.GetUserByID(c.Context(), refreshToken.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		token, err := auth.GenerateToken(user.ID, user.RoleID, cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		// Generate new refresh token
		newRefreshToken, err := auth.GenerateRefreshToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
		}

		// Delete old refresh token
		if err := refreshTokenRepo.Delete(c.Context(), req.RefreshToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete refresh token"})
		}

		// Save new refresh token
		if err := refreshTokenRepo.Create(c.Context(), newRefreshToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save refresh token"})
		}

		return c.JSON(fiber.Map{
			"token":         token,
			"refresh_token": newRefreshToken.Token,
		})
	}
}

func forgotPasswordHandler(sessionService auth.SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type ForgotPasswordRequest struct {
			Email string `json:"email"`
		}

		var req ForgotPasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		token, err := sessionService.CreatePasswordResetToken(c.Context(), req.Email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create reset token"})
		}

		// Placeholder: Send email with reset token
		return c.JSON(fiber.Map{"token": token})
	}
}

func resetPasswordHandler(sessionService auth.SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type ResetPasswordRequest struct {
			Token       string `json:"token"`
			NewPassword string `json:"new_password"`
		}

		var req ResetPasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if err := sessionService.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Remaining handler functions (getUserHandler, createUserHandler, etc.) remain unchanged...