package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type CreateContentInput struct {
	Sort			uint	`json:"sort"			binding:"required"`
	Value			string	`json:"value"			binding:"required"`
	RunID			uint	`json:"RunID"			binding:"required"`
	ContenttypeID	uint	`json:"ContenttypeID"	binding:"required"`
}

type UpdateContentInput struct {
	Sort			uint	`json:"sort"`
	Value			string	`json:"value"`
	RunID			uint	`json:"RunID"`
	ContenttypeID	uint	`json:"ContenttypeID"`
}

func GetContents(c *gin.Context) {
	var contents []models.Content
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&contents)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get contents failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalContent, 10))
		c.JSON(http.StatusOK, contents)
	}
}

func GetPublicContent(c *gin.Context) {
	var public_content []models.PublicContent
	result := models.DB.Where("ct = ?", c.Param("ct")).Find(&public_content)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get public content failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", "10") //FIXME
		c.JSON(http.StatusOK, public_content)
	}
}

func GetContent(c *gin.Context) {
	// Get model if exist
	var content models.Content
	if err := models.DB.Where("id = ?", c.Param("id")).First(&content).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get content failed.")
		return
	}
	c.JSON(http.StatusOK, content)
}

func PostContent(c *gin.Context) {
	// Validate input
	var input CreateContentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post content failed.")
		return
	}
	// Create content
	content := models.Content{
		Sort: input.Sort,
		Value: input.Value,
		RunID: input.RunID,
		ContenttypeID: input.ContenttypeID,
	}
	models.DB.Create(&content)
	c.JSON(http.StatusOK, content)
	totalContent+=1
}

func PutContent(c *gin.Context) {
	// Validate input
	var input models.Content
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Put content failed.")
		return
	}
	// Put Tribe
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchContent(c *gin.Context) {
	// Get model if exist
	var content models.Content
	if err := models.DB.Where("id = ?", c.Param("id")).First(&content).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Patch content failed.")
		return
	}
	// Validate input
	var input UpdateContentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Patch content failed.")
		return
	}
	models.DB.Model(&content).Updates(input)
	c.JSON(http.StatusOK, content)
}

func DeleteContent(c *gin.Context) {
	// Get model if exist
	var content models.Content
	if err := models.DB.Where("id = ?", c.Param("id")).First(&content).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Delete content failed.")
		return
	}
	models.DB.Delete(&content)
	c.JSON(http.StatusOK, true)
}
