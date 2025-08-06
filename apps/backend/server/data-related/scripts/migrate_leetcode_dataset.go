package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
)

func main() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("NEON_DEV_DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to connect to db: %v\n", err)
		return
	}

	fmt.Println("Connected successfully!")

	if err := db.AutoMigrate(
		&models.Problem{},
		&models.Tag{},
		&models.Hint{},
		&models.SimilarProblem{},
	); err != nil {
		log.Fatal("migration error:", err)
	}

	file, err := os.Open("../data/leetcode_problems.csv")
	if err != nil {
		log.Fatal("could not open CSV:", err)
	}
	defer file.Close()

	var csvProblems []*models.ScrapedProblem
	if err := gocsv.UnmarshalFile(file, &csvProblems); err != nil {
		log.Fatal("could not parse CSV:", err)
	}

	for _, p := range csvProblems {
		problem := models.Problem{
			Title:            p.Title,
			Slug:             p.Slug,
			Difficulty:       p.Difficulty,
			URL:              p.URL,
			DescriptionURL:   p.DescriptionURL,
			Description:      p.Description,
			PaidOnly:         p.PaidOnly,
			FrontendID:       p.FrontendID,
			AcceptanceRate:   p.AcceptanceRate,
			Category:         p.Category,
			Likes:            p.Likes,
			Dislikes:         p.Dislikes,
			Stats:            p.Stats,
			SolutionURL:      p.SolutionURL,
			SolutionSummary:  p.SolutionSummary,
			SolutionCodePy:   p.SolutionCodePy,
			SolutionCodeJava: p.SolutionCodeJava,
			SolutionCodeCpp:  p.SolutionCodeCpp,
			SolutionCodeURL:  p.SolutionCodeURL,
		}

		if err := db.Create(&problem).Error; err != nil {
			log.Println("failed to insert problem:", err)
			continue
		}

		// Add Hints
		hintLines := strings.Split(p.Hints, "\n")
		for _, line := range hintLines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				hint := models.Hint{
					ProblemID: problem.ID,
					Content:   trimmed,
				}
				if err := db.Create(&hint).Error; err != nil {
					log.Println("failed to insert hint:", err)
				}
			}
		}

		// Add SimilarProblems
		similarSlugs := strings.Split(p.SimilarQuestions, ",")
		for _, slug := range similarSlugs {
			slug = strings.TrimSpace(slug)
			if slug == "" {
				continue
			}

			var similar models.Problem
			if err := db.Where("slug = ?", slug).First(&similar).Error; err != nil {
				log.Printf("could not find similar problem for slug '%s': %v", slug, err)
				continue
			}

			sim := models.SimilarProblem{
				ProblemID: problem.ID,
				SimilarID: similar.ID,
			}
			if err := db.Create(&sim).Error; err != nil {
				log.Println("failed to insert similar problem:", err)
			}
		}
	}

	log.Println("Import completed successfully.")
}
