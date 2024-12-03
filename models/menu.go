package models

type Jenis string

const (
	Makanan Jenis = "makanan"
	Minuman Jenis = "minuman"
)

type Menu struct {
	MenuID				uint					`gorm:"primaryKey;autoIncrement"`
	MenuName			string					`gorm:"size:100;not null"`
	Type 				Jenis					`gorm:"type:enum('makanan','minuman');not null"`
	Picture				string					`gorm:"not null"`
	Description			string					`gorm:"size:255;not null"`
	Price				float64					`gorm:"not null"`
	TransactionDetails	[]TransactionDetail		`gorm:"foreignKey:IDMenu"`
}