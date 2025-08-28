package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
}

type Problem struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Slug           string `gorm:"unique;not null"`
	URL            string
	DescriptionURL string
	Description    string `gorm:"type:text"`
	Difficulty     string `gorm:"not null"` // Easy / Medium / Hard
	Category       string
	PaidOnly       bool
	FrontendID     int `gorm:"uniqueIndex"`
	AcceptanceRate float64
	Stats          string `gorm:"type:text"`
	Likes          int
	Dislikes       int

	SolutionURL      string
	SolutionSummary  string `gorm:"type:text"`
	SolutionCodePy   string `gorm:"type:text"`
	SolutionCodeJava string `gorm:"type:text"`
	SolutionCodeCpp  string `gorm:"type:text"`
	SolutionCodeURL  string

	Tags            []Tag            `gorm:"many2many:problem_tags;"`
	TestCases       []TestCase       `gorm:"constraint:OnDelete:CASCADE;"`
	Hints           []Hint           `gorm:"constraint:OnDelete:CASCADE;"`
	SimilarProblems []SimilarProblem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;"`
}

type TestCase struct {
	gorm.Model
	ProblemID uint   `gorm:"not null;index"`
	Input     string `gorm:"type:text"`
	Expected  string `gorm:"type:text"`
	IsSample  bool
}

type Language struct {
	gorm.Model
	Name string `gorm:"unique;not null"` // e.g. "python", "java", "cpp"
}

type Tag struct {
	gorm.Model
	Name     string    `gorm:"unique;not null"`
	Problems []Problem `gorm:"many2many:problem_tags;"`
}

type Hint struct {
	gorm.Model
	ProblemID uint   `gorm:"not null;index"`
	Content   string `gorm:"type:text;not null"`
}

type SimilarProblem struct {
	gorm.Model
	ProblemID uint `gorm:"not null;index"` // current problem
	SimilarID uint `gorm:"not null;index"` // ID of similar problem
}

// Used for ingesting data from CSV, not for database modeling
type ScrapedProblem struct {
	Difficulty       string  `csv:"difficulty"`
	FrontendID       int     `csv:"frontendQuestionId"`
	PaidOnly         bool    `csv:"paidOnly"`
	Title            string  `csv:"title"`
	Slug             string  `csv:"titleSlug"`
	URL              string  `csv:"url"`
	DescriptionURL   string  `csv:"description_url"`
	Description      string  `csv:"description"`
	SolutionURL      string  `csv:"solution_url"`
	SolutionSummary  string  `csv:"solution"`
	SolutionCodePy   string  `csv:"solution_code_python"`
	SolutionCodeJava string  `csv:"solution_code_java"`
	SolutionCodeCpp  string  `csv:"solution_code_cpp"`
	SolutionCodeURL  string  `csv:"solution_code_url"`
	Category         string  `csv:"category"`
	AcceptanceRate   float64 `csv:"acceptance_rate"`
	Topics           string  `csv:"topics"`
	Hints            string  `csv:"hints"`
	Likes            int     `csv:"likes"`
	Dislikes         int     `csv:"dislikes"`
	SimilarQuestions string  `csv:"similar_questions"`
	Stats            string  `csv:"stats"`
}
