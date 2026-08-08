package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"github.com/ieraHQ/Vakalat/backend/api/validation"
)

// GetOrderHandler retrieves an order by ID.
func GetOrderHandler(orderService services.OrderService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		order, err := orderService.GetOrderByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.JSON(order)
	}
}

// CreateOrderHandler creates a new order.
func CreateOrderHandler(orderService services.OrderService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var order repositories.Order
		if err := c.BodyParser(&order); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&order); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := orderService.CreateOrder(c.Context(), &order); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
		}
		return c.Status(fiber.StatusCreated).JSON(order)
	}
}

// UpdateOrderHandler updates an order.
func UpdateOrderHandler(orderService services.OrderService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var order repositories.Order
		if err := c.BodyParser(&order); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		order.ID = id
		if err := validation.Validate(&order); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := orderService.UpdateOrder(c.Context(), &order); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update order"})
		}
		return c.JSON(order)
	}
}

// DeleteOrderHandler deletes an order.
func DeleteOrderHandler(orderService services.OrderService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := orderService.DeleteOrder(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete order"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// ListOrdersByMatterHandler retrieves all orders for a matter.
func ListOrdersByMatterHandler(orderService services.OrderService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		matterID := c.Params("matterID")
		orders, err := orderService.ListOrdersByMatter(c.Context(), matterID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list orders"})
		}
		return c.JSON(orders)
	}
}

// GetMatterTimelineHandler retrieves the merged hearing/order timeline for a matter.
func GetMatterTimelineHandler(timelineService services.TimelineService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		matterID := c.Params("matterID")
		timeline, err := timelineService.GetMatterTimeline(c.Context(), matterID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get matter timeline"})
		}
		return c.JSON(timeline)
	}
}
