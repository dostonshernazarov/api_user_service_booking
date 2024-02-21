package auth

import (
	"api_user_service_booking/config"
	"api_user_service_booking/pkg/logger"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtHandler struct {
	Sub       string
	Iss       string
	Exp       string
	Iat       string
	Aud       []string
	Role      string
	Token     string
	SigninKey string
	Log       logger.Logger
	Timeout   int
}

func (jwtHandler *JwtHandler) GenerateJwt() (access, refresh string, err error) {
	var (
		accessToken, refreshToken *jwt.Token
		claims                    jwt.MapClaims
	)

	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accessToken.Claims.(jwt.MapClaims)
	claims["sub"] = jwtHandler.Sub
	claims["iss"] = jwtHandler.Iss
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["role"] = jwtHandler.Role
	claims["aud"] = jwtHandler.Aud

	sgKey := config.Load().SignInKey

	access, err = accessToken.SignedString([]byte(sgKey))
	if err != nil {
		jwtHandler.Log.Error("error generating access token", logger.Error(err))
		return
	}

	refresh, err = refreshToken.SignedString([]byte(sgKey))
	if err != nil {
		jwtHandler.Log.Error("error generating refresh token", logger.Error(err))
		return
	}

	return access, refresh, nil
}

//// ExtractClaims ...
//func (jwtHandler *JwtHandler) ExtractClaims() (jwt.MapClaims, error) {
//	var (
//		token *jwt.Token
//		err   error
//	)
//
//	token, err = jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error) {
//		return []byte(jwtHandler.SigninKey), nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if !(ok && token.Valid) {
//		jwtHandler.Log.Error("invalid jwt token")
//		return nil, err
//	}
//
//	return claims, nil
//}

// ExtractClaim extracts claims from given token
func ExtractClaim(tokenStr string, signingKey []byte) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		err = fmt.Errorf("invalid JWT Token")
		return nil, err
	}

	return claims, nil
}
