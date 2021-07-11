package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ConversationHandler struct {
	ConversationService services.ConversationService
}

func (ch *ConversationHandler) FindConversation(c *fiber.Ctx) error {
	c.Accepts("application/json")

	token := c.Get("Authorization")

	var auth domain.Authentication

	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	username := c.Params("username")

	conversation, err := ch.ConversationService.FindConversation(u.Username, username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": conversation.Messages})
}

func (ch *ConversationHandler) GetConversationPreviews(c *fiber.Ctx) error {
	c.Accepts("application/json")

	token := c.Get("Authorization")

	var auth domain.Authentication

	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	conversation, err := ch.ConversationService.GetConversationPreviews(u.Username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": conversation})
}
