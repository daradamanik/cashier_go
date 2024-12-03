package controllers

import (
	"cashier_go/db"
	"cashier_go/models"
	"os"
	"time"

	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct { //include DB
	DB *gorm.DB
}

type Claims struct { //payload buat jwt
	Email string `json:"email"`
	Role  string `json:"role"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func (uc *UserController) Login(c *gin.Context) { //gin.context itu kayak req, res kalo di js
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil { //dijadiin json dulu
		c.JSON(400, gin.H{"error": "invalid input"}) //kayak return res.json
		return
	}
	var user models.User
	if err := uc.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{"error": "invalid email or password"})
		} else {
			c.JSON(500, gin.H{"error": "internal server error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil { //cocokin password
		c.JSON(401, gin.H{"error": "invalid email or password"})
		return
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		c.JSON(500, gin.H{"error": "Internal server error: missing secret key"})
		c.Abort()
		return
	}

	jwtExpires := time.Now().Add(24 * time.Hour) //set expire time for jwt
	Claims := &Claims{                           //payload
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(jwtExpires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"data":    gin.H{"email": user.Email, "name": user.Name, "role": user.Role},
	})
}

func (uc *UserController) AddUser(c *gin.Context) { //sama kayak req, res kalau di js
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil { //kalau error pas input user
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) //[]byte ... converts password into byte slice
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := uc.DB.Create(&user).Error; err != nil { //kalau error pas buat user baru
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(200, gin.H{"message": "user added successfuly", "user": user})
}

func (uc *UserController) AllUser(c *gin.Context) {
	var users []models.User

	if err := uc.DB.Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to retrieve users"})
		return
	}

	c.JSON(200, gin.H{
		"message": "users retrieved successfuly",
		"data":    users,
	})
}

func (uc *UserController) UserById(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := uc.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(500, gin.H{"error": "failed to fetch user"})
		return
	}
	c.JSON(200, gin.H{
		"message": "user found",
		"data":    user,
	})
}

func (uc *UserController) UserByRole(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		c.JSON(400, gin.H{"error": "Role query parameter is required"})
		return
	}

	var user []models.User

	if err := uc.DB.Where("role=?", role).Find(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "error while fetching users"})
		return
	}
	if len(user) == 0 {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	c.JSON(200, gin.H{
		"data": user,
	})
}

func (uc *UserController) SearchUser(c *gin.Context) {
	var user []models.User
	var searchData struct {
		Keyword string `form:"keyword" binding:"required"`
	}

	if err := c.ShouldBindQuery(&searchData); err != nil {
		c.JSON(401, gin.H{"error": "Invalid input"})
		return
	}

	query := uc.DB.Model(&models.User{})

	if searchData.Keyword != "" {
		query = query.Where("email LIKE ? OR name LIKE ? OR role LIKE ?",
			"%"+searchData.Keyword+"%",
			"%"+searchData.Keyword+"%",
			"%"+searchData.Keyword+"%")
	}

	if err := query.Find(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve users"})
		return
	}

	if len(user) == 0 {
		c.JSON(404, gin.H{"message": "No users found"})
		return
	}

	c.JSON(200, gin.H{"users": user})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var userData models.User
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var userFound models.User
	if err := uc.DB.First(&userFound, id).Error; err != nil {
		c.JSON(404, gin.H{"message": "user not found"})
		return
	}

	if userData.Name != "" {
		userFound.Name = userData.Name
	}
	if userData.Email != "" {
		userFound.Email = userData.Email
	}
	if userData.Role != "" {
		userFound.Role = userData.Role
	}

	if err := uc.DB.Save(&userFound).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User updated successfully",
		"user":    userFound,
	})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var userFound models.User
	if err := db.DB.First(&userFound, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "User not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	if err := db.DB.Delete(&userFound).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
