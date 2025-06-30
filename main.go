package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {

	fmt.Println("Hello, World!")
	app := fiber.New()
	app.Get("/abc", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! 1234")
	})
	app.Post("/sum", AddOperation)
	app.Listen(":3000")

}

func AddOperation(c *fiber.Ctx) error {
	var request AddOperationRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	result := request.Num1 + request.Num2
	response := AddOperationResponse{Result: result}

	return c.JSON(response) 
}

type AddOperationRequest struct {
	Num1 int  `json:"num1"`
	Num2 int  `json:"num2"`
}

type AddOperationResponse struct {
	Result int `json:"result"`
}
	