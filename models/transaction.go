package models
import "time"

type Status string
const (
	Paid Status = "paid"
	Unpaid Status = "Unpaid"
)

type Transaction struct {
	TransactionID		uint 				`gorm:"primaryKey;autoIncrement"`		
	Date				time.Time			`gorm:"not null"`
	IDUser				uint				`gorm:"not null"`
	User				User				`gorm:"foreignKey:IDUser"`
	Customer			string				`gorm:"size:100;not null"`
	Status				Status				`gorm:"type:enum('paid','unpaid');not null"`
	CreatedAt			time.Time			`gorm:"autoCreateTime"`
	TransactionDetails	[]TransactionDetail `gorm:"foreignKey:IDTransaction"`
}