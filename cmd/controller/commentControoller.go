package controller

import (
	"errors"
	"fmt"
	"github.com/tenajuro12/newBackend/pkg/database"
	models2 "github.com/tenajuro12/newBackend/pkg/models"
	"github.com/tenajuro12/newBackend/pkg/util"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateComment(c *fiber.Ctx) error {
	var commentData models2.Comment
	if err := c.BodyParser(&commentData); err != nil {
		fmt.Println("Unable to parse comment body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid comment payload",
		})
	}

	cookie := c.Cookies("jwt")
	userID, err := util.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}

	var blogPost models2.Blog
	if err := database.DB.Where("id = ?", postID).First(&blogPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Blog post not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	comment := models2.Comment{
		UserID:   userID,
		PostID:   uint(postID),
		Content:  commentData.Content,
		DateTime: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to create comment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

func UpdateComment(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("commentID"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid comment ID",
		})
	}

	var updatedComment models2.Comment
	if err := c.BodyParser(&updatedComment); err != nil {
		fmt.Println("Unable to parse comment body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid comment payload",
		})
	}

	var comment models2.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Comment not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	comment.Content = updatedComment.Content

	if err := database.DB.Save(&comment).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to update comment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

func DeleteComment(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("commentID"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid comment ID",
		})
	}

	var comment models2.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Comment not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to delete comment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comment deleted successfully",
	})
}

func ReadComments(c *fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}

	var blogPost models2.Blog
	if err := database.DB.Where("id = ?", postID).First(&blogPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Blog post not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	var comments []models2.Comment
	if err := database.DB.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to retrieve comments",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Comments retrieved successfully",
		"comments": comments,
	})
}
