package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey" yaml:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string `form:"firstName" json:"firstName" yaml:"firstName"`
	LastName string `form:"lastName" json:"lastName" yaml:"lastName"`
	Password string `form:"password" json:"password" yaml:"password"`
}

func getUserId (c *gin.Context) (userId uint64) {
	userId, _ = strconv.ParseUint(c.Param("id"), 10, 64)
	return
}

// request for creating a user
func createUser (router *gin.Engine, db *gorm.DB) {
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
}

// delete a user by id
func deleteUser (router *gin.Engine, db *gorm.DB) {
	router.DELETE("/deleteUser/:id", func(c *gin.Context) {
		fmt.Println("###ID:", getUserId(c))
		db.Delete(&User{}, getUserId(c))
	})
}

// update user info, have to include all fields
func updateUserInfo (router *gin.Engine, db *gorm.DB) {
	router.PUT("/updateUserInfo", func(c *gin.Context) {
		contentType := c.Request.Header["Content-Type"][0]
		fmt.Println("content type:", contentType)

		var user User
		if contentType == "application/x-yaml" {
			if err := c.ShouldBindYAML(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		} else if contentType == "application/json" {
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		fmt.Println("user:", user)
		fmt.Println("user ID:", user.ID)
		db.Model(&user).Omit("CreatedAt", "UpdatedAt", "DeletedAt").Updates(&user)
		//var userToUpdate User
		//db.First(&userToUpdate, user.ID)
		//userToUpdate = user
		//db.Save(&userToUpdate)
	})
}

// search user by id
func searchUser (router *gin.Engine, db *gorm.DB) {
	router.GET("/searchUser", func(c *gin.Context) {
		userId := getUserId(c)
		var user User
		db.First(&user, userId)
		c.JSON(http.StatusOK, &user)
	})
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

	createUser(router, db)
	deleteUser(router, db)
	updateUserInfo(router, db)
	searchUser(router, db)

	// update user first name
	router.PUT("/updateFirstName/:id", func(c *gin.Context) {
		userId := getUserId(c)
		newFirstName := c.PostForm("newFirstName")
		db.Model(&User{}).Where("ID = ?", userId).Update("FirstName", newFirstName)
	})

	// update user last name
	router.PUT("/updateLastName/:id", func(c *gin.Context) {
		userId := getUserId(c)
		newFirstName := c.PostForm("newLastName")
		db.Model(&User{}).Where("ID = ?", userId).Update("LastName", newFirstName)
	})

	// update user password
	router.PUT("/updatePassword/:id", func(c *gin.Context) {
		userId := getUserId(c)
		newFirstName := c.PostForm("newPassword")
		db.Model(&User{}).Where("ID = ?", userId).Update("Password", newFirstName)
	})

	router.Run()
}


