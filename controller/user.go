package controller

import (
	"paywatcher/database"
	"paywatcher/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// funciones para crear crud de user

// GET /user/:id
// traer usuario sin id
func GetUser(c fiber.Ctx) error {
	// traer la db
	db := database.DB
	// guardar el parametro id de fiber
	id := c.Params("id")

	// si el user no tiene id vacio
	if id != "" {
		return GetUserById(c, db, id)
	}

	// si el user tiene id vacio
	return GetAllUsers(c, db)

}

// traer usuario por id
func GetUserById(c fiber.Ctx, db *gorm.DB, id string) error {
	// instancia de nuestra estructura user creada en models
	var user model.User
	// guardar en la variable user el id
	result := db.Find(&user, id)

	// si el id no existe. cuando creamos la estructura por defecto se pone a 0
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "Error", "message": "User not found!"})
	}

	// usuario pero sin la password. ya que lo vamos a retornar como json
	userReturn := model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		UserName: user.UserName,
	}

	// si el id si existe
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": userReturn, "message": "User found!"})
}

// traer todos los usuarios
func GetAllUsers(c fiber.Ctx, db *gorm.DB) error {
	var users []model.User
	// guardar todos los datos de la db en el array
	db.Find(&users)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": users})
}

// POST /user

// PUT /user/:id

// DELETE /user/:id
