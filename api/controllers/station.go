package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "georgslauf/models"
    "strconv"
    log "github.com/sirupsen/logrus"
)

type CreateStationInput struct {
    Name    string  `json:"name" inding:"required"`
    Short   string  `json:"short" binding:"required"`
    TribeID uint    `json:"TribeID" binding:"required"`
    Size    uint    `json:"size" binding:"required"`
    // LoginID uint    `json:"LoginID" binding:"required"`
}

type UpdateStationInput struct {
    Name    string  `json:"name"`
    Short   string  `json:"short"`
    TribeID uint    `json:"TribeID"`
    Size    uint    `json:"size"`
    // LoginID uint    `json:"LoginID"`
}

func GetStationsByTribe(c *gin.Context) {
    var stations []models.Station
    result := models.DB.Where("tribe_id = ?", c.Param("tribeid")).Find(&stations)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get stations failed.")
    } else {
        c.HTML(http.StatusOK, "station/tribe", stations)
    }
}

func GetStationsPublic(c *gin.Context) {
    var stations []models.Station
    result := models.DB.Find(&stations)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get public stations failed.")
    } else {
        c.HTML(http.StatusOK, "station/public", stations)
    }
}

func GetStations(c *gin.Context) {
    var stations []models.Station
    _start, _ :=strconv.Atoi(c.DefaultQuery("_start", "0"))
    _end, _ :=strconv.Atoi(c.DefaultQuery("_end", "10"))
    _sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
    result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&stations)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get stations failed.")
    } else {
        c.JSON(http.StatusOK, stations)
    }
}

func GetStation(c *gin.Context) {
    // Get model if exist
    var station models.Station
    if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Get station failed.")
        return
    }
    c.JSON(http.StatusOK, station)
}

func PostStation(c *gin.Context) {
    // Validate input
    var input CreateStationInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Post station failed.")
        return
    }
    // Create station
    station := models.Station{
        Name: input.Name,
        Short: input.Short,
        TribeID: input.TribeID,
        Size: input.Size,
        // LoginID: input.LoginID,
    }
    models.DB.Create(&station)
    c.JSON(http.StatusOK, station)
}

func PutStation(c *gin.Context) {
    // Validate input
    var input models.Station
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Put station failed.")
        return
    }
    // Put Tribe
    models.DB.Save(&input)
    c.JSON(http.StatusOK, input)
}

func PatchStation(c *gin.Context) {
    // Get model if exist
    var station models.Station
    if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Patch station failed.")
        return
    }
    // Validate input
    var input UpdateStationInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Patch station failed.")
        return
    }
    models.DB.Model(&station).Updates(input)
    c.JSON(http.StatusOK, station)
}

func DeleteStation(c *gin.Context) {
    // Get model if exist
    var station models.Station
    if err := models.DB.Where("id = ?", c.Param("id")).First(&station).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Delete station failed.")
        return
    }
    models.DB.Delete(&station)
    c.JSON(http.StatusOK, true)
}
