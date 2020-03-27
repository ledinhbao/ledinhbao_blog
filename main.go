package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/ledinhbao/blog/core"
	"github.com/ledinhbao/blog/packages/models"
	"github.com/ledinhbao/blog/packages/sports/strava"
)

// L93hxwPc8r
// ledinhbao_axis
// ledinhbao_blog

const (
	userkey       = "user"
	dbInstance    = "database"
	adminkey      = "admin"
	stravaAuthURL = "https://www.strava.com/oauth/authorize?client_id=44814&" +
		"redirect_uri=localhost:9096&response_type=code&approval_prompt=auto&scope=activity:read"
	authUserID = "AuthUserID"
)

func hashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

type Configuration struct {
	RootUserSetup bool `json:"RootUserSetup"`
}

func readConfig() Configuration {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Cannot read config file" + err.Error())
	}
	data := Configuration{}
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		panic("Unmarshal config data failed: " + err.Error())
	}
	return data
}

// AuthRequired is a middleware to check if the user is authorized or not.
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// unauthorize will be transfer to /admin/login
		session.AddFlash("Unauthorized!", "is-error")
		session.Save()
		c.Redirect(http.StatusPermanentRedirect, "/admin/login")
	}
}

func dbHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(dbInstance, db)
		c.Next()
	}
}

func RandString() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func main() {
	router := gin.Default()

	cookieName := RandString()
	router.Use(sessions.Sessions("ledinhbao_com_sessions", sessions.NewCookieStore([]byte(cookieName))))
	// db, err := sqlx.Connect("mysql", "ledinhbao_axis:L93hxwPc8r@/ledinhbao_blog")
	// db, err := sqlx.Connect("sqlite3", "database.db")
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("Cannot connect to database." + err.Error())
	}
	defer db.Close()
	// Set database instance for global use
	router.Use(dbHandler(db))

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(core.Setting{})

	// Serving static resources
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))

	// router.LoadHTMLGlob("templates/*")

	router.HTMLRender = ginview.New(goview.Config{
		Root:         "views/frontend",
		Extension:    ".html",
		Master:       "layout/master",
		Partials:     []string{},
		DisableCache: true,
	})

	router.GET("/setup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "setup", gin.H{
			"message": "Root Setup",
		})
	})

	router.GET("/", displayPosts)

	backendViewMiddleware := ginview.NewMiddleware(goview.Config{
		Root:         "views/backend",
		Extension:    ".html",
		Master:       "layout/master",
		Partials:     []string{},
		DisableCache: true,
	})

	adminGeneralRoute := router.Group("/admin", backendViewMiddleware)
	{
		adminGeneralRoute.GET("/login", showAdminLoginPage)
		adminGeneralRoute.POST("/postLogin", adminPostLogin)
	}

	adminRoute := router.Group("/admin", backendViewMiddleware)
	adminRoute.Use(AuthRequired)
	{
		adminRoute.GET("/dashboard", displayAdminDashboard)
		adminRoute.GET("/", displayAdminIndex)
		adminRoute.GET("/logout", adminLogout)
	}

	router.GET("/admin/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-register", gin.H{
			"message": "Admin Register",
		})
	})

	router.POST("/admin/register", func(c *gin.Context) {
		var formData = models.User{}
		formData.Username = c.PostForm("username")
		formData.SetPassword(c.PostForm("password"))
		formData.PasswordConfirm = c.PostForm("password2")
		formData.Role = 1

		message := ""

		db.Create(&formData)
		c.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	})
	initializeRoutes(router)
	inititalizePostRoutes(router)

	// main path: /admin/strava/* -> redirect: /admin/dashboard
	strava.SetConfig(strava.Config{
		ClientID:        "44814",
		ClientSecret:    "c44a13c4308b3b834320ae5e3648d6c7855980a3",
		PathPrefix:      "/admin",
		PathRedirect:    "/dashboard",
		GlobalDatabase:  dbInstance,
		URLCallbackHost: "http://bc7b66a4.ngrok.io",
	})
	strava.InitializeRoutes(router)
	db.AutoMigrate(&strava.Link{})
	db.AutoMigrate(&strava.Athlete{})

	// go strava.CreateSubscription(db)
	router.Run(":9096")
}

func setSession(c *gin.Context, key string, value string) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

func displayAdminIndex(c *gin.Context) {
	ginview.HTML(c, http.StatusOK, "admin-index", gin.H{})
}

func showAdminLoginPage(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes("is-error")
	session.Save()
	ginview.HTML(c, http.StatusOK, "admin-login.html", gin.H{
		"errors": flashes,
	})
}

func adminPostLogin(c *gin.Context) {
	db := c.MustGet(dbInstance).(*gorm.DB)
	user := models.User{}
	passwordFromRequest := c.PostForm("password")
	db.Where("username = ?", c.PostForm("username")).First(&user)

	if user.TryPassword(passwordFromRequest) {
		session := sessions.Default(c)
		session.Set(userkey, user.Username)
		session.Set(authUserID, user.ID)
		session.Save()
		c.Redirect(http.StatusFound, "/admin/dashboard")
		// c.Abort()
	} else {
		session := sessions.Default(c)
		session.AddFlash("Wrong password", "is-error")
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/admin/login")
		// c.Abort()
	}
}

func adminLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/admin/login")
}

func displayAdminDashboard(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get(authUserID)

	db := c.MustGet(dbInstance).(*gorm.DB)
	userInfo := models.User{}
	stravaInfo := strava.Athlete{}
	stravaLink := strava.Link{}

	db.Where("id = ?", userID).First(&userInfo)
	db.Where(&strava.Link{UserID: userInfo.ID}).First(&stravaLink)
	db.Where(&strava.Athlete{Username: stravaInfo.Username}).First(&stravaInfo)

	var stravaSetting core.Setting
	db.Where(core.Setting{Key: strava.ActiveConfig().SubscriptionDBKey}).First(&stravaSetting)

	ginview.HTML(c, http.StatusOK, "admin-dashboard", gin.H{
		"user":              userInfo,
		"strava_link":       stravaLink,
		"athelete":          stravaInfo,
		"IsStravaConnected": stravaLink.ID > 0,
		"StravaRevokeURL":   strava.ActiveConfig().GetRevokeURLFor(stravaInfo.Username),

		"IsStravaSubscribed":   stravaSetting.ID > 0,
		"stravaSubscriptionID": stravaSetting.Value,
	})
}
