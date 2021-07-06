package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type CreateContentTypeInput struct {
	Name	string	`json:"name"	binding:"required"`
	Public	bool	`json:"public	binding:"required"`
}

type UpdateContentTypeInput struct {
	Name	string	`json:"name"`
	Public	bool	`json:"public"`
}

func GetContentTypes(c *gin.Context) {
	var contenttypes []models.ContentType
	_start, _ :=strconv.Atoi(c.DefaultQuery("_start", "0"))
	_end, _ :=strconv.Atoi(c.DefaultQuery("_end", "10"))
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&contenttypes)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get contenttypes failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalContentType, 10))
		c.JSON(http.StatusOK, contenttypes)
	}
}

func GetContentType(c *gin.Context) {
	// Get model if exist
	var contenttype models.ContentType
	if err := models.DB.Where("id = ?", c.Param("id")).First(&contenttype).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get contenttype failed.")
		return
	}
	c.JSON(http.StatusOK, contenttype)
}

func PostContentType(c *gin.Context) {
	// Validate input
	var input CreateContentTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post contenttype failed.")
		return
	}
	// Create contenttype
	contenttype := models.ContentType{
		Name: input.Name,
		Public: input.Public,
	}
	models.DB.Create(&contenttype)
	c.JSON(http.StatusOK, contenttype)
	totalContentType+=1
}

func PutContentType(c *gin.Context) {
	// Validate input
	var input models.ContentType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Put contenttype failed.")
		return
	}
	// Put Tribe
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchContentType(c *gin.Context) {
	// Get model if exist
	var contenttype models.ContentType
	if err := models.DB.Where("id = ?", c.Param("id")).First(&contenttype).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Patch contenttype failed.")
		return
	}
	// Validate input
	var input UpdateContentTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Patch contenttype failed.")
		return
	}
	models.DB.Model(&contenttype).Updates(input)
	c.JSON(http.StatusOK, contenttype)
}

func DeleteContentType(c *gin.Context) {
	// Get model if exist
	var contenttype models.ContentType
	if err := models.DB.Where("id = ?", c.Param("id")).First(&contenttype).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Delete contenttype failed.")
		return
	}
	models.DB.Delete(&contenttype)
	c.JSON(http.StatusOK, true)
}
