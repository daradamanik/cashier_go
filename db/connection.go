package db
import (
    "gorm.io/driver/mysql"
	"github.com/joho/godotenv"
    "gorm.io/gorm"
    "log"
    "os"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get DB_URL from environment
    dsn := os.Getenv("DB_URL")
    if dsn == "" {
        log.Fatal("DB_URL not set in .env file")
    }

    // Connect to the database
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

    log.Println("Database connection established")
}