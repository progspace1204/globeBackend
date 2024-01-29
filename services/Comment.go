package services

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
	"test.com/models"
)

type CommentServiceInferface interface {
	Create(post *models.Post)
}

type CommentService struct {
	PostServiceInferface
	DbContext *gorm.DB
}

func (cs *CommentService) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	user_id := r.Form.Get("user_id")
	post_id := r.Form.Get("post_id")
	content := r.Form.Get("content")

	uid, err := strconv.ParseInt(user_id, 10, 32)
	if err != nil {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Invalid user_id",
		})
		return
	}
	_, err = GetUserById(cs.DbContext, int32(uid))
	if err != nil {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Invalid user_id",
		})
		return
	}

	pid, err := strconv.ParseInt(post_id, 10, 32)
	if err != nil {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Invalid post_id",
		})
		return
	}
	_, err = GetPostById(cs.DbContext, int32(pid))
	if err != nil {
		writeResponse(w, Response{
			Type:    ERROR,
			Message: "Invalid post_id",
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
	result := cs.DbContext.Create(&models.Comment{User_id: int32(uid), Post_id: int32(pid), Content: content})
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
