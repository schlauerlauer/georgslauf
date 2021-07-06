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
	DPSG		string	`json:"dpsg"		binding:"required"`
	Address		string	`json:"address"		binding:"required"`
	LoginID		uint	`json:"LoginID"		binding:"required"`
}

type UpdateTribeInput struct {
	Name		string	`json:"name"`
	Short		string	`json:"short"`
	DPSG		string	`json:"dpsg"`
	Address		string	`json:"address"`
	LoginID		uint	`json:"LoginID`
}

func GetTribes(c *gin.Context) {
	var tribes []models.Tribe
	_start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	_end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&tribes)
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
		DPSG: input.DPSG,
		Address: input.Address,
		LoginID: input.LoginID,
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
	}
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
