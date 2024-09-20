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

func RoleMiddleware(role string, level int) fiber.Handler{
	return func(c *fiber.Ctx) error {
		var nthDigit int
		switch role{
		case "Sales":
			nthDigit = 0
		case "Purchasing":
			nthDigit = 1
		case "Manufacturing":
			nthDigit = 2
		case "Inventory":
			nthDigit = 3
		default:
			nthDigit = -1
		}

		permission_int := getNthDigit(GetPermissionVal(string(c.Locals("username").(string))), nthDigit)
		if permission_int < 0 {
			return c.Status(401).SendString("Access Denied")
		}
		c.Locals("permission_int", permission_int)
		return c.Next()
	}
}