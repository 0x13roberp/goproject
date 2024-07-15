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
	result := db.First(&user, id)

	// si el id no existe. cuando creamos la estructura por defecto se pone a 0
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "Error", "message": "User not found!"})
	}

	// si el id si existe
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": user, "message": "User found!"})
}

// traer todos los usuarios
func GetAllUsers(c fiber.Ctx, db *gorm.DB) error {
	var users []model.User
	// guardar todos los datos de la db en el array
	db.Find(&users)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": users})
}

// POST /user
func CreateUser(c fiber.Ctx) error {
	db := database.DB
	// crear un nuevo modelo de user
	user := new(model.User)

	// comprobar que estan viniendo todos los campos. BodyParser deprecated, usar bind body pasandole el puntero del user
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "Error", "message": "Body JSON incompleted!"})
	}

	// cifrar la pass usando la funcion creada abajo
	password, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Error", "message": "Error while hashing password!"})
	}
	// si no hubo ningun error cifrandola, le asignamos la pass cifrada al user
	user.Password = password

	// guardar la informacion del usuario en la db
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "Error", "message": "Error while saving data to the database!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Success", "data": user, "message": "User created!"})
}

// PUT /user/:id
func UpdateUser(c fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user model.User

	type updateUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	// uu = user update
	var uu updateUser

	if err := c.Bind().Body(&uu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "Error", "message": "Error while updating user!"})
	}

	// para actualizar la pass del user, tambien tenemos que cifrarla
	password, err := hashPassword(uu.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Error", "message": "Error while hashing password!"})
	}

	// encontrar el user
	result := db.First(&user, id)

	// error al actualizar el user
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "Error", "message": "User not found!"})
	}

	// cambiar los datos del usuario con el de updated user
	user.Name = uu.Name
	user.Email = uu.Email
	user.UserName = uu.UserName
	user.Password = password // pass actualizada

	// guardar el nuevo user actualizado
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Error", "message": "Error while updating user!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": user, "message": "User updated!"})
}

// DELETE /user/:id
func DeleteUser(c fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user model.User

	// comprobar que existe el user
	result := db.First(&user, id)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "Error", "message": "User not found!"})
	}

	// eliminar el user de la db
	result = db.Delete(&user)

	// si ocurre un error al eliminar el user
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "Error", "message": "Error while deleting user!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "data": user, "message": "User deleted!"})
}
