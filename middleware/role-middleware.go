package middleware

import (
	"rizaldyaristyo-fiber-boiler/database"
	"strconv"

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

func GetPermissionVal(username string) int {
	var permission_int int
	database.DB.QueryRow("SELECT permission_val FROM users WHERE username = ?", username).Scan(&permission_int)
	return permission_int
}

func getNthDigit(fourDigitInteger int, nth int) int64 {
	strFourDigitInteger := strconv.Itoa(fourDigitInteger)
	if nth < 0 || nth >= len(strFourDigitInteger) {
		return -1
	}
	nthDigitChar := strFourDigitInteger[nth]
	nthDigitInt := int(nthDigitChar - '0')
	return int64(nthDigitInt)
}

// func RoleMiddleware(c *fiber.Ctx) error {
// 	username := string(c.Locals("username").(string))
// 	permission_int := GetPermissionVal(username)
// 	c.Locals("permission_val", permission_int)
// 	return c.Next()
// }

func RoleSalesLevel1Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 0)
	if permission_int < 1 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RoleSalesLevel2Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 0)
	if permission_int < 2 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RoleSalesLevel3Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 0)
	if permission_int < 3 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RolePurchasingLevel1Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 1)
	if permission_int < 1 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RolePurchasingLevel2Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 1)
	if permission_int < 2 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RolePurchasingLevel3Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 1)
	if permission_int < 3 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RoleManufacturingLevel1Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 2)
	if permission_int < 1 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RoleManufacturingLevel2Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 2)
	if permission_int < 2 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RoleManufacturingLevel3Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 2)
	if permission_int < 3 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RoleInventoryLevel1Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 3)
	if permission_int < 1 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RoleInventoryLevel2Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 3)
	if permission_int < 2 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}

func RoleInventoryLevel3Middleware(c *fiber.Ctx) error {
	permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), 3)
	if permission_int < 3 {
		return c.Status(401).SendString("Access Denied")
	}
	c.Locals("permission_int", permission_int)
	return c.Next()
}