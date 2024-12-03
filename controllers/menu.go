package controllers

import (
	// "cashier_go/db"
	"cashier_go/db"
	"cashier_go/models"
	"fmt"
	"path/filepath"
	"strconv"
	"os"
	"time"

	// "fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MenuController struct { //include DB
	DB *gorm.DB
}

func (mc *MenuController) AddMenu(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	file, err := c.FormFile("picture")
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to get the picture file"})
		return
	}

	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !validTypes[file.Header.Get("Content-Type")] {
		c.JSON(400, gin.H{"error": "Only JPEG and PNG images are allowed"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(400, gin.H{"error": "Invalid file type. Only JPG and PNG are allowed"})
		return
	}

	filePath := "./uploads/" + time.Now().Format("20060102_150405") + "_" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "failed to save picture"})
		return
	}
	menuName := c.PostForm("menu_name")
	menuTypeInput := c.PostForm("type")
	var menuType models.Jenis
	switch menuTypeInput {
	case string(models.Makanan):
		menuType = models.Makanan
	case string(models.Minuman):
		menuType = models.Minuman
	default:
		c.JSON(400, gin.H{"error": "Invalid type. Must be 'makanan' or 'minuman'"})
		return
	}

	description := c.PostForm("description")
	price := c.PostForm("price")
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid price. Must be a valid number"})
		return
	}

	menu := models.Menu{
		MenuName:    menuName,
		Type:        menuType,
		Picture:     filePath, // Store the file path
		Description: description,
		Price:       priceFloat,
	}
	if err := mc.DB.Create(&menu).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to create menu"})
		return
	}
	c.JSON(200, gin.H{"message": "menu added successfuly", "menu": menu})
}

func (mc *MenuController) AllMenu(c *gin.Context) {
	var menus []models.Menu

	if err := mc.DB.Find(&menus).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "menu retrieved successfuly", "menu": menus})
}

func (mc *MenuController) MenuByID(c *gin.Context) {
	id := c.Param("id")

	var menu models.Menu
	if err := mc.DB.First(&menu, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "menu not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "menu retrieved succesfuly", "menu": menu})
}

func (mc *MenuController) SearchMenu(c *gin.Context) {
	var menus []models.Menu
	var searchData struct {
		Keyword string `form:"keyword" binding:"required"`
	}

	if err := c.ShouldBindQuery(&searchData); err != nil {
		c.JSON(401, gin.H{"error": "invalid input"})
		return
	}
	query := mc.DB.Model(&models.Menu{})

	if searchData.Keyword != "" {
		query = query.Where(
			"menu_name LIKE ? OR description LIKE ?",
			"%"+searchData.Keyword+"%",
			"%"+searchData.Keyword+"%",
		)
	}

	if err := query.Find(&menus).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to retrieve menus"})
		return
	}

	if len(menus) == 0 {
		c.JSON(404, gin.H{"error": "menu not found"})
		return
	}

	c.JSON(200, gin.H{"menu": menus})
}

func (mc *MenuController) ByType(c *gin.Context) {
	var menus []models.Menu
	menuType := c.Query("type")
	if menuType == "" {
		c.JSON(404, gin.H{"error": "query is empty"})
		return
	}

	if err := mc.DB.Where("type=?", menuType).Find(&menus).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(menus) == 0 {
		c.JSON(404, gin.H{"error": "menu not found"})
		return
	}

	c.JSON(200, gin.H{"message": "menus retrieved successfully", "menu": menus})
}

func (mc *MenuController) UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var menuFound models.Menu

	if err := mc.DB.First(&menuFound, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "menu not found"})
		return
	}

	menuName := c.PostForm("menu_name")
	description := c.PostForm("description")
	price := c.PostForm("price")
	menuType := c.PostForm("type")

	if menuName != "" {
		menuFound.MenuName = menuName
	}
	if description != "" {
		menuFound.Description = description
	}
	if price != "" {
		parsedPrice, err := strconv.ParseFloat(price, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid price format"})
			return
		}
		menuFound.Price = parsedPrice
	}
	if menuType != "" {
		switch menuType {
		case string(models.Makanan), string(models.Minuman):
			menuFound.Type = models.Jenis(menuType)
		default:
			c.JSON(400, gin.H{"error": "Invalid menu type. Must be 'makanan' or 'minuman'"})
			return
		}
	}

	file, err := c.FormFile("picture")
	if err == nil {
		filePath := fmt.Sprintf("./uploads/%s", file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save picture"})
			return
		}
		menuFound.Picture = filePath
	} else {
		c.JSON(400, gin.H{"error": "Failed to process file upload"})
		return
	}

	if err := mc.DB.Save(&menuFound).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update menu"})
		return
	}

	c.JSON(200, gin.H{"message": "menu updated successfuly", "updated menu": menuFound})
}

func (mc *MenuController) DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	var menuFound models.Menu
	if err := db.DB.First(&menuFound, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error":"menu not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if menuFound.Picture != "" {
        if err := os.Remove(menuFound.Picture); err != nil {
            c.JSON(500, gin.H{"error": "Failed to delete picture"})
            return
        }
    }

	if err := mc.DB.Delete(&menuFound).Error; err != nil {
		c.JSON(500, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200, gin.H{"message":"menu deleted successfuly"})
}