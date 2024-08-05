package controller

import (
	"paywatcher/database"
	"paywatcher/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetCategories(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	userID := getIdUserInToken(c)

	if id != "" {
		return getCategoryByID(c, db, id)
	}
	return getAllCategories(c, db, userID)

}

func getAllCategories(c *fiber.Ctx, db *gorm.DB, userID int) error {
	categories := []model.Category{}
	if err := db.Where("user_id = ?", userID).Find(categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "Error", "message": "Error while getting categories"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": categories, "message": "Found categories"})
}

// esta funcion puede ser privada porque solamente la vamos a usar en este file
func getCategoryByID(c *fiber.Ctx, db *gorm.DB, id string) error {
	var category = model.Category{}

	if err := db.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "Error", "message": "Category not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": category, "message": "Found categories"})
}

// post /categories
func CreateCategory(c *fiber.Ctx) error {
	db := database.DB
	category := new(model.Category)

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "Error", "message": "JSON Incompleted"})
	}

	if err := db.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "Error", "message": "Error while creating the category"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": category, "message": "Category Created"})
}

// put /categories
