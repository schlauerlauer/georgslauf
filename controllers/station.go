package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
)

type CreateStationInput struct {
	Name	string	`json:"name"	binding:"required"`
	Short	string	`json:"short"	binding:"required"`
}

type UpdateStationInput struct {
	Name	string	`json:"name"`
	Short	string	`json:"short"`
}

func GetStations(c *gin.Context) {
	var stations []models.Station
	models.DB.Find(&stations)
	c.JSON(http.StatusOK, gin.H{"data": stations})
}

func GetStation(c *gin.Context) {
	// Get model if exist
	var station models.Station
	if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": station})
}

func PostStation(c *gin.Context) {
	// Validate input
	var input CreateStationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create station
	station := models.Station{Name: input.Name, Short: input.Short}
	models.DB.Create(&station)
	c.JSON(http.StatusOK, gin.H{"data": station})
}

func PatchStation(c *gin.Context) {
	// Get model if exist
	var station models.Station
	if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var input UpdateStationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&station).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": station})
}

func DeleteStation(c *gin.Context) {
	// Get model if exist
	var station models.Station
	if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&station)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
