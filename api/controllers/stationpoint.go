package controllers

import (
	"georgslauf/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateStationPointInput struct {
	GroupID     int64    `json:"GroupID" binding:"required"`
	StationID   int64    `json:"StationID" binding:"required"`
	Value       int64    `json:"value" binding:"required"`
}

type UpdateStationPointInput struct {
	Value   int64    `json:"value"`
}

func GetStationPoints(c *gin.Context) {
	var stationpoints []models.StationPoint
	_start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	_end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&stationpoints)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get stationpoints failed.")
	}
	c.JSON(http.StatusOK, stationpoints)
}

func GetStationPoint(c *gin.Context) {
	// Get model if exist
	var stationpoint models.StationPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&stationpoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get stationpoint failed.")
		return
	}
	c.JSON(http.StatusOK, stationpoint)
}

func PostStationPoint(c *gin.Context) {
	// Validate input
	var input CreateStationPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post stationpoint failed.")
		return
	}
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
		log.Warn("Put stationpoint failed.")
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
		log.Warn("Patch stationpoint failed.")
		return
	}
	// Validate input
	var input UpdateStationPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Patch stationpoint failed.")
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
		log.Warn("Delete stationpoint failed.")
		return
	}
	models.DB.Delete(&stationpoint)
	c.JSON(http.StatusOK, true)
}
