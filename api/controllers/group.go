package controllers


import (
	"georgslauf/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)


type CreateGroupInput struct {
	Short		string	`json:"short" binding:"required"`
	Name		string	`json:"name" binding:"required"`
	Size		uint	`json:"size" binding:"required"`
	GroupingID	uint	`json:"groupingID" binding:"required"`
	TribeID		uint	`json:"TribeID" binding:"required"`
}


type UpdateGroupInput struct {
	Name		string	`json:"name"`
	Size		uint	`json:"size"`
	Short		string	`json:"short"`
	TribeID		uint	`json:"TribeID"`
	// Grouping?
}


func GetGroupsByTribe(c *gin.Context) {
	var groups []models.Group
	result := models.DB.Where("tribe_id = ?", c.Param("tribeid")).Find(&groups)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get stations failed.")
	}
	c.HTML(http.StatusOK, "group/tribe", groups)
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
	}
	c.HTML(http.StatusOK, "group/public", groups)
}


func GetGroups(c *gin.Context) []models.Group {
	var groups []models.Group
	_start, _ :=strconv.Atoi(c.DefaultQuery("_start", "0"))
	_end, _ :=strconv.Atoi(c.DefaultQuery("_end", "10"))
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Limit(_end - _start).Offset(_start).Order(_sortOrder).Find(&groups)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get groups failed.")
	}
	return groups
}


func GetGroupsWithPointsByStationID(c *gin.Context) []models.GroupWithStationPoints {
	// TODO use a view here
	// soft delete is not automatically applied here
	var groups []models.GroupWithStationPoints
	result := models.DB.Table("groups").Select("groups.name, gp.value, gp.updated_at").Joins("left join group_points gp on gp.group_id = groups.id and gp.station_id = ?", c.MustGet("station")).Where("groups.deleted_at is null").Scan(&groups)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Error(result.Error)
	}
	return groups
}


func GetGroup(c *gin.Context) models.Group {
	// Get model if exist
	var group models.Group
	if err := models.DB.Where("id = ?", c.Param("id")).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get group failed.")
	}
	return group
}


func PostGroup(c *gin.Context) {
	// Validate input
	var input CreateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post group failed.")
		// fmt.Println(err)
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
		// fmt.Println(err)
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
