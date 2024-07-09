package main

import (
	"paywatcher/database"
	"paywatcher/router"

	"github.com/gofiber/fiber/v3"
)

func main() {
	database.Connect() // conectar con la db

	app := fiber.New() // instancia de fiber

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hii\n")
	})

	// iniciar el router creado en router.go
	router.Init(app)

	app.Listen(":3000")
}
