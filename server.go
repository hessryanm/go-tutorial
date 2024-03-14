package main

import(
	"log"

	"github.com/gofiber/fiber/v2"
)

func indexHandler(c *fiber.Ctx) error {
	return c.SendString("Not Implemented")
}

func postHandler(c *fiber.Ctx) error {
	return c.SendString("Not Implemented")
}

func putHandler(c *fiber.Ctx) error {
	return c.SendString("Not Implemented")
}

func deleteHandler(c *fiber.Ctx) error {
	return c.SendString("Not Implemented")
}

func main() {
	app := fiber.New()

	app.Get("/", indexHandler)
	app.Post("/", postHandler)
	app.Put("/update", putHandler)
	app.Delete("/delete", deleteHandler)

	log.Fatal(app.Listen(":3000"))
}