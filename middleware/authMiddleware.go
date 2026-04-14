package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
	Id   int
	Type string
	jwt.RegisteredClaims
}

func Refresh(tokenString string) (string, string, int, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return "", "", -1, err
	}
	if !token.Valid {
		return "", "", -1, errors.New("Not valid refresh token, please auth")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", "", -1, errors.New("No claims")
	}
	id := claims.Id
	accessToken, refreshToken, err := CreateToken(id)
	return accessToken, refreshToken, id, err
}

func GetId(tokenString string) (int, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return -1, err
	}
	if !token.Valid {
		return -1, errors.New("Not valid refresh token, please auth")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return -1, errors.New("No claims")
	}
	id := claims.Id
	return id, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Not valid signing method")
		}

		godotenv.Load()
		return []byte(os.Getenv("SECRET_WORD")), nil
	})
}

func CreateToken(id int) (string, string, error) {
	godotenv.Load()
	secretWord := os.Getenv("SECRET_WORD")
	secretKey := []byte(secretWord)
	accessClaims := Claims{
		Id:   id,
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}
	refreshClaims := Claims{
		Id:   id,
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}
	return accessString, refreshString, nil
}

func AuthMiddleware(next http.HandlerFunc, logger zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			logger.Info("Got hhtp options request")
			w.WriteHeader(http.StatusOK)
			return
		}
		logger.Info("Getting access token cookie")
		cookie, err := r.Cookie("accessToken")
		if err != nil {
			logger.Error("Failed to get access token")
			errdto := ErrorDTO{
				Message: err.Error(),
				Time:    time.Now(),
			}
			http.Error(w, errdto.ToString(), http.StatusUnauthorized)
			return
		}
		logger.Info("Getting value of access token cookie")
		tokenString := cookie.Value
		logger.Info("Parsing token")
		token, err := parseToken(tokenString)
		if err != nil {
			logger.Error("Failed to parse token")
			errdto := ErrorDTO{
				Message: err.Error(),
				Time:    time.Now(),
			}
			http.Error(w, errdto.ToString(), http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			logger.Error("Token is not valid")
			errdto := ErrorDTO{
				Message: "Invalid token",
				Time:    time.Now(),
			}
			http.Error(w, errdto.ToString(), http.StatusUnauthorized)
			return
		}
		logger.Info("Next function")
		next(w, r)
	}
}
