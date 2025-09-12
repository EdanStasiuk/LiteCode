package dto

import "github.com/EdanStasiuk/LiteCode/apps/backend/server/models"

// Data Transfer Objects

type ProblemInput struct {
	Title          string  `json:"title"`
	Slug           string  `json:"titleSlug"`
	URL            string  `json:"url"`
	DescriptionURL string  `json:"description_url"`
	Description    string  `json:"description"`
	Difficulty     string  `json:"difficulty"`
	Category       string  `json:"category"`
	PaidOnly       *bool   `json:"paidOnly"`
	FrontendID     int     `json:"frontendQuestionId"`
	AcceptanceRate float64 `json:"acceptance_rate"`
	Stats          string  `json:"stats"`
	Likes          int     `json:"likes"`
	Dislikes       int     `json:"dislikes"`

	SolutionURL      string `json:"solution_url"`
	SolutionSummary  string `json:"solution"`
	SolutionCodePy   string `json:"solution_code_python"`
	SolutionCodeJava string `json:"solution_code_java"`
	SolutionCodeCpp  string `json:"solution_code_cpp"`
	SolutionCodeURL  string `json:"solution_code_url"`

	Topics          []string               `json:"topics"`
	TestCases       []models.TestCase      `json:"test_cases"`
	Hints           []string               `json:"hints"`
	SimilarProblems []SimilarQuestionInput `json:"similar_questions"`
}

type SimilarQuestionInput struct {
	Title     string `json:"title"`
	TitleSlug string `json:"titleSlug"`
}
