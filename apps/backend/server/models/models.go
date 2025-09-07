package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string  `gorm:"unique;not null" json:"username"`
	Email          string  `gorm:"unique;not null" json:"email"`
	PasswordHash   string  `gorm:"not null" json:"password_hash"`
	Role           string  `gorm:"not null" json:"role"` // user or admin, admin users must be manually set or via a seed script
	ProblemsSolved int     `gorm:"default:0" json:"problems_solved"`
	Streak         int     `gorm:"default:0" json:"streak"`
	AcceptanceRate float64 `gorm:"default:0" json:"acceptance_rate"` // 0.0 to 1.0
}

type Problem struct {
	gorm.Model
	Title          string  `gorm:"not null" json:"title"`
	Slug           string  `gorm:"unique;not null" json:"titleSlug"`
	URL            string  `json:"url"`
	DescriptionURL string  `json:"description_url"`
	Description    string  `gorm:"type:text" json:"description"`
	Difficulty     string  `gorm:"not null" json:"difficulty"`
	Category       string  `json:"category"`
	PaidOnly       bool    `json:"paidOnly"`
	FrontendID     int     `gorm:"uniqueIndex;check:frontend_id > 0" json:"frontendQuestionId"`
	AcceptanceRate float64 `json:"acceptance_rate"`
	Stats          string  `gorm:"type:text" json:"stats"`
	Likes          int     `json:"likes"`
	Dislikes       int     `json:"dislikes"`

	SolutionURL      string `json:"solution_url"`
	SolutionSummary  string `gorm:"type:text" json:"solution"`
	SolutionCodePy   string `gorm:"type:text" json:"solution_code_python"`
	SolutionCodeJava string `gorm:"type:text" json:"solution_code_java"`
	SolutionCodeCpp  string `gorm:"type:text" json:"solution_code_cpp"`
	SolutionCodeURL  string `json:"solution_code_url"`

	Tags            []Tag            `gorm:"many2many:problem_tags;constraint:OnDelete:CASCADE;" json:"topics"`
	TestCases       []TestCase       `gorm:"constraint:OnDelete:CASCADE;" json:"test_cases"`
	Hints           []Hint           `gorm:"constraint:OnDelete:CASCADE;" json:"hints"`
	SimilarProblems []SimilarProblem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"similar_questions"`
}

type TestCase struct {
	gorm.Model
	ProblemID uint   `gorm:"not null;index" json:"problem_id"`
	Input     string `gorm:"type:text" json:"input"`
	Expected  string `gorm:"type:text" json:"expected"`
	IsSample  bool   `gorm:"type:bool" json:"is_sample"`
}

type Language struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"` // e.g. "python", "java", "cpp"
}

type Tag struct {
	gorm.Model
	Name string    `gorm:"unique;not null" json:"name"`
	Tags []Problem `gorm:"many2many:problem_tags;constraint:OnDelete:CASCADE;" json:"topics"`
}

type Hint struct {
	gorm.Model
	ProblemID uint   `gorm:"not null;index"`
	Content   string `gorm:"type:text;not null"`
}

type SimilarProblem struct {
	gorm.Model
	ProblemID uint `gorm:"not null;index" json:"problem_id"`
	SimilarID uint `gorm:"not null;index" json:"similar_id"`
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
