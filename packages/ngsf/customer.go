package ngsf

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/core"
	"github.com/qor/validations"
	"go.uber.org/zap"
)

type (
	CustomerService struct {
		DB     *gorm.DB
		Logger *zap.Logger
	}

	Customer struct {
		ID        uint       `json:"customer_id" form:"customer_id" gorm:"primary_key"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at" sql:"index"`
		Fullname  string     `json:"fullname" form:"fullname"`
		DOB       time.Time  `json:"dob" form:"dob"`
		UserID    uint       `json:"user_id" form:"user_id"`
		User      core.User  `gorm:"association_autoupdate:false;association_autocreate:false"`
	}
)

var (
	logger *zap.Logger
)

// TableName add "ngsf_" prefix to table name
func (Customer) TableName() string {
	return "ngsf_customer"
}

func (c Customer) Validate(db *gorm.DB) {
	if c.Fullname == "" {
		db.AddError(validations.NewError(c, "Fullname", "Customer's fullname cannot be empty"))
	}
	return
}

func NewCustomerService(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger) CustomerService {
	db.AutoMigrate(&Customer{})
	return CustomerService{
		DB:     db,
		Logger: logger,
	}
}

func (s *CustomerService) CreateCustomer(newCustomer *Customer) []error {
	err := s.DB.Create(&newCustomer).GetErrors()
	return err
}

func APICreateCustomer(c *gin.Context) {
	// get data from POST form
	var customer Customer
	db := c.MustGet("database").(*gorm.DB)
	if err := c.ShouldBind(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
			"data":    err.Error(),
		})
		return
	}
	if errs := db.Create(&customer).GetErrors(); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
			"data":    errs,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "updated",
		"data":    customer,
	})
}
