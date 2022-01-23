package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "georgslauf/models"
    "strconv"
    log "github.com/sirupsen/logrus"
)

type CreateRunInput struct {
    Year        uint        `json:"year" binding:"required"`
    Note        string      `json:"note" binding:"required"`
    TribeID     uint        `json:"TribeID" binding:"required"`
}

type UpdateRunInput struct {
    Year        uint    `json:"year"`
    Note        string  `json:"note"`
    TribeID     uint    `json:"TribeID"`
}

func GetRuns(c *gin.Context) {
    var runs []models.Run
    _id := c.Query("id")
    result := models.DB
    if _id != "" {
        result = models.DB.Where("id = ?", _id).Find(&runs)
    } else {
        _start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
        _end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
        _sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
        result = models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&runs)
    }
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get runs failed.")
    } else {
        c.Header("X-Total-Count", strconv.FormatInt(totalRun, 10))
        c.JSON(http.StatusOK, runs)
    }
}

func GetRun(c *gin.Context) {
    // Get model if exist
    var run models.Run
    if err := models.DB.Where("id = ?", c.Param("id")).First(&run).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Get run failed.")
        return
    }
    c.JSON(http.StatusOK, run)
}

func PostRun(c *gin.Context) {
    // Validate input
    var input CreateRunInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Post run failed.")
        return
    }
    // Create run
    run := models.Run{
        Year: input.Year,
        Note: input.Note,
        TribeID: input.TribeID,
    }
    models.DB.Create(&run)
    c.JSON(http.StatusOK, run)
    totalRun+=1
}

func PutRun(c *gin.Context) {
    // Validate input
    var input models.Run
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Put run failed.")
        return
    }
    // Put Run
    models.DB.Save(&input)
    c.JSON(http.StatusOK, input)
}

func PatchRun(c *gin.Context) {
    // Get model if exist
    var run models.Run
    if err := models.DB.Where("id = ?", c.Param("id")).First(&run).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Patch run failed.")
        return
    }
    // Validate input
    var input UpdateRunInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Patch run failed.")
        return
    }
    models.DB.Model(&run).Updates(input)
    c.JSON(http.StatusOK, run)
}

func DeleteRun(c *gin.Context) {
    // Get model if exist
    var run models.Run
    if err := models.DB.Where("id = ?", c.Param("id")).First(&run).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Delete run failed.")
        return
    }
    models.DB.Delete(&run)
    c.JSON(http.StatusOK, true)
}
