package main

import (
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/routes"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    database.Connect()

    engine := handlebars.New("./views", ".hbs")
    engine.Reload(true)
    app := fiber.New(fiber.Config{
        Views: engine,
    })

    app.Static("/assets", "./public/assets")

    // init routes
    routes.GetRoutes(app)
    routes.PostRoutes(app)
    // routes.TaskRoutes(app)

    // app.Listen(":3000")
    app.Listen("localhost:3000")
}
