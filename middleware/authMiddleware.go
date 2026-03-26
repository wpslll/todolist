package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)
}

type Claims struct {
	id int
	jwt.RegisteredClaims
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Not valid signing method")
		}
		godotenv.Load()
		return []byte(os.Getenv("SECRET_WORD")), nil
	})
}

func CreateToken(id int) (string, error) {
	godotenv.Load()
	secretWord := os.Getenv("SECRET_WORD")
	secretKey := []byte(secretWord)
	claims := Claims{
		id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "You don't have token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := parseToken(tokenString)
		if err != nil {
			errdto := ErrorDTO {
				Message: err.Error(),
				Time: time.Now(),
			}
			http.Error(w, errdto.ToString(), http.StatusUnauthorized)
		}
		if !token.Valid {
			errdto := ErrorDTO {
				Message: "Invalid token",
				Time: time.Now(),
			}
			http.Error(w, errdto.ToString(), http.StatusUnauthorized)
		}
		next(w, r)
	}
}
