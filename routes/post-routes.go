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
	app.Post("/api/GetInventory", middleware.JWTMiddleware, middleware.RoleInventoryLevel1Middleware, handlers.GetInventory)
	app.Post("/api/GetVendors", middleware.JWTMiddleware, middleware.RoleInventoryLevel1Middleware, handlers.GetVendors)
	app.Post("/api/ReplenishInventory/:inventory_id", middleware.JWTMiddleware, middleware.RoleInventoryLevel2Middleware, handlers.ReplenishInventory)
	app.Post("/api/EditInventory/:inventory_id", middleware.JWTMiddleware, middleware.RoleInventoryLevel2Middleware, handlers.EditInventory)
	// DEV
	app.Get("/api/GetInventory", middleware.JWTMiddleware, middleware.RoleInventoryLevel1Middleware, handlers.GetInventory)
	app.Get("/api/GetVendors", middleware.JWTMiddleware, middleware.RoleInventoryLevel1Middleware, handlers.GetVendors)
	
	app.Post("/api/GetmanufacturingOrder", middleware.JWTMiddleware, middleware.RoleManufacturingLevel1Middleware, handlers.GetManufacturingOrder)
	app.Post("/api/ApproveManufacturingOrder/:manufacturing_order_id", middleware.JWTMiddleware, middleware.RoleManufacturingLevel3Middleware, handlers.ApproveManufacturingOrder)
	app.Post("/api/ReceiveManufacturingOrder/:manufacturing_order_id", middleware.JWTMiddleware, middleware.RoleManufacturingLevel3Middleware, handlers.ReceiveManufacturingOrder)
	app.Post("/api/CancelManufacturingOrder/:manufacturing_order_id", middleware.JWTMiddleware, middleware.RoleManufacturingLevel3Middleware, handlers.CancelManufacturingOrder)
	// DEV
	app.Get("/api/GetmanufacturingOrder", middleware.JWTMiddleware, middleware.RoleManufacturingLevel1Middleware, handlers.GetManufacturingOrder)
}