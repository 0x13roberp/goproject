package controller

import (
	"paywatcher/config"
	"paywatcher/database"
	"paywatcher/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// funcion para cifrar la pass
func hashPassword(password string) (string, error) {
	// bcrypt sirve para hashear
	CryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	// la funcion retorna byte y un error. entonces casteamos el byte a string
	return string(CryptedPass), err
}

func getIdUserInToken(c *fiber.Ctx) int {
	user := c.Locals("user").(*jwt.Token) // obtiene el token de autenticacion del contexto
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)
	return int(id)
}

// funcion para comprobar si existe un usuario
func ExistingUser(identity string) (model.User, error) {
	// Traer la base de datos para comparar los datos del user guardado con el introducido al hacer log in
	var DB = database.DB
	var user model.User

	// En el caso de que no coincida el nombre de usuario con el que esta en la base de datos. tambien podriamos usar el email
	if err := DB.Where("user_name = ?", identity).First(&user); err.Error != nil {
		return user, err.Error
	}
	return user, nil
}

// Funcion para checkear si la passwords son iguales
func CheckPassword(hash string, password string) bool {
	// Ahora checkear la si las passwords son iguales. tambien las encriptamos en el proceso
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Funcion para crear el token de inicio de sesion ya hasheado
func CreateToken(user model.User) *jwt.Token {
	claims := jwt.MapClaims{
		"name":  user.Name,
		"id":    user.ID,
		"admin": false,
		"exp":   time.Now().Add(time.Hour * 72).Unix()}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func Login(c *fiber.Ctx) error {

	// Estructura requerida para hacer log in
	var loginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	var user model.User
	var err error

	// Si ocurre un error al hacer log in returnar un bad request
	if err = c.BodyParser(&loginInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error",
			"message": "Invalid Request!"})
	}

	user, err = ExistingUser(loginInput.Identity)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error",
			"message": "Invalid Credentials!"})
	}

	if !CheckPassword(user.Password, loginInput.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error",
			"message": "Invalid Credentials!"})
	}

	token := CreateToken(user)

	// Genera el token ya cifrado
	t, err := token.SignedString([]byte(config.SecretJWTKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error",
			"message": "Internal error!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "message": "Logged in correctly!", "data": t})
}
