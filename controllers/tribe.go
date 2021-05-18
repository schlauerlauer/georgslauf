package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
)

type CreateTribeInput struct {
	Name	string	`json:"name"	binding:"required"`
	Short	string	`json:"short"	binding:"required"`
}

type UpdateTribeInput struct {
	Name	string	`json:"name"`
	Short	string	`json:"short"`
}

func GetTribes(c *gin.Context) {
	var tribes []models.Tribe
	result := models.DB.Find(&tribes)
	if result.Error != nil {
		c.AbortWithStatus(500)
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(result.RowsAffected, 10))
		c.JSON(http.StatusOK, tribes)
	}
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
	tribe := models.Tribe{
		Name: input.Name,
		Short: input.Short,
	}
	models.DB.Create(&tribe)
	c.JSON(http.StatusOK, tribe)
}

func PutTribe(c *gin.Context) {
	// Validate input
	var input models.Tribe
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// TODO log error
		return
	} //TODO error checking (e.g. unique error)
	// Put Tribe
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
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
