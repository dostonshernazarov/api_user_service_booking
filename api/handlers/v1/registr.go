package v1

import (
	"api_user_service_booking/api/handlers/models"
	"api_user_service_booking/api/tokens"
	pbu "api_user_service_booking/genproto/user_proto"
	l "api_user_service_booking/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
)

// Registr
// @Summary Registr
// @Description Registr - Api for registring users
// @Tags registr
// @Accept json
// @Produce json
// @Param registr body models.UserRegister true "UserDetail"
// @Success 200 {object} models.ResponseUser
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/signup [post]
func (h *handlerV1) Registr(c *gin.Context) {
	var (
		body        models.UserRegister
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	body.Email = strings.TrimSpace(body.Email)
	body.Password = strings.TrimSpace(body.Password)
	body.Email = strings.ToLower(body.Email)

	existEmail, err := h.serviceManager.UserService().CheckUniqueEmail(ctx, &pbu.CheckUniqueRequest{
		Column: "email",
		Value:  body.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check email uniquess", l.Error(err))
		return
	}

	if existEmail.IsExist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email already in use, please use another email address",
		})
		h.log.Error("failed to check email unique", l.Error(err))
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	code := strconv.Itoa(rand.Int())[:6]
	body.Code = code

	userByte, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err, "error marhshalling user to json")
	}
	_, err = rdb.Set(context.Background(), body.Email, userByte, time.Minute*3).Result()
	if err != nil {
		fmt.Println(err, "error saving code to redis")
		return
	}

	from := "dostonshernazarov2001@gmail.com"
	password := "yzri faon zuix pldt"

	to := []string{
		body.Email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(code)

	auth := smtp.PlainAuth("Verification Code for registration", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}

	responsemessage := models.ResponseMessage{
		Content: "We send verification password you email",
	}

	c.JSON(http.StatusOK, responsemessage)
}

// LogIn
// @Summary LogIn User
// @Description LogIn - Api for login users
// @Tags registr
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Param password query string true "Password"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/login [get]
func (h *handlerV1) LogIn(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	email := c.Query("email")
	password := c.Query("password")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	responseUser, err := h.serviceManager.UserService().GetUserByEmail(ctx, &pbu.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect email to login",
		})
		h.log.Error("failed to get user info", l.Error(err))
		return
	}

	if responseUser.Password != password {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect password to login. Try again !!",
		})
		h.log.Error("failed to check password", l.Error(err))
		return
	}
	h.jwtHandler = tokens.JwtHandler{
		Sub:       responseUser.Id,
		Role:      responseUser.Role,
		SigninKey: h.cfg.SignInKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimout,
	}

	access, refresh, err := h.jwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "oops something went wrong!!",
		})
		h.log.Error("failed to generate JWT", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access token":  access,
		"refresh token": refresh,
		"user":          responseUser,
	})

	c.JSON(http.StatusOK, "ok")
}

// Verification
// @Summary Verification User
// @Description LogIn - Api for verification users
// @Tags registr
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Param code query string true "Code"
// @Success 200 {object} models.RegisterResponseModel
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/verify [get]
func (h *handlerV1) Verification(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	email := c.Query("email")
	code := c.Query("code")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	val, err := rdb.Get(ctx, email).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email. Try again ..",
		})
		h.log.Error("failed to get user from redis", l.Error(err))
		return
	}

	var userdetail models.UserRegister
	if err := json.Unmarshal([]byte(val), &userdetail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unmarshiling error",
		})
		h.log.Error("error unmarshalling userdetail", l.Error(err))
		return
	}
	if userdetail.Code != code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect code. Try again",
		})
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "error while generating uuid",
		})
		h.log.Error("error generate new uuid", l.Error(err))
		return
	}

	h.jwtHandler = tokens.JwtHandler{
		Sub:  id.String(),
		Iss:  "client",
		Role: "user",
		Log:  h.log,
	}

	access, refresh, err := h.jwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "error while generating jwt",
		})
		h.log.Error("error generate new jwt tokens", l.Error(err))
		return
	}

	user_id := uuid.New()

	res, err := h.serviceManager.UserService().Create(ctx, &pbu.User{
		Id:           user_id.String(),
		FirstName:    userdetail.FirstName,
		LastName:     userdetail.LastName,
		Email:        userdetail.Email,
		Password:     userdetail.Password,
		Role:         "user",
		RefreshToken: refresh,
	})

	c.JSON(http.StatusOK, &models.RegisterResponseModel{
		UserID:       id.String(),
		AccessToken:  access,
		RefreshToken: refresh,
	})

	c.JSON(http.StatusOK, res)
}

// RefreshAccessToken
// @Summary RefreshAccessToken User
// @Description Refresh token - Api for verification users
// @Tags token
// @Accept json
// @Produce json
// @Param refresh_token query string true "refreshToken"
// @Success 200 {object} models.RegisterResponseModel
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/retoken [get]
func (h *handlerV1) RefreshAccessToken(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	refreshToken := c.Query("refreshToken")

	fmt.Println(refreshToken)

	user, err := h.serviceManager.UserService().GetUserByRfshToken(context.Background(), &pbu.GetUserByRfshTokenRequest{
		Token: refreshToken,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect token. Try again",
		})
		return
	}

	h.jwtHandler = tokens.JwtHandler{
		Sub:     user.Id,
		Iss:     "client",
		Role:    user.Role,
		Log:     h.log,
		Timeout: h.cfg.AccessTokenTimout,
	}

	access, refresh, err := h.jwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "error while generating new jwt",
		})
		h.log.Error("error generate new jwt tokens", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.RegisterResponseModel{
		UserID:      user.Id,
		AccessToken: access,
	})

	h.serviceManager.UserService().UpdateUser(context.Background(), &pbu.User{
		Id:           user.Id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		Birthday:     user.Birthday,
		ImageUrl:     user.ImageUrl,
		CardNum:      user.CardNum,
		Phone:        user.Phone,
		Role:         user.Role,
		RefreshToken: refresh,
	})

}
