package auth

import (
	"fmt"
	"peer-store/config"
	"peer-store/db"
	"peer-store/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userId uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userId,
		"user_email": email,
		"exp":        time.Now().Add(30 * time.Hour * 24).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(config.CONFIG.AppSecret)
	if err != nil {
		return "", err
	}

	tokenRecord := models.AccessToken{
		UserId:  userId,
		Token:   tokenString,
		Blocked: false,
	}

	if err := db.GetDB().Create(&tokenRecord).Error; err != nil {
		return "", err
	}

	return tokenString, nil

}

func VerifyJWT(tokenString string) (string, string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.CONFIG.AppSecret, nil
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		userID := fmt.Sprintf("%v", claims["user_id"])

		userEmail, ok := claims["user_email"].(string)

		if !ok {
			return "", "", fmt.Errorf("invalid user email claim")
		}

		if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
			return "", "", fmt.Errorf("token has expired")
		}

		return userID, userEmail, nil

	}

	return "", "", fmt.Errorf("invalid token")

}
