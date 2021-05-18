package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
)

type CreateLoginInput struct {
	Name 		string		`json:"name"		binding:"required"`
	Password	string		`json:"password"	binding:"required"`
	RoleID		uint		`json:"RoleID"		binding:"required"`
}

type UpdateLoginInput struct {
	Name		string		`json:"name"`
	Password	string		`json:"password"`
	RoleID		uint		`json:"RoleID"`
}

func GetLogins(c *gin.Context) {
	var logins []models.Login
	result := models.DB.Find(&logins)
	if result.Error != nil {
		c.AbortWithStatus(500)
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(result.RowsAffected, 10))
		c.JSON(http.StatusOK, logins)
	}
}

func GetLogin(c *gin.Context) {
	// Get model if exist
	var login models.Login
	if err := models.DB.Where("id = ?", c.Param("id")).First(&login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, login)
}

func PostLogin(c *gin.Context) {
	// Validate input
	var input CreateLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} // TODO error checking multiple unique
	// Create login
	login := models.Login{
		Name: input.Name,
		Password: input.Password,
		RoleID: input.RoleID,
	}
	models.DB.Create(&login)
	c.JSON(http.StatusOK, login)
}

func PutLogin(c *gin.Context) {
	// Validate input
	var input models.Login
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// TODO log error
		return
	}
	// Put Login
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchLogin(c *gin.Context) {
	// Get model if exist
	var login models.Login
	if err := models.DB.Where("id = ?", c.Param("id")).First(&login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var input UpdateLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&login).Updates(input)
	c.JSON(http.StatusOK, login)
}

func DeleteLogin(c *gin.Context) {
	// Get model if exist
	var login models.Login
	if err := models.DB.Where("id = ?", c.Param("id")).First(&login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&login)
	c.JSON(http.StatusOK, true)
}
