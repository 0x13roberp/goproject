package router

import (
	"paywatcher/controller"

	"github.com/gofiber/fiber/v3"
)

// iniciar las rutas. toma el app de main.go
func Init(app *fiber.App) {
	// group para crear grupos de rutas que comparten algo en comun y middlewares.
	api := app.Group("/api")
	user := api.Group("/user")

	user.Get("/", controller.GetUser)
	user.Get("/:id", controller.GetUser)
	user.Post("/", controller.CreateUser)
	user.Put("/:id", controller.UpdateUser)
	user.Delete("/:id", controller.DeleteUser)
}
