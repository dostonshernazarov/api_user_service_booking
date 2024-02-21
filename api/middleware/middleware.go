package middleware

import (
	"api_user_service_booking/api/auth"
	"api_user_service_booking/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func MiddleGetInfo(ctx *gin.Context) {
	fmt.Printf("Method: %s, \nPath: %s, \nIP: %S",
		ctx.Request.Method,
		ctx.Request.URL.Path,
		ctx.ClientIP())
	ctx.Next()
}

func MiddleCheckAdmin(ctx *gin.Context) {
	key := ctx.GetHeader("Authorization")

	if key != "adminsecure" {
		ctx.AbortWithStatusJSON(403, gin.H{
			"Message": "You have no access to this page",
		})
		return
	}

	ctx.Next()
}

func Auth(ctx *gin.Context) {
	if ctx.Request.URL.Path != "/v1/users/login" {
		ctx.Next()
		return
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized request",
		})
		return
	}
	fmt.Printf(token)

	cfg := config.Load()
	key := cfg.SignInKey
	claims, err := auth.ExtractClaim(token, []byte(key))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}

	if cast.ToString(claims["role"]) != "user" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "have no access to this path",
		})
		return
	}

	ctx.Next()
}
