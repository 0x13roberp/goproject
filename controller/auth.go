package controller

import (
	"paywatcher/config"
	"paywatcher/database"
	"paywatcher/model"
	"time"

	"github.com/gofiber/fiber/v3"
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

func Login(c fiber.Ctx) error {

	// Estructura requerida para hacer log in
	var loginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	// Si ocurre un error al hacer log in returnar un bad request
	if err := c.Bind().Body(&loginInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "Error",
			"message": "Invalid Request!"})
	}

	// Usamos la estructura de user para guardar los datos del user haciendo log in y despues compararlos con la de la base de datos
	var user model.User

	// Traer la base de datos para comparar los datos del user guardado con el introducido al hacer log in
	var DB = database.DB

	// En el caso de que no coincida el nombre de usuario con el que esta en la base de datos. tambien podriamos usar el email
	if err := DB.Where("username = ? ", loginInput.Identity).First(&user); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "Error",
			"message": "Invalid Credentials!"})
	}

	// Ahora checkear la si las passwords son iguales. tambien las encriptamos en el proceso
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "Error", "message": "Internal error!"})
	}

	// Crear el claims para generar el token de inicio de sesion
	claims := jwt.MapClaims{"name": user.Name, "admin": false, "exp": time.Now().Add(time.Hour * 72).Unix()}

	// Token de inicio de sesion con su hash
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Genera el token ya cifrado
	t, err := token.SignedString([]byte(config.SecretJWTKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "Error",
			"message": "Internal error!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "message": "Logged in correctly!", "data": t})
}
