package middleware

import (
	"rizaldyaristyo-fiber-boiler/database"

	"github.com/gofiber/fiber/v2"
)

// PermissionVal
//
// 0 = Not Permitted
// 1 = View
// 2 = Add
// 3 = Approve
//
// 0 0 0 0
// │ │ │ └── Inventory
// │ │ └──── Manufacturing
// │ └────── Purchasing
// └──────── Sales
//
// PermissionVal Example:
// 1101 = View Sales, View Purchasing, No Manufacturing, View Inventory
// 1211 = View Sales, Add Purchasing, View Manufacturing, View Inventory
// 0312 = No Sales, Add Purchasing, View Manufacturing, View Inventory

func RoleMiddleware(c *fiber.Ctx) error {
	// get username from c.Locals
	username := string(c.Locals("username").(string))

	// get user permissions from database
	var permission_int int
	database.DB.QueryRow("SELECT permission_val FROM users WHERE username = ?", username).Scan(&permission_int)

	c.Locals("permission_val", permission_int)
	return c.Next()
	// return c.JSON(fiber.Map{"permission_val": permission_int})
}