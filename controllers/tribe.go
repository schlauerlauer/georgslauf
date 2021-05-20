package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type CreateTribeInput struct {
	Name		string	`json:"name"		binding:"required"`
	Short		string	`json:"short"		binding:"required"`
	Email		string	`json:"email"		binding:"required"`
	Reset		bool	`json:"reset"		binding:"required"`
	Active		bool	`json:"active"		binding:"required"`
	Confirmed	bool	`json:"confirmed	binding:"required"`
	DPSG		string	`json:"dpsg"`
	Address		string	`json:"address"`
}

type UpdateTribeInput struct {
	Name		string	`json:"name"`
	Short		string	`json:"short"`
	Email		string	`json:"email"`
	Reset		bool	`json:"reset"`
	Active		bool	`json:"active"`
	Confirmed	bool	`json:"confirmed"`
	DPSG		string	`json:"dpsg"`
	Address		string	`json:"address"`
}

func GetTribes(c *gin.Context) {
	var tribes []models.Tribe
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&tribes)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get tribes failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalTribe, 10))
		c.JSON(http.StatusOK, tribes)
	}
}

func GetTribe(c *gin.Context) {
	// Get model if exist
	var tribe models.Tribe
	if err := models.DB.Where("id = ?", c.Param("id")).First(&tribe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get tribe failed.")
		return
	}
	c.JSON(http.StatusOK, tribe)
}

func PostTribe(c *gin.Context) {
	// Validate input
	var input CreateTribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post tribe failed.")
		return
	}
	// Create tribe
	tribe := models.Tribe{
		Name: input.Name,
		Short: input.Short,
		Email: input.Email,
		Reset: input.Reset,
		Active: input.Active,
		Confirmed: input.Confirmed,
		DPSG: input.DPSG,
		Address: input.Address,
	}
	models.DB.Create(&tribe)
	c.JSON(http.StatusOK, tribe)
	totalTribe+=1
}

func PutTribe(c *gin.Context) {
	// Validate input
	var input models.Tribe
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Put tribe failed.")
		return
	} //TODO error checking (e.g. unique error)
	// Put Tribe
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchTribe(c *gin.Context) {
	// Get model if exist
	var tribe models.Tribe
	if err := models.DB.Where("id = ?", c.Param("id")).First(&tribe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Patch tribe failed.")
		return
	}
	// Validate input
	var input UpdateTribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Patch tribe failed.")
		return
	}
	models.DB.Model(&tribe).Updates(input)
	c.JSON(http.StatusOK, tribe)
}

func DeleteTribe(c *gin.Context) {
	// Get model if exist
	var tribe models.Tribe
	if err := models.DB.Where("id = ?", c.Param("id")).First(&tribe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Delete tribe failed.")
		return
	}
	models.DB.Delete(&tribe)
	c.JSON(http.StatusOK, true)
}
