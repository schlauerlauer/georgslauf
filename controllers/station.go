package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
)

type CreateStationInput struct {
	Name	string	`json:"name"	binding:"required"`
	Short	string	`json:"short"	binding:"required"`
	TribeID	uint	`json:"TribeID"	binding:"required"`
	Size	uint	`json:"size"	binding:"required"`
}

type UpdateStationInput struct {
	Name	string	`json:"name"`
	Short	string	`json:"short"`
	TribeID	uint	`json:"TribeID"`
	Size	uint	`json:"size"`
}

func GetStations(c *gin.Context) {
	var stations []models.Station
	result := models.DB.Find(&stations)
	if result.Error != nil {
		c.AbortWithStatus(500)
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(result.RowsAffected, 10))
		c.JSON(http.StatusOK, stations)
	}
}

func GetStation(c *gin.Context) {
	// Get model if exist
	var station models.Station
	if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, station)
}

func PostStation(c *gin.Context) {
	// Validate input
	var input CreateStationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} //TODO error checking (e.g. unique error)
	// Create station
	station := models.Station{
		Name: input.Name,
		Short: input.Short,
		TribeID: input.TribeID,
		Size: input.Size,
	}
	models.DB.Create(&station)
	c.JSON(http.StatusOK, station)
}

func PutStation(c *gin.Context) {
	// Validate input
	var input models.Station
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// TODO log error
		return
	}
	// Put Tribe
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
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
	c.JSON(http.StatusOK, station)
}

func DeleteStation(c *gin.Context) {
	// Get model if exist
	var station models.Station
	if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&station)
	c.JSON(http.StatusOK, true)
}
