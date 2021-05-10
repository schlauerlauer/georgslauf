package main

import (
	"georgslauf/controllers"
	"georgslauf/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:	[]string{"http://localhost:3000"},
		AllowMethods:	[]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:	[]string{"Origin", "Content-Length", "Content-Type", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:	12 * time.Hour,
		ExposeHeaders:	[]string{"x-total-count","Content-Range"},
	}))
	models.ConnectDatabase()
	v1 := r.Group("/v1")
	group := v1.Group("/groups")
	{
		group.GET("/", controllers.GetGroups)
		group.GET("/:id", controllers.GetGroup)
		group.POST("/", controllers.PostGroup)
		group.PATCH("/:id", controllers.PatchGroup)
		group.PUT("/:id", controllers.PutGroup)
		group.DELETE("/:id", controllers.DeleteGroup)
		group.OPTIONS("/", Options)
	}
	tribe := v1.Group("/tribes")
	{
		tribe.GET("/", controllers.GetTribes)
		tribe.GET("/:id", controllers.GetTribe)
		tribe.POST("/", controllers.PostTribe)
		tribe.PATCH("/:id", controllers.PatchTribe)
		tribe.DELETE("/:id", controllers.DeleteTribe)
	}
	role := v1.Group("/roles")
	{
		role.GET("/", controllers.GetRoles)
		role.GET("/:id", controllers.GetRole)
		role.POST("/", controllers.PostRole)
		role.PATCH("/:id", controllers.PatchRole)
		role.DELETE("/:id", controllers.DeleteRole)
	}
	station := v1.Group("/stations")
	{
		station.GET("/", controllers.GetStations)
		station.GET("/:id", controllers.GetStation)
		station.POST("/", controllers.PostStation)
		station.PATCH("/:id", controllers.PatchStation)
		station.DELETE("/:id", controllers.DeleteStation)
	}
	grouppoint := v1.Group("/grouppoints")
	{
		grouppoint.GET("/", controllers.GetGroupPoints)
		grouppoint.GET("/:id", controllers.GetGroupPoint)
		grouppoint.POST("/", controllers.PostGroupPoint)
		grouppoint.PATCH("/:id", controllers.PatchGroupPoint)
		grouppoint.DELETE("/:id", controllers.DeleteGroupPoint)
	}
	stationpoint := v1.Group("/stationpoints")
	{
		stationpoint.GET("/", controllers.GetStationPoints)
		stationpoint.GET("/:id", controllers.GetStationPoint)
		stationpoint.POST("/", controllers.PostStationPoint)
		stationpoint.PATCH("/:id", controllers.PatchStationPoint)
		stationpoint.DELETE("/:id", controllers.DeleteStationPoint)
	}
	login := v1.Group("/logins")
	{
		login.GET("/", controllers.GetLogins)
		login.GET("/:id", controllers.GetLogin)
		login.POST("/", controllers.PostLogin)
		login.PATCH("/:id", controllers.PatchLogin)
		login.DELETE("/:id", controllers.DeleteLogin)
	}

	r.Run()
}

// Options common response for rest options
func Options(c *gin.Context) {
	Origin := c.MustGet("CorsOrigin").(string)

	c.Writer.Header().Set("Access-Control-Allow-Origin", Origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}