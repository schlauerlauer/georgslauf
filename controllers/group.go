package controllers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
)

type CreateGroupInput struct {
	Name	string	`json:"name"	binding:"required"`
	Size	uint	`json:"size"	binding:"required"`
	TribeID	uint	`json:"TribeID"	binding:"required"`
	Short 	string	`json:"short"	binding:"required"`
	RoleID	uint	`json:"role"	binding:"required"`
	Details	string	`json:"details"	binding:"required"`
	Contact	string	`json:"contact"	binding:"required"`
}

type UpdateGroupInput struct {
	Name	string	`json:"name"`
	Size	uint	`json:"size"`
	Short 	string	`json:"short"`
	TribeID	uint	`json:"TribeID"`
	RoleID	uint	`json:"role"`
	Details	string	`json:"details"`
	Contact	string	`json:"contact"`
}

func GetGroups(c *gin.Context) {
	var groups []models.Group
	if err := models.DB.Find(&groups).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		//lastname := c.Query("_end")
		// start := c.DefaultQuery("_start", "0")
		// end := c.DefaultQuery("_end", "10")
		// fmt.Println(start)
		c.Header("X-Total-Count", "10") //FIXME
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
		return
	}
	// Create group
	group := models.Group{
		Name: input.Name,
		Size: input.Size,
		TribeID: input.TribeID,
		Short: input.Short,
		RoleID: input.RoleID,
		Details: input.Details,
		Contact: input.Contact,
	}
	models.DB.Create(&group)
	c.JSON(http.StatusOK, group)
}

func PutGroup(c *gin.Context) {
	// Validate input
	var input models.Group
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	// Put group
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
