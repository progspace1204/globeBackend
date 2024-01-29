package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"test.com/models"
)

type Response struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var SUCCESS string = "success"
var ERROR string = "error"

func writeResponse(w http.ResponseWriter, resp Response) {
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(data))
}

func hashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

type PostUserInferface interface {
	Create(w http.ResponseWriter, r *http.Request)
	CheckUserById(Id int32)
}

type UserService struct {
	PostServiceInferface
	DbContext *gorm.DB
}

func GetUserById(DbContext *gorm.DB, Id int32) (models.User, error) {
	var user models.User
	result := DbContext.Table("users").First(&user, Id)
	return user, result.Error
}

func (ps *UserService) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	username := r.Form.Get("username")
	match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]+$", username)
	if !match {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Username must start with a-Z, can contain a-Z, 0-9, _",
		})
		return
	}
	var count int64
	ps.DbContext.Table("users").Where("`username` = ?", username).Count(&count)
	if count != 0 {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Conflict Username",
		})
		return
	}
	password := r.Form.Get("password")
	if len(password) == 0 {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Password can not be empty",
		})
		return
	}
	password, _ = hashedPassword(password)

	user := models.User{Username: username, Password: password, Author: true}
	result := ps.DbContext.Create(&user)
	if result.Error != nil {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Failed to create account",
		})
	} else {
		writeResponse(w, Response{
			Type:    SUCCESS,
			Message: "",
		})
	}
}
