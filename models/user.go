package models

type User struct {
    ID             uint   `gorm:"primaryKey"`
    Username       string `gorm:"unique;not null"`
    Password       string `gorm:"not null"`
    SelectedFields string // Store selected fields as JSON or a string
    LetterheadPath string // Path to the uploaded company letterhead
}
