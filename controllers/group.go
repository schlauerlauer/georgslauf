package controllers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type CreateGroupInput struct {
	Short 		string		`json:"short"		binding:"required"`
	Name		string		`json:"name"		binding:"required"`
	Size		uint		`json:"size"		binding:"required"`
	GroupingID	uint		`json:"GroupingID"	binding:"required"`
	TribeID		uint		`json:"TribeID"		binding:"required"`
}

type UpdateGroupInput struct {
	Name		string	`json:"name"`
	Size		uint	`json:"size"`
	Short 		string	`json:"short"`
	TribeID		uint	`json:"TribeID"`
	GroupingID	uint	`json:"GroupingID"`
}

func GetGroups(c *gin.Context) {
	var groups []models.Group
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&groups)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get groups failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		fmt.Println(totalGroup)
		c.Header("X-Total-Count", strconv.FormatInt(totalGroup, 10))
		c.JSON(http.StatusOK, groups)
	}
}

func GetGroup(c *gin.Context) {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get group failed.")
		return
	}
	c.JSON(http.StatusOK, group)
}

func PostGroup(c *gin.Context) {
	// Validate input
	var input CreateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post group failed.")
		fmt.Println(err)
		return
	}
	group := models.Group{
		Short:		input.Short,
		Name:		input.Name,
		Size:		input.Size,
		GroupingID:	input.GroupingID,
		TribeID:	input.TribeID,
	}
	// Create group
	models.DB.Create(&group)
	c.JSON(http.StatusOK, group)
	totalGroup+=1
}

func PutGroup(c *gin.Context) {
	// Validate input
	var input models.Group
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Put group failed.")
		fmt.Println(err)
		return
	}
	// Put Group
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchGroup(c *gin.Context) {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Patch group failed.")
		return
	}
	// Validate input
	var input UpdateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Patch group failed.")
		return
	}
	models.DB.Model(&group).Updates(input)
	c.JSON(http.StatusOK, group)
}

func DeleteGroup(c *gin.Context) {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Delete group failed.")
		return
	}
	models.DB.Delete(&group)
	c.JSON(http.StatusOK, true)
}
