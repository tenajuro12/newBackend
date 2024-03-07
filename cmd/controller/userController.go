package controller

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tenajuro12/newBackend/pkg/database"
	models2 "github.com/tenajuro12/newBackend/pkg/models"
	"github.com/tenajuro12/newBackend/pkg/util"
	"gorm.io/gorm"
	"strconv"
)

// DeleteUser deletes a user account.
func DeleteUser(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := util.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Check if the user exists
	var user models2.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	if err := database.DB.Where("user_id = ?", userID).Delete(&models2.Blog{}).Error; err != nil {
		return c.JSON(fiber.Map{
			"message": "Failed to delete associated blogs",
		})
	}

	// Delete user from database
	if err := database.DB.Delete(&user).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User account deleted successfully",
	})
}

func UpdateUser(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := util.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Parse request body to extract updated user data
	var updatedUser models2.User
	if err := c.BodyParser(&updatedUser); err != nil {
		fmt.Println("Unable to parse user body")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid user payload",
		})
	}

	// Check if the user exists
	var user models2.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Update user information
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.Email = updatedUser.Email
	user.Phone = updatedUser.Phone

	// Save updated user to database
	if err := database.DB.Save(&user).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User information updated successfully",
		"user":    user,
	})
}

func GetUserInfo(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := util.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Query the database for user information
	var user models2.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Return user information
	return c.JSON(user)
}

// FollowUser allows a user to follow another user
func FollowUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	followerIDStr, err := util.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Convert followerID to a uint
	followerID, err := strconv.ParseUint(followerIDStr, 10, 64)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Parse followed user ID from request parameters
	followedUserID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	// Check if the user is trying to follow themselves
	if followerID == followedUserID {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Cannot follow yourself",
		})
	}
	// Check if the followed user exists
	var followedUser models2.User
	if err := database.DB.First(&followedUser, followedUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Followed user not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Check if the follow relationship already exists
	var followRelationship models2.Follow
	if err := database.DB.Where("follower_id = ? AND followed_user_id = ?", followerID, followedUserID).First(&followRelationship).Error; err == nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Already following this user",
		})
	}

	// Create a new follow relationship
	follow := models2.Follow{
		FollowerID:     uint(followerID),
		FollowedUserID: uint(followedUserID),
	}
	if err := database.DB.Create(&follow).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to follow user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully followed user",
	})
}

// UnfollowUser allows a user to unfollow another user
func UnfollowUser(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	followerID, err := util.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Parse followed user ID from request parameters
	followedUserID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	// Check if the follow relationship exists
	var followRelationship models2.Follow
	if err := database.DB.Where("follower_id = ? AND followed_user_id = ?", followerID, followedUserID).First(&followRelationship).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "Not following this user",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Delete the follow relationship
	if err := database.DB.Delete(&followRelationship).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to unfollow user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully unfollowed user",
	})
}

// In controller/userController.go

// GetPostsFromFollowedUsers retrieves posts from users that the current user is following
func GetPostsFromFollowedUsers(c *fiber.Ctx) error {
	// Extract user ID from JWT cookie
	cookie := c.Cookies("jwt")
	userID, err := util.Parsejwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Get list of users the current user is following
	var followedUsers []models2.Follow
	if err := database.DB.Where("follower_id = ?", userID).Find(&followedUsers).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to retrieve followed users",
		})
	}

	// Extract followed user IDs
	var followedUserIDs []uint
	for _, follow := range followedUsers {
		followedUserIDs = append(followedUserIDs, follow.FollowedUserID)
	}

	// Retrieve posts from followed users
	var blogs []models2.Blog
	if err := database.DB.Where("user_id IN (?)", followedUserIDs).Find(&blogs).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to retrieve posts from followed users",
		})
	}

	return c.JSON(fiber.Map{
		"posts": blogs,
	})
}
