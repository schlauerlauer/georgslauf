package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type CreateRuleInput struct {
	LoginID		uint		`json:"LoginID"	binding:"required"`
	Object		string		`json:"object"	binding:"required"`
	Action		[]string	`json:"action"	binding:"required"`
}

// type UpdateRuleInput struct {
// 	Ptype	string	`json:"ptype"`
// 	V0		string	`json:"v0"`
// 	V1		string	`json:"v1"`
// 	V2		string	`json:"v2"`
// }

func GetRules(c *gin.Context) {
	var rules []models.Rule
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&rules)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get rules failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalRule, 10))
		c.JSON(http.StatusOK, rules)
	}
}

// func GetRule(c *gin.Context) {
// 	// Get model if exist
// 	var rule models.Rule
// 	if err := models.DB.Where("id = ?", c.Param("id")).First(&rule).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
// 		log.Warn("Get rule failed.")
// 		return
// 	}
// 	c.JSON(http.StatusOK, rule)
// }

func PostRule(c *gin.Context) {
	// Validate input
	var input CreateRuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post rule failed: ", err)
		return
	}
	//log.Debug(input.Test) //FIXME array to string with ()
	var login models.Login
	if err := models.DB.Where("id = ?", input.LoginID).First(&login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get login failed.")
		return
	}
	rule, _ := models.EN.AddPolicy(login.Username, input.Object, input.Action)
	c.JSON(http.StatusOK, rule)
	totalRule+=1
}

// func PutRule(c *gin.Context) {
// 	// Validate input
// 	var input models.Rule
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		log.Warn("Put rule failed.")
// 		return
// 	}
// 	// Put Tribe
// 	models.DB.Save(&input)
// 	c.JSON(http.StatusOK, input)
// }

func DeleteRule(c *gin.Context) {
	// Get model if exist
	var rule models.Rule
	if err := models.DB.Where("id = ?", c.Param("id")).First(&rule).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Delete rule failed.")
		return
	}
	models.DB.Delete(&rule)
	c.JSON(http.StatusOK, true)
}

// TODO MOVE FROM DELETE GORM TO CASBIN API