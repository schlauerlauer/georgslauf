package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "georgslauf/models"
    "strconv"
    log "github.com/sirupsen/logrus"
)

type CreateGroupPointInput struct {
    StationID   uint    `json:"StationID" binding:"required"`
    GroupID     uint    `json:"GroupID" binding:"required"`
    Value       uint    `json:"value" binding:"required"`
}

type UpdateGroupPointInput struct {
    Value   uint    `json:"value"`
}

func GetGroupTops(c *gin.Context) {
    var grouptops []models.GroupTop
    _start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
    _end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
    _sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
    result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&grouptops)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get grouptops failed.")
    } else {
        c.Header("X-Total-Count", strconv.FormatInt(totalGroupTop, 10))
        c.JSON(http.StatusOK, grouptops)
    }
}

func GetGroupTop(c *gin.Context) {
    var grouptop models.GroupTop
    if err := models.DB.Where("id = ?", c.Param("id")).First(&grouptop).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Get grouptop failed.")
        return
    }
    c.JSON(http.StatusOK, grouptop)
}

func GetGroupPoints(c *gin.Context) {
    var grouppoints []models.GroupPoint
    _start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
    _end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
    _sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
    result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&grouppoints)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get grouppoints failed.")
    } else {
        c.Header("X-Total-Count", strconv.FormatInt(totalGroupPoint, 10))
        c.JSON(http.StatusOK, grouppoints)
    }
}

func GetGroupPoint(c *gin.Context) {
    var grouppoint models.GroupPoint
    if err := models.DB.Where("id = ?", c.Param("id")).First(&grouppoint).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Get grouppoint failed.")
        return
    }
    c.JSON(http.StatusOK, grouppoint)
}

func PostGroupPoint(c *gin.Context) {
    // Validate input
    var input CreateGroupPointInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Post grouppoint failed.")
        return
    }
    // Create grouppoint
    grouppoint := models.GroupPoint{
        StationID: input.StationID,
        GroupID: input.GroupID,
        Value: input.Value}
    models.DB.Create(&grouppoint)
    c.JSON(http.StatusOK, grouppoint)
    totalGroupPoint+=1
}

func PutGroupPoint(c *gin.Context) {
    // Validate input
    var input models.GroupPoint
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Put grouppoint failed.")
        return
    }
    // Put GroupPoint
    models.DB.Save(&input)
    c.JSON(http.StatusOK, input)
}

func PatchGroupPoint(c *gin.Context) {
    // Get model if exist
    var grouppoint models.GroupPoint
    if err := models.DB.Where("id = ?", c.Param("id")).First(&grouppoint).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Patch grouppoint failed.")
        return
    }
    // Validate input
    var input UpdateGroupPointInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Patch grouppoint failed.")
        return
    }
    models.DB.Model(&grouppoint).Updates(input)
    c.JSON(http.StatusOK, grouppoint)
}

func DeleteGroupPoint(c *gin.Context) {
    // Get model if exist
    var grouppoint models.GroupPoint
    if err := models.DB.Where("id = ?", c.Param("id")).First(&grouppoint).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Delete grouppoint failed.")
        return
    }
    models.DB.Delete(&grouppoint)
    c.JSON(http.StatusOK, true)
}
