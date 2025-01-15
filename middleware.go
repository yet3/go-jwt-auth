package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthorizationHeader struct {
	Authorization string `header:"Authorization" binding:"required,startswith=Bearer "`
}

func Authorized(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var headerData AuthorizationHeader

		if err := ctx.ShouldBindHeader(&headerData); err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Accessing unauthorized route"})
			return
		}

		accessToken := headerData.Authorization[7:]
		claims, err := ValidateAccessToken(accessToken)

		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invlaid access token"})
			return
		}

		var userId int
		err = db.QueryRow("SELECT id FROM users WHERE id = ? AND accessToken = ?", claims.UserId, accessToken).Scan(&userId)
		if err != nil {
			fmt.Println(err)
			if err == sql.ErrNoRows {
				ctx.Abort()
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Invlaid access token"})
				return
			}
			ctx.Abort()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating access token"})
			return
		}

    fmt.Println("AUTHORIZED")
		ctx.Set("userId", userId)
		ctx.Next()
	}
}
