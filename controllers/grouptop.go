package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

func GetGroupTops(c *gin.Context) {
	var grouptops []models.GroupTop
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&grouptops)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get grouptops failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalGroupTop, 10))
		c.JSON(http.StatusOK, grouptops)
	}
}

func GetGroupTop(c *gin.Context) {
	var grouptop models.GroupTop
	if err := models.DB.Where("id = ?", c.Param("id")).First(&grouptop).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get grouptop failed.")
		return
	}
	c.JSON(http.StatusOK, grouptop)
}
