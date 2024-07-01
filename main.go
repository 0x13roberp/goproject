package main

import (
	"paywatcher/database"

	"github.com/gofiber/fiber/v3"
)

func main() {
	database.Connect() // conectar con la db

	app := fiber.New() // instancia de fiber

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hii")
	})

	app.Listen(":3000")
}
