package models

type TransactionDetail struct {
	DetailID		uint		`gorm:"primaryKey;autoIncrement"`
	IDTransaction	uint		`gorm:"not null"`
	Transaction 	Transaction	`gorm:"foreignKey:IDTransaction"`
	IDMenu			uint		`gorm:"not null"`
	Menu 			Menu		`gorm:"foreignKey:IDMenu"`
	TotalItem		int 		`gorm:"not null"`
	TotalPrice		float64		`gorm:"not null"`
}