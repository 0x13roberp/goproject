package main

import (
	"paywatcher/database"
	"paywatcher/middleware"
	"paywatcher/router"

	"github.com/gofiber/fiber/v3"
)

func main() {
	database.Connect() // conectar con la db

	app := fiber.New() // instancia de fiber

	app.Use(middleware.ProtectedHandler()) // usar la funcion creada para hacer un login correcto

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hii\n")
	})

	// iniciar el router creado en router.go
	router.Init(app)

	app.Listen(":3000")
}
