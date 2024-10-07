package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetAllPostFormInputs(c *fiber.Ctx) error {
	formData := c.Request().PostArgs()
	
	var output string
	formData.VisitAll(func(key, value []byte) {
		output += string(key) + ": " + string(value) + "\n"
	})

	return c.SendString(output)
}