package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"github.com/ieraHQ/Vakalat/backend/api/validation"
)

// GetClientHandler retrieves a client by ID.
func GetClientHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		client, err := clientService.GetClientByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client not found"})
		}
		return c.JSON(client)
	}
}

// CreateClientHandler creates a new client.
func CreateClientHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var client repositories.Client
		if err := c.BodyParser(&client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := clientService.CreateClient(c.Context(), &client); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create client"})
		}
		return c.Status(fiber.StatusCreated).JSON(client)
	}
}

// UpdateClientHandler updates a client.
func UpdateClientHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var client repositories.Client
		if err := c.BodyParser(&client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		client.ID = id
		if err := validation.Validate(&client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := clientService.UpdateClient(c.Context(), &client); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update client"})
		}
		return c.JSON(client)
	}
}

// DeleteClientHandler deletes a client.
func DeleteClientHandler(clientService services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := clientService.DeleteClient(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete client"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// ListClientsHandler retrieves a list of clients.
func ListClientsHandler(clientService services.ClientService) fiber.Handler {
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
