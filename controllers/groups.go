package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
)

type CreateGroupInput struct {
	Name	string	`json:"name"	binding:"required"`
	Size	uint	`json:"size"	binding:"required"`
}

type UpdateGroupInput struct {
	Name	string	`json:"name"`
	Size	uint	`json:"size"`
}

func GetGroups(c *gin.Context) {
	var groups []models.Group
	models.DB.Find(&groups)
	c.JSON(http.StatusOK, gin.H{"data": groups})
}

func GetGroup(c *gin.Context) {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": group})
}

func PostGroup(c *gin.Context) {
	// Validate input
	var input CreateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create group
	group := models.Group{Name: input.Name, Size: input.Size}
	models.DB.Create(&group)
	c.JSON(http.StatusOK, gin.H{"data": group})
}

func PatchGroup(c *gin.Context) {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var input UpdateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&group).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": group})
}

func DeleteGroup(c *gin.Context) {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&group)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
