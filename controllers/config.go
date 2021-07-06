package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type CreateConfigInput struct {
	Name	string	`json:"name"	binding:"required"`
	ValueB	bool	`json:"valueb"	binding:"required"`
}

type UpdateConfigInput struct {
	Name	string	`json:"name"`
	ValueB	bool	`json:"valueb"`
}

func GetConfigs(c *gin.Context) {
	var configs []models.Config
	_start, _ :=strconv.Atoi(c.DefaultQuery("_start", "0"))
	_end, _ :=strconv.Atoi(c.DefaultQuery("_end", "10"))
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&configs)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get configs failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalConfig, 10))
		c.JSON(http.StatusOK, configs)
	}
}

func GetConfig(c *gin.Context) {
	// Get model if exist
	var config models.Config
	if err := models.DB.Where("id = ?", c.Param("id")).First(&config).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get config failed.")
		return
	}
	c.JSON(http.StatusOK, config)
}

func PostConfig(c *gin.Context) {
	// Validate input
	var input CreateConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post config failed.")
		return
	}
	// Create config
	config := models.Config{
		Name: input.Name,
		ValueB: input.ValueB,
	}
	models.DB.Create(&config)
	c.JSON(http.StatusOK, config)
	totalConfig+=1
}

func PutConfig(c *gin.Context) {
	// Validate input
	var input models.Config
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Put config failed.")
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
		log.Warn("Patch config failed.")
		return
	}
	// Validate input
	var input UpdateConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Patch config failed.")
		return
	}
	models.DB.Model(&config).Updates(input)
	c.JSON(http.StatusOK, config)
}

func DeleteConfig(c *gin.Context) {
	// Get model if exist
	var config models.Config
	if err := models.DB.Where("id = ?", c.Param("id")).First(&config).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Delete config failed.")
		return
	}
	models.DB.Delete(&config)
	c.JSON(http.StatusOK, true)
}
