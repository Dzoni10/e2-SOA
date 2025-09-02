package handler

import (
	"bytes"
	"database-example/auth"
	"database-example/dto"
	"database-example/model"
	"database-example/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

	///ISpis u konzoli vrednosti unosa fronta
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	fmt.Println("Request body:", string(bodyBytes))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	/////

	var userDTO dto.CreateUserDTO

	err := json.NewDecoder(req.Body).Decode(&userDTO)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user := model.User{
		Name:     userDTO.Name,
		Surname:  userDTO.Surname,
		Username: userDTO.Username,
		Password: userDTO.Password, // ovde lepo setuješ
		Role:     model.Role(userDTO.Role),
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
		Password string `json:"password"`
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

	userID := handler.ExtractUserIDFromToken(req)

	users, err := handler.UserService.GetAllUsersExcept(userID)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Failed to fetch users"})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(users)
}

func (handler *UserHandler) ExtractUserIDFromToken(req *http.Request) int {

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return 0 // ili -1 ako želiš da označiš grešku
	}

	// 2. Očekujemo format "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0
	}
	tokenString := parts[1]

	// 3. Parsiraj token
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_password_for_encoding_messages"), nil
	})
	if err != nil {
		return 0
	}

	// 4. Izvuci UserId iz claims-a
	if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
		return claims.UserId
	}

	return 0
}

func (handler *UserHandler) BlockUser(writer http.ResponseWriter, req *http.Request) {
	userID_str := mux.Vars(req)["id"]
	userID, err := strconv.Atoi(userID_str)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"error": "Invalid ID format"}`))
		return
	}

	blockedUser, err := handler.UserService.BlockUser(userID)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Failed to block user"})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blockedUser)
}

func (handler *UserHandler) UpdateUserProfile(writer http.ResponseWriter, req *http.Request) {
	tokenUserID := handler.ExtractUserIDFromToken(req)
	if tokenUserID == 0 {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Invalid or missing token"})
		return
	}

	user, err := handler.UserService.FindUser(tokenUserID)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": "bad user id"})
	}

	// Parse request body
	var profileDTO dto.UserProfile
	if err := json.NewDecoder(req.Body).Decode(&profileDTO); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}

	// Update user profile using service
	handler.UserService.UpdateUserProfile(&user, &profileDTO)

	// Return success response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"message": "Profile updated successfully"})
}
