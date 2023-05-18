package controllers

import (
	"fmt"
	"georgslauf/models"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateGroupInput struct {
    Short       string  `json:"short" binding:"required"`
    Name        string  `json:"name" binding:"required"`
    Size        uint    `json:"size" binding:"required"`
    GroupingID  uint    `json:"groupingID" binding:"required"`
    TribeID     uint    `json:"TribeID" binding:"required"`
}

type UpdateGroupInput struct {
    Name        string  `json:"name"`
    Size        uint    `json:"size"`
    Short       string  `json:"short"`
    TribeID     uint    `json:"TribeID"`
// GROuping?
}

func GetGroupsByTribe(c *gin.Context) {
    var groups []models.Group
    result := models.DB.Where("tribe_id = ?", c.Param("tribeid")).Find(&groups)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get stations failed.")
    } else {
        c.HTML(http.StatusOK, "group/tribe", groups)
    }
}

func GetGroupsPublic(c *gin.Context) {
    var groups []models.Group
    _start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
    _end, _ := strconv.Atoi(c.DefaultQuery("_end", "10"))
    _sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
    result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&groups)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get public groups failed.")
    } else {
        c.HTML(http.StatusOK, "group/public", groups)
    }
}

func GetGroups(c *gin.Context) {
    var groups []models.Group
    claims := jwt.ExtractClaims(c)
    user, _ := c.Get("id")
    log.Warn("CLAIM ", claims["permissions"])
    log.Warn("CONTE ", user.(*models.Login).ID)
    _start, _ :=strconv.Atoi(c.DefaultQuery("_start", "0"))
    _end, _ :=strconv.Atoi(c.DefaultQuery("_end", "10"))
    _sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
    result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&groups)
    if result.Error != nil {
        c.AbortWithStatus(500)
        log.Warn("Get groups failed.")
    } else {
        c.JSON(http.StatusOK, groups)
    }
}

func GetGroup(c *gin.Context) {
    // Get model if exist
    var group models.Group
    if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Get group failed.")
        return
    }
    c.JSON(http.StatusOK, group)
}

func PostGroup(c *gin.Context) {
    // Validate input
    var input CreateGroupInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Post group failed.")
        fmt.Println(err)
        return
    }
    group := models.Group{
        Short:		input.Short,
        Name:		input.Name,
        Size:		input.Size,
        // Grouping:	input.Grouping, // TODO
        TribeID:	input.TribeID,
    }
    // Create group
    models.DB.Create(&group)
    c.JSON(http.StatusOK, group)
}

func PutGroup(c *gin.Context) {
    // Validate input
    var input models.Group
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Put group failed.")
        fmt.Println(err)
        return
    }
    // Put Group
    models.DB.Save(&input)
    c.JSON(http.StatusOK, input)
}

func PatchGroup(c *gin.Context) {
    // Get model if exist
    var group models.Group
    if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Patch group failed.")
        return
    }
    // Validate input
    var input UpdateGroupInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        log.Warn("Patch group failed.")
        return
    }
    models.DB.Model(&group).Updates(input)
    c.JSON(http.StatusOK, group)
}

func DeleteGroup(c *gin.Context) {
    // Get model if exist
    var group models.Group
    if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        log.Warn("Delete group failed.")
        return
    }
    models.DB.Delete(&group)
    c.JSON(http.StatusOK, true)
}
