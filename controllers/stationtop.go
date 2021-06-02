package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

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
