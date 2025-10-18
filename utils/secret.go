package utils

import (
	"soccer-team-api/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(user model.User) (string, error) {
	env := LoadEnv()
	expDuration, err := time.ParseDuration(env.JWTExpiresIn)
	if err != nil {
		return "", err
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     jwt.NewNumericDate(time.Now().Add(expDuration)),
	}).SignedString([]byte(env.JWTSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}
func ParseJWT(tokenStr string) (*jwt.Token, error) {
	env := LoadEnv()
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(env.JWTSecret), nil
	})
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
