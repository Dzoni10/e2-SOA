package handler

import (
	"database-example/auth"
	"database-example/model"
	"database-example/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *service.UserService
}

func (handler *UserHandler) Get(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("User with id %s", idStr)

	id, err := strconv.Atoi(idStr)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"error": "Invalid ID format"}`))
		return
	}

	user, err := handler.UserService.FindUser(id)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"error": "User not found"}`))
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(user)
}

func (handler *UserHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var user model.User

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.UserService.Create(&user)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(user)
}

func (handler *UserHandler) Login(writer http.ResponseWriter, req *http.Request) {

	var ceredentials struct {
		Username string `json:"username"`
		Password string
	}

	err := json.NewDecoder(req.Body).Decode(&ceredentials)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"error": "Invalid input"}`))
		return
	}

	user, err := handler.UserService.Authenticate(ceredentials.Username, ceredentials.Password)

	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"error": "Invalid credentials"}`))
		return
	}

	token, err := auth.GenerateJWT(user.ID, int(user.Role))

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"error": "Failed to generate token"}`))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{
		"token": token,
	})
}


func (handler *UserHandler) GetAllUsers(writer http.ResponseWriter, req *http.Request) {
	users, err := handler.UserService.GetAllUsers()
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Failed to fetch users"})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(users)
}
