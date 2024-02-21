package controllers

import (
	"georgslauf/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateTribeInput struct {
	Name    string `json:"name" binding:"required"`
	Short   string `json:"short" binding:"required"`
	DPSG    string `json:"dpsg" binding:"required"`
	Address string `json:"address" binding:"required"`
	URL     string `json:"url"`
}

type UpdateTribeInput struct {
	Name    string `json:"name"`
	Short   string `json:"short"`
	DPSG    string `json:"dpsg"`
	Address string `json:"address"`
	URL     string `json:"url"`
}

func GetTribes(c *gin.Context) {
	var tribes []models.Tribe
	_start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	_end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&tribes)
	if result.Error != nil {
		c.AbortWithStatus(500)
		slog.Warn("Get tribes failed.")
	}
	c.JSON(http.StatusOK, tribes)
}

func GetTribe(c *gin.Context) {
	// Get model if exist
	var tribe models.Tribe
	if err := models.DB.Where("id = ?", c.Param("id")).First(&tribe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		slog.Warn("Get tribe failed.")
		return
	}
	c.JSON(http.StatusOK, tribe)
}

func PostTribe(c *gin.Context) {
	// Validate input
	var input CreateTribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Warn("Post tribe failed.")
		return
	}
	// Create tribe
	tribe := models.Tribe{
		Name:    input.Name,
		Short:   input.Short,
		DPSG:    input.DPSG,
		Address: input.Address,
	}
	models.DB.Create(&tribe)
	c.JSON(http.StatusOK, tribe)
}

func PutTribe(c *gin.Context) {
	// Validate input
	var input models.Tribe
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Warn("Put tribe failed.")
		return
	}
	// Put Tribe
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchTribe(c *gin.Context) {
	// Get model if exist
	var tribe models.Tribe
	if err := models.DB.Where("id = ?", c.Param("id")).First(&tribe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		slog.Warn("Patch tribe failed.")
		return
	}
	// Validate input
	var input UpdateTribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Warn("Patch tribe failed.")
		return
	}
	models.DB.Model(&tribe).Updates(input)
	c.JSON(http.StatusOK, tribe)
}
