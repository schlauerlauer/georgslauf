package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "georgslauf/models"
    "strconv"
    log "github.com/sirupsen/logrus"
)

type RuleInput struct {
    Ptype       string      `json:"ptype"   binding:"required"`
    LoginID     uint        `json:"LoginID" binding:"required"`
    Object      string      `json:"object"  binding:"required"`
    Actions     []string    `json:"action"  binding:"required"`
}

// type UpdateRuleInput struct {
//  Ptype   string  `json:"ptype"`
//  V0      string  `json:"v0"`
//  V1      string  `json:"v1"`
//  V2      string  `json:"v2"`
// }

func GetRules(c *gin.Context) {
    var rules []models.Rule
    _start, _ :=strconv.Atoi(c.DefaultQuery("_start", "0"))
    _end, _ :=strconv.Atoi(c.DefaultQuery("_end", "10"))
    _sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
    result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&rules)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get rules failed.")
    } else {
        c.Header("X-Total-Count", strconv.FormatInt(totalRule, 10))
        c.JSON(http.StatusOK, rules)
    }
}

func PostRule(c *gin.Context) {
    // Validate input
    var input RuleInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Post rule failed: ", err)
        return
    }
    var login models.Login
    if err := models.DB.Where("id = ?", input.LoginID).First(&login).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Get login failed.")
        return
    }
    action := ""
    if len(input.Actions) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No action defined!"})
        log.Warn("Empty action posted.")
        return
    }
    if len(input.Actions) == 1 {
        action = input.Actions[0]
    } else {
        for i := 0; i < len(input.Actions); i++ {
            action += "(" + input.Actions[i] + ")"
            if i < len(input.Actions) -1 {
                action += "|"
            }
        }
    }
    rule, _ := models.EN.AddPolicy(login.Username, input.Object, action)
    c.JSON(http.StatusOK, rule)
    totalRule+=1
}

func DeleteRule(c *gin.Context) {
    var rule models.Rule
    if err := models.DB.Where("id = ?", c.Param("id")).First(&rule).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Delete rule failed.")
        return
    }
    removed, err := models.EN.RemovePolicy(rule.V0, rule.V1, rule.V2)
    if err != nil {
        log.Error("Error removing policy: ", err)
    }
    c.JSON(http.StatusOK, removed)
}
