package models

type Role string

const (
	Admin Role = "admin"
	Kasir Role = "kasir"
)

type User struct {
	UserID       uint          `gorm:"primaryKey;autoIncrement"`
	Name         string        `gorm:"size:100;not null"`
	Email        string        `gorm:"size:200;uniqueIndex;not null"`
	Password     string        `gorm:"not null"`
	Role         Role          `gorm:"type:enum('admin', 'kasir');not null"`
	Transactions []Transaction `gorm:"foreignKey:IDUser"`
}