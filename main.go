package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	gorm.Model
	FirstName string
	LastName string
	Password string
}

func main() {
	dsn := "host=localhost user=postgres password=123456 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// sync table info
	db.AutoMigrate(&User{})

	router := gin.Default()

	// request for creating a user
	router.POST("/createUser", func(c *gin.Context) {
		firstName := c.Query("firstName")
		lastName := c.Query("lastName")
		password := c.PostForm("password")

		user := User{
			FirstName: firstName,
			LastName:  lastName,
			Password:  password,
		}

		db.Create(&user)
	})

	// delete a user by id
	router.DELETE("/deleteUser", func(c *gin.Context) {
		userId, _ := strconv.ParseUint(c.Query("id"), 10, 64)
		db.Delete(&User{}, userId)
	})

	// update user first name
	router.PUT("/updateFirstName", func(c *gin.Context) {
		userId, _ := strconv.ParseUint(c.Query("id"), 10, 64)
		newFirstName := c.PostForm("newFirstName")
		db.Model(&User{}).Where("ID = ?", userId).Update("FirstName", newFirstName)
	})

	// update user last name
	router.PUT("/updateLastName", func(c *gin.Context) {
		userId, _ := strconv.ParseUint(c.Query("id"), 10, 64)
		newFirstName := c.PostForm("newLastName")
		db.Model(&User{}).Where("ID = ?", userId).Update("LastName", newFirstName)
	})

	// update user password
	router.PUT("/updatePassword", func(c *gin.Context) {
		userId, _ := strconv.ParseUint(c.Query("id"), 10, 64)
		newFirstName := c.PostForm("newPassword")
		db.Model(&User{}).Where("ID = ?", userId).Update("Password", newFirstName)
	})

	// search user by id
	router.GET("/searchUser", func(c *gin.Context) {
		userId, _ := strconv.ParseUint(c.Query("id"), 10, 64)
		var user User
		db.First(&user, userId)
		c.JSON(http.StatusOK, gin.H{
			"ID": user.ID,
			"First Name": user.FirstName,
			"Last Name": user.LastName,
		})
	})


	router.Run()
}


