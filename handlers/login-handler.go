package handlers

import (
	"os"
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
    // load secret from env
    jwtSecret := []byte(os.Getenv("JWT_SECRET"))

    // parse form data or JSON
    var user models.UserLogin
    contentType := c.Get("Content-Type")
    if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
        user.Username = c.FormValue("username")
        user.Password = c.FormValue("password")
    } else {
        err := c.BodyParser(&user)
        if err != nil {
            // return c.Status(400).SendString(err.Error())
            return c.Render("login", fiber.Map{
                "loginMessage": "Invalid username or password",
            })
        }
    }

    // check if user exists
    var storedPassword string
    err := database.DB.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedPassword)
    if err != nil {
        // return c.Status(401).SendString("Invalid username or password")
        return c.Render("login", fiber.Map{
            "loginMessage": "Invalid username or password",
        })
    }

    // compare password
    err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
    if err != nil {
        // return c.Status(401).SendString("Invalid username or password")
        return c.Render("login", fiber.Map{
            "loginMessage": "Invalid username or password",
        })
    }

    // generate token
    token := jwt.New(jwt.SigningMethodHS256)
    texp := time.Now().Add(time.Hour * 72) // token expiration: 72 hours
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = user.Username
    claims["exp"] = texp.Unix()

    t, err := token.SignedString(jwtSecret)
    if err != nil {
        return c.Status(500).SendString("Could not login")
    }

    // return c.JSON(fiber.Map{"token": t,})

    // Set the JWT cookie
    jwtCookie := new(fiber.Cookie)
    jwtCookie.Name = "jwt"
    jwtCookie.Value = t
    jwtCookie.Expires = texp
    jwtCookie.HTTPOnly = true
    jwtCookie.SameSite = "Lax"
    jwtCookie.Secure = c.Protocol() == "https" // use in dev
    // jwtCookie.Secure = true // use in prod
    c.Cookie(jwtCookie)

    return c.Redirect("/dashboard?login=success")
}

// for dev
func DevGetJWT(c *fiber.Ctx) error {

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = "dev"
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
    t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return c.Status(500).SendString("Could not login")
    }

    return c.Status(200).SendString(t)
}