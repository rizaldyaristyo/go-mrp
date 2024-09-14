package routes

import (
	"rizaldyaristyo-fiber-boiler/handlers"
	"rizaldyaristyo-fiber-boiler/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetRoutes(app *fiber.App) {
    // Public Routes
    app.Get("/", handlers.Index)
    app.Get("/register", handlers.RegisterHbs)
    app.Get("/login", handlers.LoginHbs)

    // Protected Routes
    app.Get("/logout", middleware.JWTMiddleware, handlers.Logout)
    app.Get("/dashboard", middleware.JWTMiddleware, handlers.DashboardHbs)
    
    // Role Specific Routes
    app.Get("/sales", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.SalesHbs)
    app.Get("/purchasing", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.PurchasingHbs)
    app.Get("/manufacturing", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.ManufacturingHbs)
    app.Get("/inventory", middleware.JWTMiddleware, middleware.RoleMiddleware, handlers.InventoryHbs)
}