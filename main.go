package main

import (
	"net/http"
	"georgslauf/controllers"
	"georgslauf/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
	log "github.com/sirupsen/logrus"
	jwt "github.com/appleboy/gin-jwt/v2"
	"os"
)

var (
	identityKey =	"id"
	permissionKey =	"permissions"
	contactKey =	"fullName"
	avatarKey =		"avatar"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.Print("Log level ", log.GetLevel())
}

func main() {
	models.ConnectDatabase()
	models.SetEnforcer()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(cors.Default())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	

	controllers.InitTotal()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "action-copier-shredding-landless-marrow-backhand-vacation-doorway-regulator-truck"
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "georgslauf.de",
		Key:         []byte(secret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Login); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
					contactKey: v.Contact,
					avatarKey: v.Avatar,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.Login{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: controllers.Login,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.Login); ok {
				//sub := v.Username
				//obj := c.Request.URL.RequestURI()
				//act := c.Request.Method
				en, _ := models.EN.Enforce(v.Username, c.Request.URL.RequestURI(), c.Request.Method)
				// log.Debug("Enforce(\"", sub, "\",\"", obj, "\",\"", act, "\") is ", en)
				// log.Debug("Reason: ", reason)
				if en {
					//log.Debug("Enforcer passed.")
					return true
				}
			}
			log.Debug("Enforcer blocked.")
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	r.Static("/static", "uploads")
	r.POST("/login/", authMiddleware.LoginHandler)
	r.GET("/refresh/", authMiddleware.RefreshHandler)
	r.GET("/logout/", authMiddleware.LogoutHandler)
	v1 := r.Group("/v1")
	v1.Use(authMiddleware.MiddlewareFunc())
	test := v1.Group("/test")
	test.GET("/", func(c *gin.Context) {
		log.Info("Hello received a GET request..")
	})
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
	rule := v1.Group("/rules")
	{
		rule.GET("/", controllers.GetRules)
		//rule.GET("/:id", controllers.GetRule)
		rule.POST("/", controllers.PostRule)
		//rule.PUT("/:id", controllers.PutRule)
		rule.DELETE("/:id", controllers.DeleteRule)
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
	stationtop := v1.Group("/stationtops")
	{
		stationtop.GET("/", controllers.GetStationTops)
		stationtop.GET("/:id", controllers.GetStationTop)
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
	content := v1.Group("/contents")
	{
		content.GET("/", controllers.GetContents)
		content.GET("/:id", controllers.GetContent)
		content.POST("/", controllers.PostContent)
		content.PUT("/:id", controllers.PutContent)
		content.DELETE("/:id", controllers.DeleteContent)
		content.PATCH("/:id", controllers.PatchContent)
	}
	contenttype := v1.Group("/contenttypes")
	{
		contenttype.GET("/", controllers.GetContentTypes)
		contenttype.GET("/:id", controllers.GetContentType)
		contenttype.POST("/", controllers.PostContentType)
		contenttype.PUT("/:id", controllers.PutContentType)
		contenttype.DELETE("/:id", controllers.DeleteContentType)
		contenttype.PATCH("/:id", controllers.PatchContentType)
	}
	run := v1.Group("/runs")
	{
		run.GET("/", controllers.GetRuns)
		run.GET("/:id", controllers.GetRun)
		run.POST("/", controllers.PostRun)
		run.PUT("/:id", controllers.PutRun)
		run.DELETE("/:id", controllers.DeleteRun)
		run.PATCH("/:id", controllers.PatchRun)
	}
	config := v1.Group("/configs")
	{
		config.GET("/", controllers.GetConfigs)
		config.GET("/:id", controllers.GetConfig)
		config.POST("/", controllers.PostConfig)
		config.PUT("/:id", controllers.PutConfig)
		config.DELETE("/:id", controllers.DeleteConfig)
		config.PATCH("/:id", controllers.PatchConfig)
	}
	log.Info("API ready.")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
