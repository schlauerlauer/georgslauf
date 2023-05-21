package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	log "github.com/sirupsen/logrus"
)

func GetPublic(c *gin.Context) {
	var stations []models.Station
	stationResult := models.DB.Joins("Tribe").Find(&stations)
	if stationResult.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get public info failed.")
	}
	c.HTML(http.StatusOK, "public", gin.H{
		"stations": stations,
	})
}
