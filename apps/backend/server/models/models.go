package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	Submissions  []Submission
}

type Problem struct {
	ID               uint   `gorm:"primaryKey"`
	Title            string `gorm:"not null"`
	Slug             string `gorm:"unique;not null"`
	URL              string
	DescriptionURL   string
	Description      string `gorm:"type:text"`
	Difficulty       string `gorm:"not null"` // Easy / Medium / Hard
	Category         string
	PaidOnly         bool
	FrontendID       int `gorm:"uniqueIndex"`
	AcceptanceRate   float64
	Hints            string `gorm:"type:text"`
	Likes            int
	Dislikes         int
	Stats            string `gorm:"type:text"`
	SimilarQuestions string `gorm:"type:text"`

	SolutionURL      string
	SolutionSummary  string `gorm:"type:text"`
	SolutionCodePy   string `gorm:"type:text"`
	SolutionCodeJava string `gorm:"type:text"`
	SolutionCodeCpp  string `gorm:"type:text"`
	SolutionCodeURL  string

	CreatedAt   time.Time
	Tags        []Tag `gorm:"many2many:problem_tags;"`
	Submissions []Submission
}

type TestCase struct {
	ID        uint   `gorm:"primaryKey"`
	ProblemID uint   `gorm:"not null"`
	Input     string `gorm:"type:text"`
	Expected  string `gorm:"type:text"`
	IsSample  bool
}

type Submission struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	ProblemID uint   `gorm:"not null"`
	Language  string `gorm:"not null"`
	Code      string `gorm:"type:text"`
	Status    string // Accepted, Wrong Answer, etc.
	RuntimeMs *int
	MemoryKb  *int
	CreatedAt time.Time
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type ProblemTag struct {
	ProblemID uint
	TagID     uint
}

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
