package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ledinhbao/blog/packages/models"
)

func initCustomerRoute(engine *gin.Engine, middlewares ...gin.HandlerFunc) {
	r := engine.Group("/admin", middlewares...)
	{
		r.GET("/customer/list", pageCustomerList)
		r.GET("/customer/new", pageCustomerNew)
		r.POST("/customer/update", updateCustomerDetail)
	}
}

func displayCustomerDetail(c *gin.Context, id uint, args ...string) {
	user, _ := authUserFromSession(c)
	var customer models.Customer
	gdb.Where("id=?", id).First(&customer)
	customerJSON, _ := json.Marshal(customer)
	ginview.HTML(c, http.StatusOK, "admin-customer-detail", gin.H{
		"user":         user,
		"customer":     customer,
		"customerJSON": customerJSON,
	})
}

func pageCustomerNew(c *gin.Context) {
	displayCustomerDetail(c, 0)
	// c.String(http.StatusOK, "Fine")
	// user, _ := authUserFromSession(c)
	// ginview.HTML(c, http.StatusOK, "admin-customer-detail", gin.H{
	// 	"user": user,
	// })
}

func pageCustomerList(c *gin.Context) {
	user, _ := authUserFromSession(c)
	var customers []models.Customer
	db := c.MustGet(dbInstance).(*gorm.DB)
	db.Find(&customers)
	// gdb.Find(&customers)
	customersJSON, _ := json.Marshal(customers)

	s := sessions.Default(c)
	ginview.HTML(c, http.StatusOK, "admin-customer-list", gin.H{
		"user":                user,
		"customers":           customers,
		"customersJSON":       string(customersJSON),
		"errorNotification":   s.Flashes("error"),
		"successNotification": s.Flashes("success"),
	})
	s.Save()
}

func updateCustomerDetail(c *gin.Context) {
	var customer models.Customer
	var errs []error

	if err := c.ShouldBind(&customer); err != nil {
		// c.String(http.StatusBadRequest, fmt.Sprintf("Error when binding form data: %s", err.Error()))
		errs = append(errs, fmt.Errorf("Cannot binding form data to Customer object, %s", err.Error()))
	}

	action := c.PostForm("action")
	if action == "" {
		errs = append(errs, fmt.Errorf("Missing 'action' param in PostForm"))
	}

	var dbErrors []error
	tx := gdb.Begin()
	if customer.ID == 0 {
		dbErrors = tx.Create(&customer).GetErrors()
	} else {
		dbErrors = tx.Model(models.Customer{}).Updates(customer).GetErrors()
	}

	errs = append(errs, dbErrors...)
	if len(errs) > 0 {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error when saving customer data: %v", errs))
		return
	}

	s := sessions.Default(c)
	s.AddFlash("Customer saved", "success")
	s.Save()

	c.Redirect(http.StatusFound, "/admin/customer/list")
}
