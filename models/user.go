package models

type User struct {
	ID                    uint   `gorm:"primaryKey"`
	Username              string `gorm:"unique;not null"`
	Password              string `gorm:"not null"`
	CompanyName           string
	CompanyEmail          string
	CompanyAddress        string
	CompanyPhone          string
	LetterheadPath        string
	SelectedFields        string                 // JSON-encoded string
	FinancialInstitutions []FinancialInstitution `gorm:"foreignKey:UserID"`
	TermsConditions       string                 `gorm:"type:text"`
}

type FinancialInstitution struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null"`
	Name       string `gorm:"not null"`
	BankNumber string
	Link       string
}
