package middleware

import (
	// "fmt"
	// "strings"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(c *fiber.Ctx) error {

	// // get token from Authorization header
	// authHeader := c.Get("Authorization")
	// if authHeader == "" {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "missing token"})
	// }
	// tokenString := strings.TrimPrefix(authHeader, "Bearer ") // split Bearer from token

	// get token from the "jwt" HTTPOnly-cookie
	tokenString := c.Cookies("jwt")
	if tokenString == "" {
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "missing token"})
		c.ClearCookie("jwt")
		return c.Redirect("/login?login=missing_token")
	}

	// parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// HS256 can just be validated using HMAC 
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
		}
		// all okay, return secret key from env
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// check parsing error or invalid token
	if err != nil || !token.Valid {
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid or expired token", "error": err.Error()})
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid or expired token"})
		c.ClearCookie("jwt")
		return c.Redirect("/login?login=invalid_token")
	}

	// check claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token claims"})
		c.ClearCookie("jwt")
		return c.Redirect("/login?login=invalid_token")
	}

	// set locals to pass to the handler, add if any
	// access via c.Locals("<key-name>")
	c.Locals("username", claims["username"])

	// fmt.Printf("token: %+v\n", token)
	// fmt.Printf("claims: %+v\n", claims)

	return c.Next()
}
