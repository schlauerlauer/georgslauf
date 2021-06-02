package main

import (
	"georgslauf/controllers"
	"georgslauf/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	//"time"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.Print("Log level ", log.GetLevel())
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())
	models.ConnectDatabase()
	controllers.InitTotal()
	v1 := r.Group("/v1")
	login := v1.Group("/logins")
	{
		login.GET("/", controllers.GetLogins)
		login.GET("/:id", controllers.GetLogin)
		login.POST("/", controllers.PostLogin)
		login.PUT("/:id", controllers.PutLogin)
		login.DELETE("/:id", controllers.DeleteLogin)
		login.PATCH("/:id", controllers.PatchLogin)
	}
	group := v1.Group("/groups")
	{
		group.GET("/", controllers.GetGroups)
		group.GET("/:id", controllers.GetGroup)
		group.POST("/", controllers.PostGroup)
		group.PUT("/:id", controllers.PutGroup)
		group.DELETE("/:id", controllers.DeleteGroup)
		group.PATCH("/:id", controllers.PatchGroup)
	}
	tribe := v1.Group("/tribes")
	{
		tribe.GET("/", controllers.GetTribes)
		tribe.GET("/:id", controllers.GetTribe)
		tribe.POST("/", controllers.PostTribe)
		tribe.PUT("/:id", controllers.PutTribe)
		tribe.DELETE("/:id", controllers.DeleteTribe)
		tribe.PATCH("/:id", controllers.PatchTribe)
	}
	role := v1.Group("/roles")
	{
		role.GET("/", controllers.GetRoles)
		role.GET("/:id", controllers.GetRole)
		role.POST("/", controllers.PostRole)
		role.PUT("/:id", controllers.PutRole)
		role.DELETE("/:id", controllers.DeleteRole)
		role.PATCH("/:id", controllers.PatchRole)
	}
	station := v1.Group("/stations")
	{
		station.GET("/", controllers.GetStations)
		station.GET("/:id", controllers.GetStation)
		station.POST("/", controllers.PostStation)
		station.PUT("/:id", controllers.PutStation)
		station.DELETE("/:id", controllers.DeleteStation)
		station.PATCH("/:id", controllers.PatchStation)
	}
	grouppoint := v1.Group("/grouppoints")
	{
		grouppoint.GET("/", controllers.GetGroupPoints)
		grouppoint.GET("/:id", controllers.GetGroupPoint)
		grouppoint.POST("/", controllers.PostGroupPoint)
		grouppoint.PUT("/:id", controllers.PutGroupPoint)
		grouppoint.DELETE("/:id", controllers.DeleteGroupPoint)
		grouppoint.PATCH("/:id", controllers.PatchGroupPoint)
	}
	grouptop := v1.Group("/grouptops")
	{
		grouptop.GET("/", controllers.GetGroupTops)
		grouptop.GET("/:id", controllers.GetGroupTop)
	}
	stationpoint := v1.Group("/stationpoints")
	{
		stationpoint.GET("/", controllers.GetStationPoints)
		stationpoint.GET("/:id", controllers.GetStationPoint)
		stationpoint.POST("/", controllers.PostStationPoint)
		stationpoint.PUT("/:id", controllers.PutStationPoint)
		stationpoint.DELETE("/:id", controllers.DeleteStationPoint)
		stationpoint.PATCH("/:id", controllers.PatchStationPoint)
	}
	grouping := v1.Group("/groupings")
	{
		grouping.GET("/", controllers.GetGroupings)
		grouping.GET("/:id", controllers.GetGrouping)
		grouping.POST("/", controllers.PostGrouping)
		grouping.PUT("/:id", controllers.PutGrouping)
		grouping.DELETE("/:id", controllers.DeleteGrouping)
		grouping.PATCH("/:id", controllers.PatchGrouping)
	}
	content := v1.Group("/content")
	{
		content.GET("/", controllers.GetContents)
		content.GET("/:id", controllers.GetContent)
		content.POST("/", controllers.PostContent)
		content.PUT("/:id", controllers.PutContent)
		content.DELETE("/:id", controllers.DeleteContent)
		content.PATCH("/:id", controllers.PatchContent)
	}
	log.Info("API ready.")
	r.Run()
}
