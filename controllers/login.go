package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type CreateLoginInput struct {
	Username 	string		`json:"username"	binding:"required"`
	Password	string		`json:"password"	binding:"required"`
	RoleID		uint		`json:"RoleID"		binding:"required"`
	Salt		string		`json:"salt"		binding:"required"`
	Reset		bool		`json:"reset"		binding:"required"`
	Active		bool		`json:"active"		binding:"required"`
	Confirmed	bool		`json:"confirmed"	binding:"required"`
	Phone		string		`json:"phone"		binding:"required"`
	Email		string		`json:"email"		binding:"required"`
	Contact		string		`json:"contact"		binding:"required"`
}

type UpdateLoginInput struct {
	Username	string		`json:"username"`
	Password	string		`json:"password"`
	RoleID		uint		`json:"RoleID"`
	Reset		bool		`json:"reset"`
	Active		bool		`json:"active"`
	Confirmed	bool		`json:"confirmed"`
	Phone		string		`json:"phone"`
	Email		string		`json:"email"`
	Contact		string		`json:"contact"`
}

func GetLogins(c *gin.Context) {
	var logins []models.Login
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&logins)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get logins failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalLogin, 10))
		c.JSON(http.StatusOK, logins)
	}
}

func GetLogin(c *gin.Context) {
	// Get model if exist
	var login models.Login
	if err := models.DB.Where("id = ?", c.Param("id")).First(&login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get login failed.")
		return
	}
	c.JSON(http.StatusOK, login)
}

func PostLogin(c *gin.Context) {
	// Validate input
	var input CreateLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post login failed.")
		return
	}
	// Create login
	login := models.Login{
		Username: input.Username,
		Password: input.Password,
		RoleID: input.RoleID,
		Salt: input.Salt,
		Reset: input.Reset,
		Active: input.Active,
		Confirmed: input.Confirmed,
		Phone: input.Phone,
		Email: input.Email,
		Contact: input.Contact,
	}
	models.DB.Create(&login)
	c.JSON(http.StatusOK, login)
	totalLogin+=1
}

func PutLogin(c *gin.Context) {
	// Validate input
	var input models.Login
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Put login failed.")
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
		log.Warn("Patch login failed.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var input UpdateLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn("Past login failed.")
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
		log.Warn("Delete login failed.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&login)
	c.JSON(http.StatusOK, true)
}
