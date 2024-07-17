package router

import (
	"paywatcher/controller"
	"paywatcher/middleware"

	"github.com/gofiber/fiber/v2"
)

// iniciar las rutas. toma el app de main.go
func Init(app *fiber.App) {
	// group para crear grupos de rutas que comparten algo en comun y middlewares.
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", controller.Login)

	user := api.Group("/user")
	user.Post("/", controller.CreateUser)
	user.Get("/", middleware.ProtectedHandler(), controller.GetUser)
	user.Get("/:id", middleware.ProtectedHandler(), controller.GetUser)
	user.Put("/:id", middleware.ProtectedHandler(), controller.UpdateUser)
	user.Delete("/:id", middleware.ProtectedHandler(), controller.DeleteUser)
}
