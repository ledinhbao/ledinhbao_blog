package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// L93hxwPc8r
// ledinhbao_axis
// ledinhbao_blog

const (
	userkey = "user"
)

type User struct {
	gorm.Model
	ID              uint `gorm:"PRIMARY_KEY"`
	Username        string
	Password        string
	Role            int
	PasswordConfirm string `gorm:"-"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

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

func (f User) ValidatePassword() bool {
	if f.Password != f.PasswordConfirm {
		return false
	}
	return true
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "StatusUnauthorized",
		})
	}
}

func main() {
	router := gin.Default()

	router.Use(sessions.Sessions("MainSession", sessions.NewCookieStore([]byte("ledinhbao"))))
	// db, err := sqlx.Connect("mysql", "ledinhbao_axis:L93hxwPc8r@/ledinhbao_blog")
	// db, err := sqlx.Connect("sqlite3", "database.db")
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("Cannot connect to database." + err.Error())
	}
	defer db.Close()
	db.AutoMigrate(&User{})

	router.LoadHTMLGlob("templates/*")

	router.GET("/setup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "setup", gin.H{
			"message": "Root Setup",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world",
		})
	})

	adminRoute := router.Group("/admin")
	adminRoute.Use(AuthRequired)
	{
		adminRoute.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Admin dashboard",
			})
		})
	}

	router.GET("/admin/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-register", gin.H{
			"message": "Admin Register",
		})
	})

	router.POST("/admin/register", func(c *gin.Context) {
		var formData = User{}
		formData.Username = c.PostForm("username")
		formData.Password, _ = hashPassword(c.PostForm("password"))
		formData.PasswordConfirm = c.PostForm("password2")
		formData.Role = 1

		message := ""

		db.Create(&formData)
		c.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	})
	initializeRoutes(router)
	router.Run(":9096")
}
