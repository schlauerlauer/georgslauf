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

func GetGroups(c *gin.Context) {
	var groups []models.Group
	result := models.DB.Find(&groups)
	if result.Error != nil {
		c.AbortWithStatus(404) // TODO do this everywhere
		fmt.Println(result.Error)
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		//lastname := c.Query("_end") // TODO
		// start := c.DefaultQuery("_start", "0")
		// end := c.DefaultQuery("_end", "10")
		// fmt.Println(start)
		c.Header("X-Total-Count", strconv.FormatInt(result.RowsAffected, 10))
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
	}
	// Create group
	models.DB.Create(&group)
	c.JSON(http.StatusOK, group) // TODO don't send the whole json back
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
	fmt.Println(&input)
	c.JSON(http.StatusOK, input)
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
