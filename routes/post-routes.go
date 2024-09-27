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
	app.Post("/api/GetInventory", middleware.JWTMiddleware, middleware.RoleMiddleware("Inventory", 1), handlers.GetInventory)
	app.Post("/api/GetVendors", middleware.JWTMiddleware, middleware.RoleMiddleware("Inventory", 1), handlers.GetVendors)
	app.Post("/api/ReplenishInventory/:inventory_id", middleware.JWTMiddleware, middleware.RoleMiddleware("Inventory", 2), handlers.ReplenishInventory)
	app.Post("/api/EditInventory/:inventory_id", middleware.JWTMiddleware, middleware.RoleMiddleware("Inventory", 2), handlers.EditInventory)
	// app.Post("/api/GetManufacturingOrder", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.GetManufacturingOrder)
	app.Post("/api/GetManufacturingOrder", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.OptimizedGetManufacturingOrder)
	app.Post("/api/ApproveManufacturingOrder/:manufacturing_order_id", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 3), handlers.ApproveManufacturingOrder)
	app.Post("/api/ReceiveManufacturingOrder/:manufacturing_order_id", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 3), handlers.ReceiveManufacturingOrder)
	app.Post("/api/CancelManufacturingOrder/:manufacturing_order_id", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 3), handlers.CancelManufacturingOrder)
	app.Post("/api/GetRecipes", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.GetRecipes)
	app.Post("/api/EditRecipe/:product_id", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.EditRecipe)
	app.Post("/api/GetMaterials", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.GetMaterials)
	
	// DEV
	app.Get("/api/GetInventory", middleware.JWTMiddleware, middleware.RoleMiddleware("Inventory", 1), handlers.GetInventory)
	app.Get("/api/GetVendors", middleware.JWTMiddleware, middleware.RoleMiddleware("Inventory", 1), handlers.GetVendors)
	// app.Get("/api/GetmanufacturingOrder", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.GetManufacturingOrder)
	app.Get("/api/GetmanufacturingOrder2", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.OptimizedGetManufacturingOrder)
	app.Get("/api/GetRecipes", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.GetRecipes)
	app.Get("/api/GetMaterials", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.GetMaterials)
	app.Get("/dev/GetJWT", handlers.DevGetJWT)
	app.Post("/dev/GetJWT", handlers.DevGetJWT)
}