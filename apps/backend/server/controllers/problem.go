package controllers

import (
	"net/http"
	"strconv"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProblemByID handles GET /problems/:id
// Returns a single problem by its ID, including its tags.
func GetProblemByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var problem models.Problem
		if err := db.Preload("Tags").First(&problem, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
			return
		}
		c.JSON(http.StatusOK, problem)
	}
}

// GetProblems handles GET /problems
// Returns a list of all problems, each with their tags.
func GetProblems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var problems []models.Problem
		if err := db.Preload("Tags").Find(&problems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch problems"})
			return
		}
		c.JSON(http.StatusOK, problems)
	}
}

// CreateProblem handles POST /problems
// Creates a new problem from the provided JSON body.
func CreateProblem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var problem models.Problem
		if err := c.ShouldBindJSON(&problem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&problem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create problem"})
			return
		}
		c.JSON(http.StatusCreated, problem)
	}
}

// UpdateProblem handles PUT /problems/:id
// Updates an existing problem by its ID with the provided JSON body.
func UpdateProblem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var problem models.Problem
		if err := db.First(&problem, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
			return
		}
		if err := c.ShouldBindJSON(&problem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Save(&problem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update problem"})
			return
		}
		c.JSON(http.StatusOK, problem)
	}
}

// DeleteProblem handles DELETE /problems/:id
// Deletes a problem by its ID.
func DeleteProblem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		if err := db.Delete(&models.Problem{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete problem"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Problem deleted"})
	}
}

// RestoreProblem handles POST /problems/:id/restore
// Restores a soft-deleted problem by clearing its deleted_at timestamp.
func RestoreProblem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var problem models.Problem
		// Include soft-deleted rows
		if err := db.Unscoped().First(&problem, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
			return
		}

		// Restore using UpdateColumn with zero value
		if err := db.Unscoped().Model(&problem).UpdateColumn("deleted_at", gorm.DeletedAt{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore problem"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Problem restored"})
	}
}
