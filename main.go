package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/th3khan/Go-Authentication-with-Fiber/database"
	"github.com/th3khan/Go-Authentication-with-Fiber/models"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidSignupRequest(c *fiber.Ctx, req *SignupRequest) bool {
	if err := c.BodyParser(req); err == nil {

		if req.Name == "" || req.Email == "" || req.Password == "" {
			return false
		}

		return true
	}
	return false
}

func ValidLoginRequest(c *fiber.Ctx, req *LoginRequest) bool {
	if err := c.BodyParser(req); err == nil {

		if req.Email == "" || req.Password == "" {
			return false
		}

		return true
	}
	return false
}

func main() {
	app := fiber.New()

	engine, err := database.CreateDBEngine()
	if err != nil {
		panic(err)
	}

	app.Post("/signup", func(c *fiber.Ctx) error {
		req := new(SignupRequest)

		// Validate Request
		if err := ValidSignupRequest(c, req); !err {
			c.JSON(fiber.Map{
				"success": false,
				"message": "Invalid request",
			})
			return c.SendStatus(fiber.StatusBadRequest)
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		user := &models.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: string(hashPass),
		}

		// Insert user into database
		if _, err := engine.Insert(user); err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Create Token
		token, exp, err := createJwtToken(*user)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.JSON(fiber.Map{
			"success": true,
			"message": "User created successfully",
			"token":   token,
			"exp":     exp,
			"user":    user,
		})
		return c.SendStatus(fiber.StatusCreated)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		req := new(LoginRequest)

		// Validate Request
		if err := ValidLoginRequest(c, req); !err {
			c.JSON(fiber.Map{
				"success": false,
				"message": "Invalid credentials",
			})
			return c.SendStatus(fiber.StatusBadRequest)
		}

		user := new(models.User)
		has, err := engine.Where("email = ?", req.Email).Desc("id").Get(user)

		if err != nil || !has {
			c.JSON(fiber.Map{
				"success": false,
				"message": "Invalid credentials",
			})
			return c.SendStatus(fiber.StatusBadRequest)
		}

		token, exp, err := createJwtToken(*user)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.JSON(fiber.Map{
			"success": true,
			"message": "Login Successfully",
			"token":   token,
			"exp":     exp,
			"user":    user,
		})
		return c.SendStatus(fiber.StatusCreated)
	})

	app.Get("/private", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"path":    "private",
		})
	})

	app.Get("/public", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"path":    "public",
		})
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}

func createJwtToken(user models.User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 5).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.Id
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret_token"))

	return t, exp, err
}
