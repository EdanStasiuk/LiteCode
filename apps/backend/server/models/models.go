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
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Slug        string `gorm:"unique;not null"`
	Difficulty  string `gorm:"not null"` // Easy / Medium / Hard
	Description string `gorm:"type:text"`
	Constraints string `gorm:"type:text"`
	CreatedAt   time.Time
	TestCases   []TestCase
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
