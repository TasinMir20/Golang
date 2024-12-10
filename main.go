package main

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New()

	// Define a simple GET route
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the Fiber API!",
		})
	})

	// Define another route
	app.Get("/greet/:name", func(c fiber.Ctx) error {
		name := c.Params("name")
		return c.JSON(fiber.Map{
			"greeting": "Hello, " + name + "!",
		})
	})

	// CPU intensive route
	app.Get("/fibonacci/:n", func(c fiber.Ctx) error {
		n, err := strconv.Atoi(c.Params("n"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Please provide a valid number",
			})
		}

		// Recursive Fibonacci calculation (intentionally CPU intensive)
		var fib func(int) int
		fib = func(n int) int {
			if n <= 1 {
				return n
			}
			return fib(n-1) + fib(n-2)
		}

		result := fib(n)
		return c.JSON(fiber.Map{
			"number": n,
			"result": result,
		})
	})

	// Start the server on port 3000
	if err := app.Listen(":4000"); err != nil {
		panic(err)
	}

}
