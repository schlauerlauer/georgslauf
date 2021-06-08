package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type CreateStationPointInput struct {
	GroupID 	uint	`json:"GroupID"		binding:"required"`
	StationID	uint	`json:"StationID"	binding:"required"`
	Value		uint	`json:"value"		binding:"required"`
}

type UpdateStationPointInput struct {
	Value		uint	`json:"value"`
}

func GetStationTops(c *gin.Context) {
	var stationtops []models.StationTop
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&stationtops)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get stationtops failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalStationTop, 10))
		c.JSON(http.StatusOK, stationtops)
	}
}

func GetStationTop(c *gin.Context) {
	var stationtop models.StationTop
	if err := models.DB.Where("id = ?", c.Param("id")).First(&stationtop).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get stationtop failed.")
		return
	}
	c.JSON(http.StatusOK, stationtop)
}

func GetStationPoints(c *gin.Context) {
	var stationpoints []models.StationPoint
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&stationpoints)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get stationpoints failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalStationPoint, 10))
		c.JSON(http.StatusOK, stationpoints)
	}
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
	totalStationPoint+=1
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
