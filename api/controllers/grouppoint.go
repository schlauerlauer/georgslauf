package controllers

import (
	"georgslauf/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetGroupPoints(c *gin.Context) {
	var grouppoints []models.GroupPoint
	_start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	_end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&grouppoints)
	if result.Error != nil {
		c.AbortWithStatus(500)
		slog.Warn("Get grouppoints failed.")
	}
	c.JSON(http.StatusOK, grouppoints)
}

func GetGroupPoint(c *gin.Context) {
	var grouppoint models.GroupPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&grouppoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		slog.Warn("Get grouppoint failed.")
		return
	}
	c.JSON(http.StatusOK, grouppoint)
}

func PutGroupPointByStationID(c *gin.Context) {
	var input models.PutPoint
	if err := c.Bind(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if input.Value > 100 {
		input.Value = 100
	}
	if input.Value < 0 {
		input.Value = 0
	}

	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	stationID, err := strconv.ParseInt(c.GetString("station"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	grouppoint := models.GroupPoint{
		StationID: stationID,
		GroupID:   groupID,
		Value:     input.Value,
	}

	models.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "station_id"}, {Name: "group_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&grouppoint)

	c.HTML(http.StatusOK, "station/putpoint", input)
}
