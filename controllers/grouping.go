package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "georgslauf/models"
    "strconv"
    log "github.com/sirupsen/logrus"
)

type CreateGroupingInput struct {
    Name    string  `json:"name" binding:"required"`
    Short   string  `json:"short" binding:"required"`
}

type UpdateGroupingInput struct {
    Name    string  `json:"name"`
    Short   string  `json:"short"`
}

func GetGroupings(c *gin.Context) {
    var groupings []models.Grouping
    _start, _ :=strconv.Atoi(c.DefaultQuery("_start", "0"))
    _end, _ :=strconv.Atoi(c.DefaultQuery("_end", "10"))
    _sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
    result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&groupings)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get groupings failed.")
    } else {
        c.Header("X-Total-Count", strconv.FormatInt(totalGrouping, 10))
        c.JSON(http.StatusOK, groupings)
    }
}

func GetGrouping(c *gin.Context) {
    // Get model if exist
    var grouping models.Grouping
    if err := models.DB.Where("id = ?", c.Param("id")).First(&grouping).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Get grouping failed.")
        return
    }
    c.JSON(http.StatusOK, grouping)
}

func PostGrouping(c *gin.Context) {
    // Validate input
    var input CreateGroupingInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Post grouping failed.")
        return
    }
    // Create grouping
    grouping := models.Grouping{
        Name: input.Name,
        Short: input.Short,
    }
    models.DB.Create(&grouping)
    c.JSON(http.StatusOK, grouping)
    totalGrouping+=1
}

func PutGrouping(c *gin.Context) {
    // Validate input
    var input models.Grouping
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Put grouping failed.")
        return
    }
    // Put Tribe
    models.DB.Save(&input)
    c.JSON(http.StatusOK, input)
}

func PatchGrouping(c *gin.Context) {
    // Get model if exist
    var grouping models.Grouping
    if err := models.DB.Where("id = ?", c.Param("id")).First(&grouping).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Patch grouping failed.")
        return
    }
    // Validate input
    var input UpdateGroupingInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Patch grouping failed.")
        return
    }
    models.DB.Model(&grouping).Updates(input)
    c.JSON(http.StatusOK, grouping)
}

func DeleteGrouping(c *gin.Context) {
    // Get model if exist
    var grouping models.Grouping
    if err := models.DB.Where("id = ?", c.Param("id")).First(&grouping).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Delete grouping failed.")
        return
    }
    models.DB.Delete(&grouping)
    c.JSON(http.StatusOK, true)
}
