package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/nurovic/hmall/models"
	"github.com/nurovic/hmall/pkg"
	"github.com/nurovic/hmall/store"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) 
	defer cancel()
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := store.CreateUser(ctx,user); err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode("Created User")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Geçersiz giriş verisi", http.StatusBadRequest)
		return
	}

	storedUser, err := store.GetUserByEmail(ctx, user.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Kullanıcı bulunamadı: %v", err), http.StatusNotFound)
		return
	}
    checkPassword := pkg.CheckPasswordHash(user.Password, storedUser.Password)
    if !checkPassword {
        http.Error(w, "Kullanıcı veya Şifre Yanlış", http.StatusUnauthorized)
		return
	}

    token, err := pkg.GenerateJWT(storedUser.Email)
    if err != nil {
        log.Printf("Error generating token: %v", err)
        http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
        return
    }


    response := map[string]string{
        "email": storedUser.Email,
        "token": token,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
    
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) 
	defer cancel()
    vars := mux.Vars(r)
    userID := vars["id"]

    user, err := store.GetUserByID(ctx, userID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}


func GetMe(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Token bulunamadı", http.StatusUnauthorized)
        return
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        http.Error(w, "Geçersiz token formatı", http.StatusUnauthorized)
        return
    }
    tokenString := parts[1]
    email, err := pkg.ValidateJWT(tokenString)
    if err != nil {
        if validationErr, ok := err.(*jwt.ValidationError); ok {
            if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
                http.Error(w, "Token süresi dolmuş", http.StatusUnauthorized)
                return
            }
        }

        http.Error(w, fmt.Sprintf("Token doğrulama hatası: %v", err), http.StatusUnauthorized)
        return
    }

    response := map[string]string{
        "email": email,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

