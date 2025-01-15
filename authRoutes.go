package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

type AuthRoutes struct {
	router *gin.Engine
	db     *sql.DB
}

func NewAuthRoutes(router *gin.Engine, db *sql.DB) *AuthRoutes {
	routes := &AuthRoutes{router, db}

	group := router.Group("/auth")

	group.POST("/sign-in", routes.signInHandler)
	group.POST("/sign-up", routes.signUpHandler)
	group.POST("/refresh-token", routes.refreshToken)
	group.GET("/sign-in-with-token", Authorized(db), routes.signInWithTokenHandler)
	group.GET("/sign-out", Authorized(db), routes.signOutHandler)

	return routes
}

func (r *AuthRoutes) signInWithTokenHandler(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	ctx.JSON(http.StatusOK, gin.H{"userId": userId.(int)})
}

type SignInUserData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,ascii,min=6,max=24"`
}

func (r *AuthRoutes) signInHandler(ctx *gin.Context) {
	var userData SignInUserData

	err := ctx.ShouldBindBodyWith(&userData, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userId int
	var userPassword string
	err = r.db.QueryRow("SELECT id, password FROM users WHERE email = ?", userData.Email).Scan(&userId, &userPassword)

	if err == sql.ErrNoRows {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No user found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(userData.Password))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid password"})
		return
	}

	accessToken, refreshToken, err := GenerateAuthTokens(userId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing-in"})
		return
	}

	_, err = r.db.Exec(`
    UPDATE users SET accessToken = ?, refreshToken = ? WHERE id = ?
  `, accessToken, refreshToken, userId)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing-in"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userId": userId, "accessToken": accessToken, "refreshToken": refreshToken})
}

type SignUpUserData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,ascii,min=6,max=24"`
}

func (r *AuthRoutes) signUpHandler(ctx *gin.Context) {
	var userData SignUpUserData

	err := ctx.ShouldBindBodyWith(&userData, binding.JSON)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if r.db.QueryRow("SELECT id FROM users WHERE email = ?", userData.Email).Scan() != sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is already used!"})
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 12)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating new user"})
		return
	}

	res, err := r.db.Exec(`
    INSERT INTO users (email, password) VALUES (?, ?)
  `, userData.Email, string(hashedBytes))

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating new user"})
		return
	}

	userId, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retriving user id"})
		return
	}

	accessToken, refreshToken, err := GenerateAuthTokens(int(userId))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing-in"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userId": userId, "accessToken": accessToken, "refreshToken": refreshToken})
}

type RefreshTokenData struct {
	RefreshToken string `json:"refreshToken" binding:"required,jwt"`
}

func (r *AuthRoutes) refreshToken(ctx *gin.Context) {
	var userData RefreshTokenData

	err := ctx.ShouldBindBodyWith(&userData, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = ValidateRefreshToken(userData.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid refresh token"})
		return
	}

	var userId int
	err = r.db.QueryRow("SELECT id FROM users WHERE refreshToken = ?", userData.RefreshToken).Scan(&userId)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "RefreshToken is invalid"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error refreshing access token"})
		return
	}

	accessToken, refreshToken, err := GenerateAuthTokens(userId)

	_, err = r.db.Exec(`
    UPDATE users SET accessToken = ?, refreshToken = ? WHERE id = ?
  `, accessToken, refreshToken, userId)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error refreshing access token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userId": userId, "accessToken": accessToken, "refreshToken": refreshToken})
}

type SignOutData struct {
	Authorization string `header:"Authorization" binding:"required,startswith=Bearer "`
}

func (r *AuthRoutes) signOutHandler(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	_, err := r.db.Exec("UPDATE users SET refreshToken = NULL, accessToken = NULL WHERE id = ?", userId.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing-out"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userId": userId})

}
