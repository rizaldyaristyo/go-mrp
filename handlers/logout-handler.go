package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
    // Clear the JWT cookie
    c.ClearCookie("jwt")
    return c.Redirect("/login")
}