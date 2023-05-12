package main

import (
    "net/http"
    "georgslauf/controllers"
    "georgslauf/models"
    "github.com/gin-gonic/gin"
    // "time"
    log "github.com/sirupsen/logrus"
    // jwt "github.com/appleboy/gin-jwt/v2"
    "gopkg.in/yaml.v2"
    "os"
)

var (
    // identityKey =   "id"
    // permissionKey = "permissions"
    // emailKey =      "email"
    // avatarKey =     "avatar"
    // loginKey =      "login"
    cfg = newConfig("./config.yaml")
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

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "https://admin.georgslauf.de")
        //c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Range, X-Total-Count")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Total-Count")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}

func main() {
    models.ConnectDatabase(cfg.Database.Postgresql)
    // controllers.InitTotal()

    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(CORS())

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "api.georgslauf.de",
            "version": "23.2.0-alpha",
        })
    })
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    // authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
    //     Realm:       "georgslauf.de",
    //     Key:         []byte(cfg.Server.Secret),
    //     Timeout:     time.Hour,
    //     MaxRefresh:  time.Hour,
    //     IdentityKey: identityKey,
    //     PayloadFunc: func(data interface{}) jwt.MapClaims {
    //         if v, ok := data.(*models.Login); ok {
    //             return jwt.MapClaims{
    //                 identityKey: v.Username,
    //                 emailKey: v.Email,
    //                 avatarKey: v.Avatar,
    //                 permissionKey: v.Permissions,
    //                 loginKey: v.ID,
    //             }
    //         }
    //         return jwt.MapClaims{}
    //     },
    //     IdentityHandler: func(c *gin.Context) interface{} {
    //         claims := jwt.ExtractClaims(c)
    //         return &models.Login{
    //             Username: claims[identityKey].(string),
    //         }
    //     },
    //     Authenticator: controllers.Login,
    //     Authorizator: func(data interface{}, c *gin.Context) bool {
    //         if v, ok := data.(*models.Login); ok {
    //             //sub := v.Username
    //             //obj := c.Request.URL.RequestURI()
    //             //act := c.Request.Method
    //             en, _ := models.EN.Enforce(v.Username, c.Request.URL.RequestURI(), c.Request.Method)
    //             // log.Debug("Enforce(\"", sub, "\",\"", obj, "\",\"", act, "\") is ", en)
    //             // log.Debug("Reason: ", reason)
    //             if en {
    //                 //log.Debug("Enforcer passed.")
    //                 return true
    //             }
    //         }
    //         log.Debug("Enforcer blocked.")
    //         return false
    //     },
    //     Unauthorized: func(c *gin.Context, code int, message string) {
    //         c.JSON(code, gin.H{
    //             "code":    code,
    //             "message": message,
    //         })
    //     },
    //     TokenLookup: "header: Authorization, query: token, cookie: jwt",
    //     TokenHeadName: "Bearer",
    //     TimeFunc: time.Now,
    // })
    // if err != nil {
    //     log.Fatal("JWT Error:" + err.Error())
    // }
    public := r.Group("/public")
    {
        public.Static("/media", "uploads")
        public.GET("/content/:ct", controllers.GetPublicContent)
        public.GET("/stations", controllers.GetPublicStations)
        public.GET("/stations/:id", controllers.GetPublicStation)
        public.GET("/groups", controllers.GetPublicGroups)
        public.GET("/groups/:id", controllers.GetPublicGroup)
    }
    // auth := r.Group("/auth")
    // {
    //     auth.POST("/login", authMiddleware.LoginHandler)
    //     auth.GET("/refresh", authMiddleware.RefreshHandler)
    //     auth.GET("/logout", authMiddleware.LogoutHandler)
    // }
    login := r.Group("/logins")
    // login.Use(authMiddleware.MiddlewareFunc())
    {
        login.GET("", controllers.GetLogins)
        login.GET(":id", controllers.GetLogin)
        login.POST("", controllers.PostLogin)
        login.PUT(":id", controllers.PutLogin)
        login.DELETE(":id", controllers.DeleteLogin)
        login.PATCH(":id", controllers.PatchLogin)
    }
    group := r.Group("/groups")
    // group.Use(authMiddleware.MiddlewareFunc())
    {
        group.GET("", controllers.GetGroups)
        group.GET(":id", controllers.GetGroup)
        group.POST("", controllers.PostGroup)
        group.PUT(":id", controllers.PutGroup)
        group.DELETE(":id", controllers.DeleteGroup)
        group.PATCH(":id", controllers.PatchGroup)
    }
    tribe := r.Group("/tribes")
    // tribe.Use(authMiddleware.MiddlewareFunc())
    {
        tribe.GET("", controllers.GetTribes)
        tribe.GET(":id", controllers.GetTribe)
        tribe.POST("", controllers.PostTribe)
        tribe.PUT(":id", controllers.PutTribe)
        tribe.DELETE(":id", controllers.DeleteTribe)
        tribe.PATCH(":id", controllers.PatchTribe)
        tribe.GET("/stations:loginid", controllers.GetStationsByLogin)
        tribe.GET("/groups:loginid", controllers.GetGroupsByLogin)
    }
    station := r.Group("/stations")
    // station.Use(authMiddleware.MiddlewareFunc())
    {
        station.GET("", controllers.GetStations)
        station.GET(":id", controllers.GetStation)
        station.POST("", controllers.PostStation)
        station.PUT(":id", controllers.PutStation)
        station.DELETE(":id", controllers.DeleteStation)
        station.PATCH(":id", controllers.PatchStation)
    }
    grouppoint := r.Group("/grouppoints")
    // grouppoint.Use(authMiddleware.MiddlewareFunc())
    {
        grouppoint.GET("", controllers.GetGroupPoints)
        grouppoint.GET(":id", controllers.GetGroupPoint)
        grouppoint.POST("", controllers.PostGroupPoint)
        grouppoint.PUT(":id", controllers.PutGroupPoint)
        grouppoint.DELETE(":id", controllers.DeleteGroupPoint)
        grouppoint.PATCH(":id", controllers.PatchGroupPoint)
    }
    grouptop := r.Group("/grouptops")
    // grouptop.Use(authMiddleware.MiddlewareFunc())
    {
        grouptop.GET("", controllers.GetGroupTops)
        grouptop.GET(":id", controllers.GetGroupTop)
    }
    stationpoint := r.Group("/stationpoints")
    // stationpoint.Use(authMiddleware.MiddlewareFunc())
    {
        stationpoint.GET("", controllers.GetStationPoints)
        stationpoint.GET(":id", controllers.GetStationPoint)
        stationpoint.POST("", controllers.PostStationPoint)
        stationpoint.PUT(":id", controllers.PutStationPoint)
        stationpoint.DELETE(":id", controllers.DeleteStationPoint)
        stationpoint.PATCH(":id", controllers.PatchStationPoint)
    }
    stationtop := r.Group("/stationtops")
    // stationtop.Use(authMiddleware.MiddlewareFunc())
    {
        stationtop.GET("", controllers.GetStationTops)
        stationtop.GET(":id", controllers.GetStationTop)
    }
    grouping := r.Group("/groupings")
    // grouping.Use(authMiddleware.MiddlewareFunc())
    {
        grouping.GET("", controllers.GetGroupings)
        grouping.GET(":id", controllers.GetGrouping)
        grouping.POST("", controllers.PostGrouping)
        grouping.PUT(":id", controllers.PutGrouping)
        grouping.DELETE(":id", controllers.DeleteGrouping)
        grouping.PATCH(":id", controllers.PatchGrouping)
    }
    content := r.Group("/content")
    // content.Use(authMiddleware.MiddlewareFunc())
    {
        content.GET("", controllers.GetContents)
        content.GET(":id", controllers.GetContent)
        content.POST("", controllers.PostContent)
        content.PUT(":id", controllers.PutContent)
        content.DELETE(":id", controllers.DeleteContent)
        content.PATCH(":id", controllers.PatchContent)
    }
    contenttype := r.Group("/contenttypes")
    // contenttype.Use(authMiddleware.MiddlewareFunc())
    {
        contenttype.GET("", controllers.GetContentTypes)
        contenttype.GET(":id", controllers.GetContentType)
        contenttype.POST("", controllers.PostContentType)
        contenttype.PUT(":id", controllers.PutContentType)
        contenttype.DELETE(":id", controllers.DeleteContentType)
        contenttype.PATCH(":id", controllers.PatchContentType)
    }
    run := r.Group("/runs")
    // run.Use(authMiddleware.MiddlewareFunc())
    {
        run.GET("", controllers.GetRuns)
        run.GET(":id", controllers.GetRun)
        run.POST("", controllers.PostRun)
        run.PUT(":id", controllers.PutRun)
        run.DELETE(":id", controllers.DeleteRun)
        run.PATCH(":id", controllers.PatchRun)
    }
    config := r.Group("/config")
    // config.Use(authMiddleware.MiddlewareFunc())
    {
        config.GET("", controllers.GetConfigs)
        config.GET(":id", controllers.GetConfig)
        config.POST("", controllers.PostConfig)
        config.PUT(":id", controllers.PutConfig)
        config.DELETE(":id", controllers.DeleteConfig)
        config.PATCH(":id", controllers.PatchConfig)
    }
    r.GET("/metrics", gin.BasicAuth(gin.Accounts{
        cfg.Server.Metrics.Username: cfg.Server.Metrics.Password,
    }), controllers.MetricsHandler())
    
    log.Info("Listening on ", cfg.Server.Host, ":", cfg.Server.Port)
    if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
        log.Fatal(err)
    }
}
