package controllers

import (
	"errors"
	"georgslauf/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateStationInput struct {
	Name    string  `json:"name" inding:"required"`
	Short   string  `json:"short" binding:"required"`
	TribeID int64    `json:"TribeID" binding:"required"`
	Size    int64    `json:"size" binding:"required"`
}

type UpdateStationInput struct {
	Name    string  `json:"name"`
	Short   string  `json:"short"`
	TribeID int64    `json:"TribeID"`
	Size    int64    `json:"size"`
}

func GetStationsByTribe(c *gin.Context) {
	var stations []models.Station
	result := models.DB.Where("tribe_id = ?", c.Param("tribeid")).Find(&stations)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get stations failed.")
	}
	c.HTML(http.StatusOK, "station/tribe", stations)
}

func GetStationsPublic(c *gin.Context) {
	var stations []models.Station
	result := models.DB.Find(&stations)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get public stations failed.")
	}
	c.HTML(http.StatusOK, "station/public", stations)
}

func GetStations(c *gin.Context) {
	var stations []models.Station
	_start, _ :=strconv.Atoi(c.DefaultQuery("_start", "0"))
	_end, _ :=strconv.Atoi(c.DefaultQuery("_end", "10"))
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&stations)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get stations failed.")
	}
	c.JSON(http.StatusOK, stations)
}

func GetStation(c *gin.Context) models.Station {
	// Get model if exist
	var station models.Station
	if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get station failed.")
	}
	return station
}

func GetStationByID(id string) (models.Station, error) {
	var station models.Station
	if err := models.DB.Where("stations.id = ?", id).Joins("Tribe").First(&station).Error; err != nil {
		log.Warn("Get station failed.")
		return models.Station{}, errors.New("station not found")
	}
	return station, nil
}

func PostStation(c *gin.Context) {
	// Validate input
	var input CreateStationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post station failed.")
		return
	}
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
		log.Warn("Put station failed.")
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
		log.Warn("Patch station failed.")
		return
	}
	// Validate input
	var input UpdateStationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Patch station failed.")
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
		log.Warn("Delete station failed.")
		return
	}
	models.DB.Delete(&station)
	c.JSON(http.StatusOK, true)
}
