package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
    log.Println("HERE",user.Password)
    checkPassword := pkg.CheckPasswordHash(user.Password, storedUser.Password)
    if !checkPassword {
        http.Error(w, "Kullanıcı veya Şifre Yanlış", http.StatusUnauthorized)
		return
	}

    token, err := generateJWT(storedUser.Email)
    if err != nil {
        log.Printf("Error generating token: %v", err)
        http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
        return
    }

    log.Println("HERE Token generated:", token)

    // Return token and email in response
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

func generateJWT(email string) (string, error) {
    mongoURI := os.Getenv("SECRET_KEY")
    secretKey := mongoURI

    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}
func GetMe(w http.ResponseWriter, r *http.Request) {
    // Authorization header'dan JWT token alınır
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Token bulunamadı", http.StatusUnauthorized)
        return
    }

    // Bearer prefix kontrolü ve token ayrıştırma
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        http.Error(w, "Geçersiz token formatı", http.StatusUnauthorized)
        return
    }
    tokenString := parts[1]

    // Token doğrulama
    email, err := validateJWT(tokenString)
    if err != nil {
        // Token süresi dolmuşsa kontrol et
        if validationErr, ok := err.(*jwt.ValidationError); ok {
            if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
                http.Error(w, "Token süresi dolmuş", http.StatusUnauthorized)
                return
            }
        }

        // Diğer doğrulama hatalarını döndür
        http.Error(w, fmt.Sprintf("Token doğrulama hatası: %v", err), http.StatusUnauthorized)
        return
    }

    // Kullanıcı bilgilerini döndür
    response := map[string]string{
        "email": email,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}


// func validateJWT(tokenString string) (string, error) {
//     secretKey := "yourSecretKey"

//     token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//         if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//             return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//         }
//         return []byte(secretKey), nil
//     })

//     if err != nil {
//         // Token süresi dolmuşsa özel bir hata ile döner
//         if validationErr, ok := err.(*jwt.ValidationError); ok && validationErr.Errors&jwt.ErrTokenExpired != 0 {
//             return "", jwt.ErrTokenExpired
//         }
//         return "", err
//     }

//     if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//         email, ok := claims["email"].(string)
//         if !ok {
//             return "", fmt.Errorf("email bilgisi bulunamadı")
//         }
//         return email, nil
//     }

//     return "", fmt.Errorf("geçersiz token")
// }

func validateJWT(tokenString string) (string, error) {
    mongoURI := os.Getenv("SECRET_KEY")
    secretKey := mongoURI

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        // Eğer token süresi dolmuşsa `jwt.ValidationErrorExpired` hata kodu ile kontrol edebilirsiniz
        if validationErr, ok := err.(*jwt.ValidationError); ok {
            if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
                return "", fmt.Errorf("token süresi dolmuş")
            }
        }
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        email, ok := claims["email"].(string)
        if !ok {
            return "", fmt.Errorf("token'dan e-posta alınamadı")
        }
        return email, nil
    }

    return "", fmt.Errorf("geçersiz token")
}
