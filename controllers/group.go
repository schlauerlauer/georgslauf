package controllers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
)

type CreateGroupInput struct {
	Short 		string		`json:"short"	binding:"required"`
	Name		string		`json:"name"	binding:"required"`
	Size		uint		`json:"size"	binding:"required"`
	RoleID		uint		`json:"RoleID"	binding:"required"`
	TribeID		uint		`json:"TribeID"	binding:"required"`
	Details		string		`json:"details"	binding:"required"`
	Contact		string		`json:"contact"	binding:"required"`
}

type UpdateGroupInput struct {
	Name	string	`json:"name"`
	Size	uint	`json:"size"`
	Short 	string	`json:"short"`
	TribeID	uint	`json:"TribeID"`
	RoleID	uint	`json:"RoleID"`
	Details	string	`json:"details"`
	Contact	string	`json:"contact"`
}

// TODO move total to redis

func GetGroups(c *gin.Context) {
	var groups []models.Group
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&groups)
	if result.Error != nil {
		c.AbortWithStatus(500)
		// TODO add logging everywhere
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		fmt.Println(totalGroup)
		c.Header("X-Total-Count", strconv.FormatInt(totalGroup, 10)) //FIXME everywhere -> only shows current page not all
		c.JSON(http.StatusOK, groups)
	}
}

func GetGroup(c *gin.Context) {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, group)
}

func PostGroup(c *gin.Context) {
	// Validate input
	var input CreateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	group := models.Group{
		Short: input.Short,
		Name: input.Name,
		Size: input.Size,
		RoleID: input.RoleID,
		TribeID: input.TribeID,
		Details: input.Details,
		Contact: input.Contact,
	} //TODO error checking (e.g. unique error)
	// Create group
	models.DB.Create(&group)
	c.JSON(http.StatusOK, "") // TODO don't send the whole json back (everywhere) (get id back! tutorial)
	totalGroup+=1
}

func PutGroup(c *gin.Context) {
	// Validate input
	var input models.Group
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	// Put Group
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input) // TODO here too
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
	c.JSON(http.StatusOK, group)
}

func DeleteGroup(c *gin.Context) {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&group)
	c.JSON(http.StatusOK, true)
}
