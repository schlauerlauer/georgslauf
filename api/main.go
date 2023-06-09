package main

import (
	"net/http"
	"georgslauf/controllers"
	"georgslauf/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	ory "github.com/ory/client-go"
	"errors"
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cfg = newConfig("./config.yaml")
	systemCfg = &models.Config{}
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.Print("Log level ", log.GetLevel(), ".")
	checkConfig()
}

func checkConfig() {
	checkEmptyString(cfg.Server.Port, "api port")
	checkEmptyString(cfg.Server.Secret, "api secret")
	checkEmptyString(cfg.Database.Postgresql.Hostname, "DB hostname")
	checkEmptyString(cfg.Database.Postgresql.Port, "DB port")
	checkEmptyString(cfg.Database.Postgresql.Database, "DB database name")
	checkEmptyString(cfg.Database.Postgresql.Username, "DB username")
	checkEmptyString(cfg.Database.Postgresql.Password, "DB password")
}

func checkEmptyString(checkThis string, description string) {
	if checkThis == "" {
		log.Fatal("needed config var ", description, " is empty.")
	}
}

func newConfig(configPath string) (*models.APIConfig) {
	config := &models.APIConfig{}
	file, err := os.Open(configPath)
	if err != nil {
		log.Error(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		log.Error(err)
	}
	return config
}

type kratosMiddleware struct {
	ory *ory.APIClient
}


func (k *kratosMiddleware) validateSession(request *http.Request) (*ory.Session, error) {
	cookie, err := request.Cookie("ory_kratos_session")
	if err != nil {
		return nil, err
	}
	if cookie == nil {
		return nil, errors.New("no session found in cookie")
	}
	resp, _, err := k.ory.FrontendApi.ToSession(context.Background()).Cookie(cookie.String()).Execute()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        // c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Range, X-Total-Count")
        // c.Writer.Header().Set("Access-Control-Allow-Headers", "HX-Request, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        // c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}


func KratosMiddleware() *kratosMiddleware {
	configuration := ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{
			URL: "http://127.0.0.1:11433", // Kratos Public API
		},
	}
	return &kratosMiddleware{
		ory: ory.NewAPIClient(configuration),
	}
}


func (k *kratosMiddleware) Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := k.validateSession(c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized) // not logged in
			return
		}
		if !*session.Active {
			c.AbortWithStatus(http.StatusUnauthorized) // session not active
			return
		}

		c.Set("station", session.Identity.MetadataPublic["station"])
		c.Set("tribe", session.Identity.MetadataPublic["tribe"])
		c.Set("identity", session.Identity.Id)

		c.Next()
	}
}


func BooleanPermission(permission bool, code int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if (!permission) {
			c.AbortWithStatus(code)
		}
		c.Next()
	}
}


func updateSystemConfig() {
	newSystemCfg, _ := controllers.GetConfig()
	systemCfg = &newSystemCfg
	log.Info("Updated system config")
}


func main() {
	models.ConnectDatabase(cfg.Database.Postgresql)
	gin.SetMode(gin.ReleaseMode)
	updateSystemConfig()

	router := gin.Default()
	router.Use(CORS())
	router.LoadHTMLGlob("templates/*")

	k := KratosMiddleware()

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "https://georgslauf.de/")
	})
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, "23.5.0-alpha")
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	public := router.Group("/public")
	{
		public.Static("/media", "media")
		public.GET("", BooleanPermission(systemCfg.System.PublicStationsVisible, http.StatusNoContent), controllers.GetPublic)
		public.GET("/message", func(c *gin.Context) {
			c.String(http.StatusOK, systemCfg.Notice)
		})
	}

	// group := router.Group("/group", k.Session())
	// {
	// 	group.GET("", controllers.GetGroups)
	// 	group.GET(":id", controllers.GetGroup)
	// 	group.POST("", controllers.PostGroup)
	// 	group.PUT(":id", controllers.PutGroup)
	// 	group.DELETE(":id", controllers.DeleteGroup)
	// 	group.PATCH(":id", controllers.PatchGroup)
	// }

	// TODO check context for station / tribe / none
	home := router.Group("/home", k.Session())
	{
		home.GET("", func(c *gin.Context) {
			// TODO tribe / admin / both / all
			stationID := c.GetString("station")
			// tribeID := c.GetString("tribe")

			// logged in as station
			if (stationID != "") {
				groups := controllers.GetGroupsWithPointsByStationID(c)
				station, _ := controllers.GetStationByID(stationID)

				c.HTML(http.StatusOK, "station/points", gin.H{
					"station": station,
					"groups": groups,
					"groupings": systemCfg.Groupings,
					"enableEdit": systemCfg.System.AllowGroupPoints,
				})
			}
		})
		home.
			PUT("group/:id", BooleanPermission(systemCfg.System.AllowGroupPoints, http.StatusNoContent), controllers.PutGroupPointByStationID)
	}

	settings := router.Group("/settings", k.Session())
	{
		settings.GET("", func(c *gin.Context) {
			// TODO tribe / admin / both / all
			// if settings allow chaning posten settings (size usw)
			stationID := c.GetString("station")
			station, _ := controllers.GetStationByID(stationID)
			c.HTML(http.StatusOK, "settings", gin.H{
				"station": station,
			})
		})
	}

	tribe := router.Group("/tribe", k.Session())
	{
		tribe.GET("/info", func(c *gin.Context) {
			c.HTML(http.StatusOK, "tribe/info", systemCfg.Contact)
		})
		tribe.GET("/stations/:tribeid", controllers.GetStationsByTribe)
		tribe.GET("/groups/:tribeid", controllers.GetGroupsByTribe)
	// tribe.GET("", controllers.GetTribes)
	// 	tribe.GET(":id", controllers.GetTribe)
	// 	tribe.POST("", controllers.PostTribe)
	// 	tribe.PUT(":id", controllers.PutTribe)
	// 	tribe.PATCH(":id", controllers.PatchTribe)
	// 	tribe.GET("/groups:loginid", controllers.GetGroupsByLogin)
	}

	// station := router.Group("/station", k.Session())
	// {
	// 	station.GET("", controllers.GetStations)
	// 	station.GET(":id", controllers.GetStation)
	// 	station.POST("", controllers.PostStation)
	// 	station.PUT(":id", controllers.PutStation)
	// 	station.DELETE(":id", controllers.DeleteStation)
	// 	station.PATCH(":id", controllers.PatchStation)
	// }

	// grouppoint := router.Group("/grouppoint", k.Session())
	// {
	// 	grouppoint.GET("", controllers.GetGroupPoints)
	// 	grouppoint.GET(":id", controllers.GetGroupPoint)
	// 	grouppoint.POST("", controllers.PostGroupPoint)
	// 	grouppoint.PUT(":id", controllers.PutGroupPoint)
	// 	grouppoint.PATCH(":id", controllers.PatchGroupPoint)
	// }

	// stationpoint := router.Group("/stationpoint", k.Session()))
	// {
	// 	stationpoint.GET("", controllers.GetStationPoints)
	// 	stationpoint.GET(":id", controllers.GetStationPoint)
	// 	stationpoint.POST("", controllers.PostStationPoint)
	// 	stationpoint.PUT(":id", controllers.PutStationPoint)
	// 	stationpoint.PATCH(":id", controllers.PatchStationPoint)
	// }

	// config := router.Group("/config", k.Session())
	// {
	// 	config.POST("", controllers.PostConfig)
	// 	config.PUT(":id", controllers.PutConfig)
	// 	config.PATCH(":id", controllers.PatchConfig)
	// }

	log.Info("Listening on ", cfg.Server.Host, ":", cfg.Server.Port)
	if err := router.Run(":"+cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
