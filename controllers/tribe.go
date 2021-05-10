package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
)

type CreateTribeInput struct {
	Name	string	`json:"name"	binding:"required"`
}

type UpdateTribeInput struct {
	Name	string	`json:"name"`
}

func GetTribes(c *gin.Context) {
	var tribes []models.Tribe
	models.DB.Find(&tribes)
	c.Header("Access-Control-Expose-Headers", "X-Total-Count")
	c.Header("X-Total-Count", "10") //FIXME
	c.JSON(http.StatusOK, tribes)
}

func GetTribe(c *gin.Context) {
	// Get model if exist
	var tribe models.Tribe
	if err := models.DB.Where("id = ?", c.Param("id")).First(&tribe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, tribe)
}

func PostTribe(c *gin.Context) {
	// Validate input
	var input CreateTribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create tribe
	tribe := models.Tribe{Name: input.Name}
	models.DB.Create(&tribe)
	c.JSON(http.StatusOK, tribe)
}

func PatchTribe(c *gin.Context) {
	// Get model if exist
	var tribe models.Tribe
	if err := models.DB.Where("id = ?", c.Param("id")).First(&tribe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var input UpdateTribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&tribe).Updates(input)
	c.JSON(http.StatusOK, tribe)
}

func DeleteTribe(c *gin.Context) {
	// Get model if exist
	var tribe models.Tribe
	if err := models.DB.Where("id = ?", c.Param("id")).First(&tribe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&tribe)
	c.JSON(http.StatusOK, true)
}
