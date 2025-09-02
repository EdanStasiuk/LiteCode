package controllers

import (
	"net/http"
	"strconv"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListTags handles GET /tags
// Returns a list of all tags.
func ListTags(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tags []models.Tag
		if err := db.Find(&tags).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tags": tags})
	}
}

// CreateTag handles POST /tags
// Creates a new tag with the provided JSON body.
func CreateTag(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
			return
		}

		tag := models.Tag{Name: input.Name}
		if err := db.Create(&tag).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"tag": tag})
	}
}

// DeleteTag handles DELETE /tags/:id
// Deletes a tag by its ID.
func DeleteTag(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
			return
		}

		if err := db.Delete(&models.Tag{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Tag deleted"})
	}
}
