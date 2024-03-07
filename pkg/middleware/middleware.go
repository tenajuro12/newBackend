package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tenajuro12/newBackend/pkg/util"
)

func IsAuthenticate(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	userID, err := util.Parsejwt(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "Token parsing failed",
		})
	}

	// Optionally, you could set the user ID in the context for use in subsequent handlers
	c.Locals("userID", userID)

	return c.Next()
}
