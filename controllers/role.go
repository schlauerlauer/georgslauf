package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
)

type CreateRoleInput struct {
	Name	string	`json:"name"	binding:"required"`
}

type UpdateRoleInput struct {
	Name	string	`json:"name"`
}

func GetRoles(c *gin.Context) {
	var roles []models.Role
	result := models.DB.Find(&roles)
	if result.Error != nil {
		c.AbortWithStatus(500)
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(result.RowsAffected, 10))
		c.JSON(http.StatusOK, roles)
	}
}

func GetRole(c *gin.Context) {
	// Get model if exist
	var role models.Role
	if err := models.DB.Where("id = ?", c.Param("id")).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, role)
}

func PostRole(c *gin.Context) {
	// Validate input
	var input CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} //TODO error checking (e.g. unique error)
	// Create role
	role := models.Role{
		Name: input.Name,
	}
	models.DB.Create(&role)
	c.JSON(http.StatusOK, role)
}

func PutRole(c *gin.Context) {
	// Validate input
	var input models.Role
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// TODO log error
		return
	}
	// Put Tribe
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchRole(c *gin.Context) {
	// Get model if exist
	var role models.Role
	if err := models.DB.Where("id = ?", c.Param("id")).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&role).Updates(input)
	c.JSON(http.StatusOK, role)
}

func DeleteRole(c *gin.Context) {
	// Get model if exist
	var role models.Role
	if err := models.DB.Where("id = ?", c.Param("id")).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&role)
	c.JSON(http.StatusOK, true)
}
