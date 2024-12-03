package main

import (
	"cashier_go/controllers"
	"github.com/joho/godotenv"
	"log"
	"cashier_go/db"
	"cashier_go/models"
	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDB()

	err = db.DB.AutoMigrate(
		&models.User{},
		&models.Menu{},
		&models.Transaction{},
		&models.TransactionDetail{},
	)
	if err != nil {
		panic("Failed to run migration" + err.Error())
	}
	println("Migrations completed")

	r := gin.Default()
	userControl := controllers.UserController{DB: db.DB}
	menuControl := controllers.MenuController{DB: db.DB}

	r.POST("/user/add", userControl.AddUser)
	r.GET("/user/all-user", userControl.AllUser)
	r.GET("/user/:id", userControl.UserById)
	r.GET("/user/role", userControl.UserByRole)
	r.POST("/user/login", userControl.Login)
	r.GET("/user/search", userControl.SearchUser)
	r.PATCH("/user/update/:id", userControl.UpdateUser)
	r.DELETE("/user/delete/:id", userControl.DeleteUser)
	r.POST("/menu/add", menuControl.AddMenu)
	r.GET("/menu/all-menu", menuControl.AllMenu)
	r.GET("/menu/:id", menuControl.MenuByID)
	r.GET("/menu/search", menuControl.SearchMenu)
	r.GET("/menu/type", menuControl.ByType)
	r.PATCH("/menu/update/:id", menuControl.UpdateMenu)
	r.DELETE("/menu/delete/:id", menuControl.DeleteMenu)
	r.Run(":8000")
}
