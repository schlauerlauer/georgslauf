package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
)

type CreateGroupPointInput struct {
	FromID 		uint		`json:"from"	binding:"required"`
	ToID		uint		`json:"to"	binding:"required"`
	Value		uint		`json:"value"	binding:"required"`
}

type UpdateGroupPointInput struct {
	Value		uint		`json:"value"`
}

func GetGroupPoints(c *gin.Context) {
	var grouppoints []models.GroupPoint
	models.DB.Find(&grouppoints)
	c.JSON(http.StatusOK, grouppoints)
}

func GetGroupPoint(c *gin.Context) {
	// Get model if exist
	var grouppoint models.GroupPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&grouppoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, grouppoint)
}

func PostGroupPoint(c *gin.Context) {
	// Validate input
	var input CreateGroupPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create grouppoint
	grouppoint := models.GroupPoint{FromID: input.FromID, ToID: input.ToID, Value: input.Value}
	models.DB.Create(&grouppoint)
	c.JSON(http.StatusOK, grouppoint)
}

func PutGroupPoint(c *gin.Context) {
	// Validate input
	var input models.GroupPoint
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// TODO log error
		return
	}
	// Put GroupPoint
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchGroupPoint(c *gin.Context) {
	// Get model if exist
	var grouppoint models.GroupPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&grouppoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var input UpdateGroupPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&grouppoint).Updates(input)
	c.JSON(http.StatusOK, grouppoint)
}

func DeleteGroupPoint(c *gin.Context) {
	// Get model if exist
	var grouppoint models.GroupPoint
	if err := models.DB.Where("id = ?", c.Param("id")).First(&grouppoint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&grouppoint)
	c.JSON(http.StatusOK, true)
}
