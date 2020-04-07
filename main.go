package main

import (
	"fmt"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ledinhbao/blog/core"
	"github.com/ledinhbao/blog/packages/sports/strava"
)

const (
	userkey       = "user"
	dbInstance    = "database"
	adminkey      = "admin"
	stravaAuthURL = "https://www.strava.com/oauth/authorize?client_id=44814&" +
		"redirect_uri=localhost:9096&response_type=code&approval_prompt=auto&scope=activity:read"
	authUserID = "AuthUserID"
	authUser   = string("AuthUser")
	appName    = string("LDBBlog")
	appUUID    = string("29fjkshD87HSfwy")

	RoleSuperAdmin = int(89)
	RoleAdmin      = int(55)
	RoleModerator  = int(34)
	RoleWriter     = int(21)
	RoleViewer     = int(1)
)

var log *zap.Logger
var logCfg zap.Config

// Application required, panic if one of these failed to achieve:
// 	- Config object is loaded.
// 	- Logger is successully initialized.
//	- An connection to database.
func main() {
	var err error

	// Load Configuration file: config.json
	var config core.Config
	config, err = core.NewConfigFromJSONFile("config.json")
	if err != nil {
		panic("Unabled to load Config object: " + err.Error())
	}

	// Load logger
	initLogger()

	// Setup environment
	appMode, err := config.StringValueForKey("application.mode")
	if err == nil && appMode == "release" {
		// Log Info level in release mode
		logCfg.Level.SetLevel(zap.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	// LOAD database
	dbProfile, _ := config.StringValueForKey("application.db-profile")
	dbConfig, err := config.ConfigValueForKey("database." + dbProfile)
	db, err := loadDatabase(dbConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to load database information: %s", err.Error()))
	}
	defer db.Close()
	// Model Migration
	modelMigration(db)

	router := gin.New()
	if appMode == "debug" {
		router = gin.Default()
	}

	// router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(log, true))
	router.Use(location.Default())

	initSession(router)

	// Set database instance for global use
	router.Use(dbHandler(db))

	// Serving static resources
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))

	// Front-end Template setup
	router.HTMLRender = ginview.New(goview.Config{
		Root:         "views/frontend",
		Extension:    ".html",
		Master:       "layout/master",
		Partials:     []string{},
		DisableCache: true,
	})

	// Back-end Template setup
	backendViewMiddleware := ginviewBackendTemplateMiddleware()

	adminGeneralRoute := router.Group("/admin", backendViewMiddleware)
	initNonAuthAdminRoutes(adminGeneralRoute)

	adminRouter := router.Group("/admin", backendViewMiddleware)
	adminRouter.Use(AuthRequired)
	initAdminRoutes(adminRouter)

	suRouter := router.Group("/su", backendViewMiddleware)
	suRouter.Use(SuperAdminRequired())
	initSuperAdminRoutes(suRouter)

	initializeRoutes(router)
	inititalizePostRoutes(router)

	initStravaModule(appMode, &config, router, backendViewMiddleware)
	router.Run(":9096")
}

func initLogger() {
	var err error
	// Setup logger using uber-go/zap
	logCfg = zap.NewDevelopmentConfig()
	// logCfg.OutputPaths = []string{"logs/blog.log"}
	log, err = logCfg.Build()
	defer log.Sync()
	if err != nil {
		panic("Failed to init log module" + err.Error())
	}
}

func initStravaModule(mode string, cfg *core.Config, r *gin.Engine, mdws ...gin.HandlerFunc) {
	strKey := fmt.Sprintf("strava.%s.webhook-callback", mode)
	callback, _ := cfg.StringValueForKey(strKey)
	clientID, _ := cfg.StringValueForKey("client-id")
	clientSecret, _ := cfg.StringValueForKey("client-secret")
	// main path: /admin/strava/* -> redirect: /admin/strava/dashboard
	strava.SetConfig(strava.Config{
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		PathPrefix:      "/admin",
		PathRedirect:    "/dashboard",
		GlobalDatabase:  dbInstance,
		URLCallbackHost: callback,
	})
	strava.InitializeRoutes(r)

	for _, middleware := range mdws {
		stravaClubRouter := r.Group("/admin/strava/clubs", middleware)
		initStravaRoutes(stravaClubRouter)
	}
}
