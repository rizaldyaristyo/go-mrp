package routes

import (
	"rizaldyaristyo-fiber-boiler/handlers"
	"rizaldyaristyo-fiber-boiler/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetRoutes(app *fiber.App) {
    // Public Routes
    // app.Get("/", handlers.Index)
    app.Get("/", middleware.JWTMiddleware, func(c *fiber.Ctx) error {return c.Redirect("/dashboard")})

    app.Get("/register", handlers.RegisterHbs)
    app.Get("/login", handlers.LoginHbs)

    // Protected Routes
    app.Get("/logout", middleware.JWTMiddleware, handlers.Logout)
    app.Get("/dashboard", middleware.JWTMiddleware, handlers.DashboardHbs)
    
    // Role Specific Routes
    app.Get("/sales", middleware.JWTMiddleware, middleware.RoleMiddleware("Sales", 1), handlers.SalesHbs)
    app.Get("/purchasing", middleware.JWTMiddleware, middleware.RoleMiddleware("Purchasing", 1), handlers.PurchasingHbs)
    app.Get("/manufacturing", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.ManufacturingHbs)
    app.Get("/manufacturing/recipes", middleware.JWTMiddleware, middleware.RoleMiddleware("Manufacturing", 1), handlers.ManufacturingRecipesHbs)
    app.Get("/inventory", middleware.JWTMiddleware, middleware.RoleMiddleware("Inventory", 1), handlers.InventoryHbs)
}