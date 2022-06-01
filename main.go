package main

import (
	"log"

	"rajivharlalka/imagery-v2/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/app", routes.RootRoute)
	app.Post("/activity-route", routes.ActivityRoute)

	log.Fatal(app.Listen(":3000"))
}
