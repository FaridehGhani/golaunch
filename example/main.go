package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mohsensamiei/golaunch/packages/convert"
	"github.com/mohsensamiei/golaunch/packages/restful"
)

type User struct {
	ID             uint   `gorm:"primary_key"`
	Username       string `gorm:"size:25;unique"`
	PrimaryRoleID  uint
	PrimaryRole    *Role   `gorm:"foreignkey:PrimaryRoleID"`
	SecondaryRoles []*Role `gorm:"many2many:user_roles"`
}

type Role struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"size:25;unique"`
}

type UserDTO struct {
	Username       string   `json:"username"`
	PrimaryRole    string   `json:"primary_role" map:"PrimaryRole.Name"`
	SecondaryRoles []string `json:"secondary_roles" map:"SecondaryRoles.Name"`
}

type RoleDTO struct {
	Name string `json:"name"`
}

func main() {
	db, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		log.Panic("Database connecting failed")
	}
	db.LogMode(true)
	defer db.Close()

	db.AutoMigrate(&Role{})
	db.AutoMigrate(&User{})

	router := mux.NewRouter()

	router.HandleFunc("/roles", func(response http.ResponseWriter, request *http.Request) {
		createRole(db, response, request)
	}).Methods(http.MethodPost)

	router.HandleFunc("/users", func(response http.ResponseWriter, request *http.Request) {
		createUser(db, response, request)
	}).Methods(http.MethodPost)

	router.HandleFunc("/users", func(response http.ResponseWriter, request *http.Request) {
		getUsers(db, response, request)
	}).Methods(http.MethodGet)

	addr := fmt.Sprintf(":%v", 3000)
	log.Printf("Application serve on http://localost%v", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func createRole(db *gorm.DB, response http.ResponseWriter, request *http.Request) {
	var dto RoleDTO
	if err := restful.ParseRequestBody(request, &dto); err != nil {
		restful.SendResponse(response, request, http.StatusBadRequest, err)
		return
	}

	role := Role{
		Name: dto.Name,
	}
	db.Create(&role)

	restful.SendResponse(response, request, http.StatusCreated, nil)
}

func createUser(db *gorm.DB, response http.ResponseWriter, request *http.Request) {
	var dto UserDTO
	if err := restful.ParseRequestBody(request, &dto); err != nil {
		restful.SendResponse(response, request, http.StatusBadRequest, err)
		return
	}

	var primaryRole Role
	if err := db.Where("name = ?", dto.PrimaryRole).First(&primaryRole).Error; err != nil {
		restful.SendResponse(response, request, http.StatusNotFound, err)
		return
	}

	user := User{
		Username:      dto.Username,
		PrimaryRoleID: primaryRole.ID,
	}
	db.Create(&user)

	var secondaryRoles []*Role
	db.Where("name IN (?)", dto.SecondaryRoles).Find(&secondaryRoles)
	if len(secondaryRoles) > 0 {
		db.Model(&user).Association("SecondaryRoles").Append(secondaryRoles)
	}

	restful.SendResponse(response, request, http.StatusCreated, nil)
}

func getUsers(db *gorm.DB, response http.ResponseWriter, request *http.Request) {
	var users []*User
	db.Preload("PrimaryRole").Preload("SecondaryRoles").Find(&users)

	var dto []UserDTO
	convert.MapToDTO(&users, &dto)

	log.Printf("%+v", dto)

	restful.SendResponse(response, request, http.StatusOK, dto)
}
