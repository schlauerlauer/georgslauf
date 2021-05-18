package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
)

type CreateStationPointInput struct {
	GroupID 	uint	`json:"GroupID"		binding:"required"`
	StationID	uint	`json:"StationID"	binding:"required"`
	Value		uint	`json:"value"		binding:"required"`
}

type UpdateStationPointInput struct {
	Value		uint	`json:"value"`
}

func GetStationPoints(c *gin.Context) {
	var stationpoints []models.StationPoint
	result := models.DB.Find(&stationpoints)
	if result.Error != nil {
		c.AbortWithStatus(500)
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(result.RowsAffected, 10))
		c.JSON(http.StatusOK, stationpoints)
	}
}

func GetStationPoint(c *gin.Context) {
	// Get model if exist
	var stationpoint models.StationPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&stationpoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, stationpoint)
}

func PostStationPoint(c *gin.Context) {
	// Validate input
	var input CreateStationPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} //TODO error checking (e.g. unique error)
	// Create stationpoint
	stationpoint := models.StationPoint{
		GroupID: input.GroupID,
		StationID: input.StationID,
		Value: input.Value,
	}
	models.DB.Create(&stationpoint)
	c.JSON(http.StatusOK, stationpoint)
}

func PutStationPoint(c *gin.Context) {
	// Validate input
	var input models.StationPoint
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// TODO log error
		return
	}
	// Put StationPoint
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
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
	c.JSON(http.StatusOK, stationpoint)
}

func DeleteStationPoint(c *gin.Context) {
	// Get model if exist
	var stationpoint models.StationPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&stationpoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&stationpoint)
	c.JSON(http.StatusOK, true)
}
