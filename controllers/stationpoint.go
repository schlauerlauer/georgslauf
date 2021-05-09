package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
)

type CreateStationPointInput struct {
	FromID 		uint		`json:"from"	binding:"required"`
	ToID		uint		`json:"to"	binding:"required"`
	Value		uint		`json:"value"	binding:"required"`
}

type UpdateStationPointInput struct {
	Value		uint		`json:"value"`
}

func GetStationPoints(c *gin.Context) {
	var stationpoints []models.StationPoint
	models.DB.Find(&stationpoints)
	c.JSON(http.StatusOK, gin.H{"data": stationpoints})
}

func GetStationPoint(c *gin.Context) {
	// Get model if exist
	var stationpoint models.StationPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&stationpoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stationpoint})
}

func PostStationPoint(c *gin.Context) {
	// Validate input
	var input CreateStationPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create stationpoint
	stationpoint := models.StationPoint{FromID: input.FromID, ToID: input.ToID, Value: input.Value}
	models.DB.Create(&stationpoint)
	c.JSON(http.StatusOK, gin.H{"data": stationpoint})
}

func PatchStationPoint(c *gin.Context) {
	// Get model if exist
	var stationpoint models.StationPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&stationpoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var input UpdateStationPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&stationpoint).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": stationpoint})
}

func DeleteStationPoint(c *gin.Context) {
	// Get model if exist
	var stationpoint models.StationPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&stationpoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&stationpoint)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
