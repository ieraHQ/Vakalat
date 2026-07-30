package main

import (
	"github.com/gofiber/fiber/v2"
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

	// Initialize services
	userService := services.NewUserService(userRepo)
	clientService := services.NewClientService(clientRepo)
	matterService := services.NewMatterService(matterRepo)
	documentService := services.NewDocumentService(documentRepo)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize Fiber app
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware())

	// API routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", loginHandler(userService))

	// User routes
	users := api.Group("/users", middleware.AuthMiddleware(userService))
	users.Get("/:id", getUserHandler(userService))
	users.Post("/", createUserHandler(userService))
	users.Put("/:id", updateUserHandler(userService))
	users.Delete("/:id", deleteUserHandler(userService))

	// Client routes
	clients := api.Group("/clients", middleware.AuthMiddleware(userService))
	clients.Get("/:id", getClientHandler(clientService))
	clients.Post("/", createClientHandler(clientService))
	clients.Put("/:id", updateClientHandler(clientService))
	clients.Delete("/:id", deleteClientHandler(clientService))
	clients.Get("/", listClientsHandler(clientService))

	// Matter routes
	matters := api.Group("/matters", middleware.AuthMiddleware(userService))
	matters.Get("/:id", getMatterHandler(matterService))
	matters.Post("/", createMatterHandler(matterService))
	matters.Put("/:id", updateMatterHandler(matterService))
	matters.Delete("/:id", deleteMatterHandler(matterService))
	matters.Get("/client/:clientID", listMattersByClientHandler(matterService))
	matters.Get("/advocate/:advocateID", listMattersByAdvocateHandler(matterService))

	// Document routes
	documents := api.Group("/documents", middleware.AuthMiddleware(userService))
	documents.Get("/:id", getDocumentHandler(documentService))
	documents.Post("/", createDocumentHandler(documentService))
	documents.Put("/:id", updateDocumentHandler(documentService))
	documents.Delete("/:id", deleteDocumentHandler(documentService))
	documents.Get("/matter/:matterID", listDocumentsByMatterHandler(documentService))

	// WebSocket route
	app.Get("/ws", websocket.WebSocketHandler(hub))

	// Start server
	logger.GetLogger().Info("Server started", zap.String("port", cfg.App.Port))
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		logger.GetLogger().Fatal("Failed to start server", zap.Error(err))
	}
}

// Handler functions (placeholder implementations)
func loginHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Implement login logic
		return c.JSON(fiber.Map{"token": "example-token"})
	}
}

func getUserHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		user, err := userService.GetUserByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.JSON(user)
	}
}

func createUserHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user repositories.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := userService.CreateUser(c.Context(), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
		}
		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

func updateUserHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var user repositories.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		user.ID = id
		if err := userService.UpdateUser(c.Context(), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
		}
		return c.JSON(user)
	}
}

func deleteUserHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := userService.DeleteUser(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func getClientHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		client, err := clientService.GetClientByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client not found"})
		}
		return c.JSON(client)
	}
}

func createClientHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var client repositories.Client
		if err := c.BodyParser(&client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := clientService.CreateClient(c.Context(), &client); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create client"})
		}
		return c.Status(fiber.StatusCreated).JSON(client)
	}
}

func updateClientHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var client repositories.Client
		if err := c.BodyParser(&client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		client.ID = id
		if err := clientService.UpdateClient(c.Context(), &client); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update client"})
		}
		return c.JSON(client)
	}
}

func deleteClientHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := clientService.DeleteClient(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete client"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func listClientsHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)
		clients, err := clientService.ListClients(c.Context(), limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list clients"})
		}
		return c.JSON(clients)
	}
}

func getMatterHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		matter, err := matterService.GetMatterByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Matter not found"})
		}
		return c.JSON(matter)
	}
}

func createMatterHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var matter repositories.Matter
		if err := c.BodyParser(&matter); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := matterService.CreateMatter(c.Context(), &matter); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create matter"})
		}
		return c.Status(fiber.StatusCreated).JSON(matter)
	}
}

func updateMatterHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var matter repositories.Matter
		if err := c.BodyParser(&matter); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		matter.ID = id
		if err := matterService.UpdateMatter(c.Context(), &matter); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update matter"})
		}
		return c.JSON(matter)
	}
}

func deleteMatterHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := matterService.DeleteMatter(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete matter"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func listMattersByClientHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientID := c.Params("clientID")
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)
		matters, err := matterService.ListMattersByClient(c.Context(), clientID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list matters"})
		}
		return c.JSON(matters)
	}
}

func listMattersByAdvocateHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		advocateID := c.Params("advocateID")
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)
		matters, err := matterService.ListMattersByAdvocate(c.Context(), advocateID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list matters"})
		}
		return c.JSON(matters)
	}
}

func getDocumentHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		document, err := documentService.GetDocumentByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}
		return c.JSON(document)
	}
}

func createDocumentHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var document repositories.Document
		if err := c.BodyParser(&document); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := documentService.CreateDocument(c.Context(), &document); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create document"})
		}
		return c.Status(fiber.StatusCreated).JSON(document)
	}
}

func updateDocumentHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var document repositories.Document
		if err := c.BodyParser(&document); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		document.ID = id
		if err := documentService.UpdateDocument(c.Context(), &document); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}
		return c.JSON(document)
	}
}

func deleteDocumentHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := documentService.DeleteDocument(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete document"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func listDocumentsByMatterHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		matterID := c.Params("matterID")
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)
		documents, err := documentService.ListDocumentsByMatter(c.Context(), matterID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list documents"})
		}
		return c.JSON(documents)
	}
}