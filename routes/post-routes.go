package routes

import (
	"rizaldyaristyo-fiber-boiler/handlers"
	"rizaldyaristyo-fiber-boiler/middleware"

	"github.com/gofiber/fiber/v2"
)

func PostRoutes(app *fiber.App) {
    // Public Routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Protected Routes
	app.Post("/logout", middleware.JWTMiddleware, handlers.Logout)

	// Role Specific Routes
	app.Post("/api/GetInventory", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.GetInventory)
	app.Post("/api/ReplenishInventory/:inventory_id", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.ReplenishInventory)
	// DEV
	app.Get("/api/getinventory", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.GetInventory)
	
	app.Post("/api/GetmanufacturingOrder", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.GetManufacturingOrder)
	// DEV
	app.Get("/api/GetmanufacturingOrder", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.GetManufacturingOrder)
}
