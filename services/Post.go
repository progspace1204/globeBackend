package services

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
	"test.com/models"
)

type PostServiceInferface interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type PostService struct {
	PostServiceInferface
	DbContext *gorm.DB
}

func GetPostById(DbContext *gorm.DB, Id int32) (models.Post, error) {
	var post models.Post
	result := DbContext.Table("posts").First(&post, Id)
	return post, result.Error
}

func (ps *PostService) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (ps *PostService) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	user_id := r.Form.Get("user_id")

	uid, err := strconv.ParseInt(user_id, 10, 32)
	if err != nil {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Invalid user_id",
		})
		return
	}
	user, err := GetUserById(ps.DbContext, int32(uid))
	if err != nil {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Invalid user_id",
		})
		return
	}
	if !user.Author {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Has no Author permission",
		})
		return
	}
	if len(title) == 0 {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Title can not be empty",
		})
		return
	}
	if len(content) == 0 {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Content can not be empty",
		})
		return
	}

	result := ps.DbContext.Create(&models.Post{Title: title, Content: content, User_id: int32(uid)})
	if result.Error != nil {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Failed to create a new post",
		})
	} else {
		writeResponse(w, Response{
			Type:    SUCCESS,
			Message: "",
		})
	}
}
