package auth

import (
	"encoding/json"

	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("JWTSecret")

func CreateToken(id string, email string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    id,
			"email": email,
			"role":  role,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if cookie.Value == "" {
			json.NewEncoder(w).Encode("User is not logged in")
			return
		}
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		e := VerifyToken(cookie.Value)
		if e != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Invalid token")
			return
		}
		fmt.Println("You are authorised")
		next.ServeHTTP(w, r)
	})
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func VerifyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if cookie.Value == "" {
			json.NewEncoder(w).Encode("User is not logged in")
			return
		}
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		role := GetRole(cookie.Value)
		if role != "admin" {
			http.Error(w, "Forbidden access, you are not admin", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetRole(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err.Error()
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Access claims
		role := claims["role"]
		return role.(string)
	} else {
		return "role not present"
	}
}

func GetUserID(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err.Error()
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Access claims
		userID := claims["id"]
		return userID.(string)
	} else {
		return "user not present"
	}
}
