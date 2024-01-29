package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"test.com/models"
	"test.com/services"
)

type DBConfig struct {
	Username string `mapstructure:"database.username"`
	Password string `mapstructure:"database.password"`
	DB       string `mapstructure:"database.db"`
	Host     string `mapstructure:"database.host"`
	Port     string `mapstructure:"database.port"`
}

var gdbContext *gorm.DB

func getPostService() *services.PostService {
	return &services.PostService{DbContext: gdbContext}
}

func getUserService() *services.UserService {
	return &services.UserService{DbContext: gdbContext}
}

func getCommentService() *services.CommentService {
	return &services.CommentService{DbContext: gdbContext}
}

func main() {
	/*
		viper.SetConfigFile("config.yaml")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}

		var config DBConfig
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Println("Error unmarshaling config:", err)
			return
		}

		fmt.Printf("%+v", config)
		fmt.Println("Username:", config.Username)
		fmt.Println("Password:", config.Password)
		fmt.Println("DB:", config.DB)
		fmt.Println("Host:", config.Host)
		fmt.Println("Port:", config.Port)
		return
	*/

	var err error
	dsn := "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8mb3&parseTime=True&loc=Local"
	gdbContext, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("DB connection error:", err)
	}
	defer gdbContext.AutoMigrate(&models.User{}, &models.Comment{}, &models.Post{})

	err = http.ListenAndServe(":8080", Install())
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func Install() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handleIndex).Methods("PUT")

	router.HandleFunc("/user", handlePutUser).Methods("PUT")
	router.HandleFunc("/post", handlePutPost).Methods("PUT")
	router.HandleFunc("/post", handleGetPost).Methods("GET")
	router.HandleFunc("/comment", handlePutComment).Methods("PUT")
	return router
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Hello, World!")
}

func handlePutUser(w http.ResponseWriter, r *http.Request) {
	getUserService().Create(w, r)
}

func handlePutPost(w http.ResponseWriter, r *http.Request) {
	getPostService().Create(w, r)
}

func handlePutComment(w http.ResponseWriter, r *http.Request) {
	getCommentService().Create(w, r)
}

func handleGetPost(w http.ResponseWriter, r *http.Request) {
	getPostService().GetAll(w, r)
}
