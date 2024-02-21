package controllers

import (
	"georgslauf/models"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPublic(c *gin.Context) {
	var stations []models.Station
	stationResult := models.DB.Joins("Tribe").Find(&stations)
	if stationResult.Error != nil {
		c.AbortWithStatus(500)
		slog.Warn("Get public info failed.")
	}
	c.HTML(http.StatusOK, "public", gin.H{
		"stations": stations,
	})
}
