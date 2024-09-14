package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	return c.SendFile("./public/index.html")
	// return c.Status(200).SendString("Hello, World!")
}

func LoginHbs(c *fiber.Ctx) error {
	// check if user already logged in
	if c.Cookies("jwt") != "" {
        return c.Redirect("/dashboard")
	} else {
    	loginStatus := c.Query("login")
		switch loginStatus {
		case "invalid_token":
			return c.Render("login", fiber.Map{
				"loginMessage": "Session Expired, Try to Login Again",
			})

		case "missing_token":
			return c.Render("login", fiber.Map{
				"loginMessage": "Please Login First",
			})

		case "wrong_login":
			return c.Render("login", fiber.Map{
				"loginMessage": "Wrong Username or Password!",
			})

		default:
			return c.Render("login", nil)
		}
	}
}

func RegisterHbs(c *fiber.Ctx) error {
	return c.Render("register", nil)
}

func DashboardHbs(c *fiber.Ctx) error {
	// check "login" URL query
    loginStatus := c.Query("login")
	if loginStatus == "success" {
		return c.Render("dashboard", fiber.Map{
			"loginMessage": "Login success!",
		})
	} else {
		return c.Render("dashboard", nil)
	}
}

func SalesHbs(c *fiber.Ctx) error {	
	permissionInt := getNthDigit(c.Locals("permission_val").(int), 0)

	if permissionInt > 0 {
		return c.Render("sales", fiber.Map{
			"permission_val": permissionIntToString(permissionInt),
			"permission_int": permissionInt,
		})
	}

	return c.Redirect("/dashboard")
}

func PurchasingHbs(c *fiber.Ctx) error {
	permissionInt := getNthDigit(c.Locals("permission_val").(int), 1)

	if permissionInt > 0 {
		return c.Render("purchasing", fiber.Map{
			"permission_val": permissionIntToString(permissionInt),
			"permission_int": permissionInt,
		})
	}

	return c.Redirect("/dashboard")
}

func ManufacturingHbs(c *fiber.Ctx) error {
	permissionInt := getNthDigit(c.Locals("permission_val").(int), 2)

	if permissionInt > 0 {
		return c.Render("manufacturing", fiber.Map{
			"permission_val": permissionIntToString(permissionInt),
			"permission_int": permissionInt,
		})
	}

	return c.Redirect("/dashboard")
}

func InventoryHbs(c *fiber.Ctx) error {
	permissionInt := getNthDigit(c.Locals("permission_val").(int), 3)

	if permissionInt > 0 {
		return c.Render("inventory", fiber.Map{
			"permission_val": permissionIntToString(permissionInt),
			"permission_int": permissionInt,
		})
	}

	return c.Redirect("/dashboard")
}