package router

import (
	"example.com/app/handlers"
	"example.com/app/repo"
	"example.com/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {

	mh := handlers.MessageHandler{MessageService: services.NewMessageService(repo.NewMessageRepoImpl())}
	ch := handlers.ConversationHandler{ConversationService: services.NewConversationService(repo.NewConversationRepoImpl())}

	app.Use(recover.New())

	api := app.Group("", logger.New())

	messages := api.Group("/messages")
	messages.Post("/", mh.CreateMessage)
	messages.Delete("/multi", mh.DeleteByIDs)
	messages.Delete("/", mh.DeleteByID)

	conversations := api.Group("/conversation")
	conversations.Get("/:username", ch.FindConversation)
	conversations.Get("/", ch.GetConversationPreviews)
}

func Setup() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		ExposeHeaders: "Authorization",
	}))

	SetupRoutes(app)
	return app
}
