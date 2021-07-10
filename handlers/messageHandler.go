package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

type MessageHandler struct {
	MessageService services.MessageService
}

func (mh *MessageHandler) CreateMessage(c *fiber.Ctx) error {
	c.Accepts("application/json")

	token := c.Get("Authorization")

	var auth domain.Authentication

	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	message := new(domain.Message)

	message.From = u.Username

	err = c.BodyParser(message)

	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = mh.MessageService.Create(message)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (mh *MessageHandler) DeleteByID(c *fiber.Ctx) error {
	c.Accepts("application/json")

	token := c.Get("Authorization")

	var auth domain.Authentication

	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	message := new(domain.DeleteMessage)
	err = c.BodyParser(message)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = mh.MessageService.DeleteByID(u.Username, message.Id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}