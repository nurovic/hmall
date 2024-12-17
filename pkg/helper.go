package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateJWT(tokenString string) (string, error) {
    mongoURI := os.Getenv("SECRET_KEY")
    secretKey := mongoURI

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
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

func GenerateJWT(email string) (string, error) {
    mongoURI := os.Getenv("SECRET_KEY")
    secretKey := mongoURI

    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}