package controllers

import (
	"net/http"
	"strconv"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/dto"
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
		if err := db.Preload("Tags").Order("frontend_id ASC").Find(&problems).Error; err != nil {
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
		var input dto.ProblemInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert topics strings to Tag models
		var tags []models.Tag
		for _, t := range input.Topics {
			var tag models.Tag
			if err := db.Where("name = ?", t).First(&tag).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					tag = models.Tag{Name: t}
					db.Create(&tag)
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query tags"})
					return
				}
			}
			tags = append(tags, tag)
		}

		// Handle hints
		var hints []models.Hint
		for _, h := range input.Hints {
			hints = append(hints, models.Hint{Content: h})
		}

		// Create problem
		problem := models.Problem{
			Title:            input.Title,
			Slug:             input.Slug,
			Difficulty:       input.Difficulty,
			Category:         input.Category,
			FrontendID:       input.FrontendID,
			PaidOnly:         input.PaidOnly,
			Description:      input.Description,
			DescriptionURL:   input.DescriptionURL,
			URL:              input.URL,
			SolutionURL:      input.SolutionURL,
			SolutionSummary:  input.SolutionSummary,
			SolutionCodePy:   input.SolutionCodePy,
			SolutionCodeJava: input.SolutionCodeJava,
			SolutionCodeCpp:  input.SolutionCodeCpp,
			SolutionCodeURL:  input.SolutionCodeURL,
			AcceptanceRate:   input.AcceptanceRate,
			Stats:            input.Stats,
			Likes:            input.Likes,
			Dislikes:         input.Dislikes,
			Tags:             tags,
			Hints:            hints,
			TestCases:        input.TestCases,
		}

		if err := db.Create(&problem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create problem"})
			return
		}

		var similarProblems []models.SimilarProblem
		for _, s := range input.SimilarProblems {
			var similar models.Problem
			if err := db.Where("slug = ?", s.TitleSlug).First(&similar).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// Skip or handle missing similar problem
					continue
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query similar problems"})
					return
				}
			}

			similarProblems = append(similarProblems, models.SimilarProblem{
				ProblemID: problem.ID,
				SimilarID: similar.ID,
			})
		}

		if len(similarProblems) > 0 {
			if err := db.Create(&similarProblems).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create similar problems"})
				return
			}
		}

		if err := db.Preload("SimilarProblems").First(&problem, problem.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload problem with similar questions"})
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
		if err := db.Preload("Tags").Preload("Hints").Preload("TestCases").Preload("SimilarProblems").First(&problem, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
			return
		}

		var input dto.ProblemInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := map[string]interface{}{}

		if input.Title != "" {
			updates["title"] = input.Title
		}
		if input.Slug != "" {
			updates["slug"] = input.Slug
		}
		if input.URL != "" {
			updates["url"] = input.URL
		}
		if input.Description != "" {
			updates["description"] = input.Description
		}
		if input.DescriptionURL != "" {
			updates["description_url"] = input.DescriptionURL
		}
		if input.Difficulty != "" {
			updates["difficulty"] = input.Difficulty
		}
		if input.Category != "" {
			updates["category"] = input.Category
		}
		if input.FrontendID != 0 {
			updates["frontend_id"] = input.FrontendID
		}
		if input.AcceptanceRate != 0 {
			updates["acceptance_rate"] = input.AcceptanceRate
		}
		if input.Stats != "" {
			updates["stats"] = input.Stats
		}
		if input.Likes != 0 {
			updates["likes"] = input.Likes
		}
		if input.Dislikes != 0 {
			updates["dislikes"] = input.Dislikes
		}

		// Boolean field: use pointer to detect if it's passed
		if input.PaidOnly != problem.PaidOnly {
			updates["paid_only"] = input.PaidOnly
		}

		if err := db.Model(&problem).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update problem"})
			return
		}

		// --- Update Tags ---
		if input.Topics != nil {
			var tags []models.Tag
			for _, t := range input.Topics {
				var tag models.Tag
				if err := db.Where("name = ?", t).First(&tag).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						tag = models.Tag{Name: t}
						db.Create(&tag)
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query tags"})
						return
					}
				}
				tags = append(tags, tag)
			}
			db.Model(&problem).Association("Tags").Replace(tags)
		}

		// --- Update Hints ---
		if input.Hints != nil {
			var hints []models.Hint
			for _, h := range input.Hints {
				hints = append(hints, models.Hint{Content: h})
			}
			// Clear old hints and add new ones
			db.Model(&problem).Association("Hints").Replace(hints)
		}

		// --- Update TestCases ---
		if input.TestCases != nil {
			db.Model(&problem).Association("TestCases").Replace(input.TestCases)
		}

		// --- Update SimilarProblems ---
		if input.SimilarProblems != nil {
			var similarProblems []models.SimilarProblem
			for _, s := range input.SimilarProblems {
				var similar models.Problem
				if err := db.Where("slug = ?", s.TitleSlug).First(&similar).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						continue
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query similar problems"})
						return
					}
				}
				similarProblems = append(similarProblems, models.SimilarProblem{
					ProblemID: problem.ID,
					SimilarID: similar.ID,
				})
			}
			db.Model(&problem).Association("SimilarProblems").Replace(similarProblems)
		}

		// Reload the updated problem with associations
		if err := db.Preload("Tags").Preload("Hints").Preload("TestCases").Preload("SimilarProblems").First(&problem, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload problem"})
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
