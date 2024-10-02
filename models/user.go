package models

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Username       string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	CompanyName    string
	CompanyAddress string
	CompanyPhone   string
	LetterheadPath string
	SelectedFields string // JSON-encoded string
}
