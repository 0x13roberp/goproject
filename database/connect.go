package database

import (
	"fmt"
	"paywatcher/config"
	"paywatcher/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// variable de la dataabase usando gorm como orm
var DB *gorm.DB

// conexion a la base de datos
func Connect() {
	var err error
	//variable para conectar los datos con el dsn

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DataBase.User,
		config.DataBase.Pass,
		config.DataBase.Host,
		config.DataBase.Port,
		config.DataBase.Name)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database because: %s!", err.Error()))
	}
	fmt.Println("Database connected!")

	// migrar esquemas
	fmt.Println("Migrating the schema...")
	DB.AutoMigrate(model.User{})    // migrar a la db nuestro modelo de usuario
	DB.AutoMigrate(model.Payment{}) // migrar a la db nuestro modelo de pagos
	DB.AutoMigrate(model.Category{})
	fmt.Println("Schema migrated!")
}
