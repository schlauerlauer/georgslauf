package controllers

import (
	"errors"
	"georgslauf/models"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetConfigGin(c *gin.Context) {
	// Get model if exist
	var config models.Config
	if err := models.DB.First(&config).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		slog.Warn("Get config failed.")
		return
	}
	c.JSON(http.StatusOK, config)
}

func GetConfig() (models.Config, error) {
	var cfg models.Config
	if err := models.DB.First(&cfg).Error; err != nil {
		slog.Warn("System config not found.")
		return models.Config{}, errors.New("system config not found")
	}
	return cfg, nil
}

func PutConfig(c *gin.Context) {
	// Validate input
	var input models.Config
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Warn("Put config failed.")
		return
	}
	// Put Tribe
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchConfig(c *gin.Context) {
	// Get model if exist
	var config models.Config
	if err := models.DB.Where("id = ?", c.Param("id")).First(&config).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		slog.Warn("Patch config failed.")
		return
	}
	// Validate input
	var input models.Config
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Warn("Patch config failed.")
		return
	}
	models.DB.Model(&config).Updates(input)
	c.JSON(http.StatusOK, config)
}
