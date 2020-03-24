package main

import (
	"encoding/json"
	"fmt"
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

	"github.com/ledinhbao/blog/packages/models"
)

// L93hxwPc8r
// ledinhbao_axis
// ledinhbao_blog

const (
	userkey    = "user"
	dbInstance = "database"
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
	fmt.Println(user)
	if user == nil {
		// unauthorize will be transfer to /admin/login
		c.Redirect(http.StatusFound, "/admin/login")
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

	// Serving static resources
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))

	router.LoadHTMLGlob("templates/*")

	router.GET("/setup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "setup", gin.H{
			"message": "Root Setup",
		})
	})

	router.GET("/", displayPosts)

	adminRoute := router.Group("/admin")
	adminRoute.Use(AuthRequired)
	{
		adminRoute.GET("/dashboard", func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{})
		})
		adminRoute.GET("/", displayAdminIndex)
	}
	router.GET("/admin/login", showAdminLoginPage)

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
	router.Run(":9096")
}

func displayAdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-index.html", gin.H{})
}

func showAdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-login.html", gin.H{})
}
