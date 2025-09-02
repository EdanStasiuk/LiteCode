package controllers

import (
	"net/http"
	"strconv"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserProfile handles GET /users/:id
// Returns the profile of a user by their ID, hiding sensitive fields.
func GetUserProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Hide sensitive fields
		c.JSON(http.StatusOK, gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"role":      user.Role,
			"createdAt": user.CreatedAt,
		})
	}
}

// UpdateUserProfile handles PUT /users/:id
// Updates the current user's profile (username and email). Admins can update any user's profile.
func UpdateUserProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Ensure user can only update their own profile unless admin
		userID, _ := c.Get("userID")
		var currentUser models.User
		if err := db.First(&currentUser, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if currentUser.ID != uint(toUint(id)) && currentUser.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update other user's profile"})
			return
		}

		var body struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if body.Username != "" {
			user.Username = body.Username
		}
		if body.Email != "" {
			user.Email = body.Email
		}

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		})
	}
}

// GetUserStats handles GET /users/:id/stats
// Returns statistics for a user, including problems solved, current streak, and acceptance rate.
func GetUserStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Return stats directly from the User model
		c.JSON(http.StatusOK, gin.H{
			"userId":         user.ID,
			"problemsSolved": user.ProblemsSolved,
			"streak":         user.Streak,
			"acceptanceRate": user.AcceptanceRate,
		})
	}
}

// GetSubmissionByID handles GET /submissions/:id
// Returns a submission's details and code for the authenticated user by submission ID.
func toUint(s string) uint {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(i)
}
