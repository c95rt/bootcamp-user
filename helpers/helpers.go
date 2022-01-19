package helpers

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"

	"github.com/c95rt/bootcamp-user/grpc/models"
)

const (
	bcryptCost int = 14
)

func CheckPassword(hash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func GenerateToken(user *models.User, jwtSecret string) (string, error) {
	claims := struct {
		User map[string]interface{} `json:"u"`
		jwt.StandardClaims
	}{
		map[string]interface{}{
			"i":         user.ID,
			"email":     user.Email,
			"lastname":  user.Lastname,
			"firstname": user.Firstname,
		},
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}
