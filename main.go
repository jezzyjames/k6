package main

import (
	"k6demo/repositories"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dial := mysql.Open("root:P@ssw0rd@tcp(localhost:3306)/arise")
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	app.Get("/products", func(c fiber.Ctx) error {
		productRepo := repositories.NewProductRepository(db)
		products, err := productRepo.FindAll(false)
		if err != nil {
			return err
		}
		return c.JSON(products)
	})

	app.Listen(":8000")
}
