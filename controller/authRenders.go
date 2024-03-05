package controller

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"text/template"
)

// RenderRegisterPage renders the register page using HTML template
func RenderRegisterPage(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(ui.TemplatePath("register.html")) // Use the TemplatePath function to get the correct path
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to parse HTML template",
		})
	}

	// Render the template to a string
	renderedTemplate := &bytes.Buffer{}
	if err := tmpl.Execute(renderedTemplate, nil); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to render HTML template",
		})
	}

	// Send the rendered template as a response
	return c.SendString(renderedTemplate.String())
}

// RenderLoginPage renders the login page using HTML template
func RenderLoginPage(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles("./ui/html/login.html")
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to parse HTML template",
		})
	}

	// Render the template to a string
	renderedTemplate := &bytes.Buffer{}
	if err := tmpl.Execute(renderedTemplate, nil); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to render HTML template",
		})
	}

	// Send the rendered template as a response
	return c.SendString(renderedTemplate.String())
}
