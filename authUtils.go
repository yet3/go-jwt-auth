package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserJwtClaims struct {
	UserId int `json:"userId"`
	jwt.RegisteredClaims
}

func generateJwtToken(key string, expIn time.Duration, userId int) (string, error) {
	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, UserJwtClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	accessToken, err := jwtAccessToken.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return accessToken, nil
}

func generateAccessToken(userId int) (string, error) {
	return generateJwtToken(os.Getenv("ACCESS_TOKEN"), time.Hour, userId)
}

func generateRefreshToken(userId int) (string, error) {
	return generateJwtToken(os.Getenv("REFRESH_TOKEN"), 7*24*time.Hour, userId)
}

func GenerateAuthTokens(userId int) (string, string, error) {
	accessToken, err := generateAccessToken(userId)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken(userId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func validateAuthToken(token string, key string) (*UserJwtClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &UserJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, err
	} else if claims, ok := parsed.Claims.(*UserJwtClaims); ok {
		return claims, err
	}

	return nil, err
}

func ValidateAccessToken(token string) (*UserJwtClaims, error) {
	return validateAuthToken(token, os.Getenv("ACCESS_TOKEN"))
}

func ValidateRefreshToken(token string) (*UserJwtClaims, error) {
	return validateAuthToken(token, os.Getenv("REFRESH_TOKEN"))
}
