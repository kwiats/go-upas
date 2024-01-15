package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

func CreateJWT(userID, username string) (*TokenResponse, error) {
	JWT_secret_token := os.Getenv("JWT_SECRET")
	tokenResponse := new(TokenResponse)
	claims := &jwt.MapClaims{
		"exp":           time.Now().Add(time.Minute * 15).Unix(),
		"user_id":       userID,
		"user_username": username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(JWT_secret_token))
	if err != nil {
		return nil, err
	}

	tokenResponse.AuthToken = signedToken

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = 1
	refreshTokenClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte(JWT_secret_token))
	if err != nil {
		return nil, err
	}

	tokenResponse.RefreshToken = signedRefreshToken
	return tokenResponse, nil
}

func ValidateJWT(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

}
