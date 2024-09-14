package handlers

import (
	"fmt"
	"net/mail"
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
    // Parse form data or JSON
    var user models.User
    contentType := c.Get("Content-Type")
    if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
        user.UserID = c.FormValue("employee_id")
        user.Username = c.FormValue("username")
        user.Password = c.FormValue("password")
        user.Email = c.FormValue("email")
        user.PermissionVal = (c.FormValue("sales_permission") + c.FormValue("purchasing_permission") + c.FormValue("manufacturing_permission") + c.FormValue("inventory_permission"))
        fmt.Println(user.PermissionVal)
    } else {
        if err := c.BodyParser(&user); err != nil {
            // return c.Status(400).SendString(err.Error())
            return c.Render("register", fiber.Map{
                "registerMessage": "Invalid data",
            })
        }
    }

    // Check if email valid
    _, err := mail.ParseAddress(user.Email)
    if err != nil {
        // return c.Status(400).SendString("Invalid email")
        return c.Render("register", fiber.Map{
            "registerMessage": "Invalid email",
        })
    }

    // Check if user exists
    rows, err := database.DB.Query("SELECT username FROM users WHERE username = ? OR email = ?", user.Username, user.Email)
    if err != nil {
        return c.Render("register", fiber.Map{
            "registerMessage": "Could not check if user exists",
        })
    }
    defer rows.Close()
    if rows.Next() {
        // User already exists
        return c.Render("register", fiber.Map{
            "registerMessage": "User or E-mail already exists",
        })
    }

    // bcrypt hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        // return c.Status(500).SendString("Could not hash the password")
        return c.Render("register", fiber.Map{
            "registerMessage": "Could not hash the password",
        })
    }
    user.Password = string(hashedPassword)

    // init transaction
    tx, err := database.DB.Begin()
    if err != nil {
        return c.Render("register", fiber.Map{
            "registerMessage": "Could not create user, Transaction error",
        })
    }
    
    _, err = tx.Exec(`INSERT INTO users (user_id, username, password, email, permission_val) VALUES (?, ?, ?, ?, ?)`, user.UserID, user.Username, user.Password, user.Email, user.PermissionVal)
    if err != nil {
        tx.Rollback()
        return c.Render("register", fiber.Map{
            "registerMessage": "Could not create user, Insert error",
        })
    }
    
    // commit transaction
    err = tx.Commit()
    if err != nil {
        return c.Render("register", fiber.Map{
            "registerMessage": "Could not create user, Commit error",
        })
    }

    return c.Redirect("/login?register=success")
}
