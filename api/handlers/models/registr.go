package models

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/smtp"
	"regexp"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/redis/go-redis/v9"
)

type UserRegister struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Code      string `json:"code"`
}

type ResponseMessage struct {
	Content string `json:"content"`
}

type VerifyUserRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type VerifyUserResponse struct {
	UserInfo User `json:"userinfo"`
}

type ResponseUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Code      string `json:"code"`
}

type UserValidationStauts struct {
	Status bool `json:"status"`
}

type CheckUser struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (rm *UserRegister) Validate() error {
	return validation.ValidateStruct(
		rm,
		validation.Field(&rm.Email, validation.Required, is.Email),
		validation.Field(
			&rm.Password,
			validation.Required,
			validation.Length(8, 30),
			validation.Match(regexp.MustCompile("[a-z]|[A-Z][1-9]")),
		),
	)
}

func GenerateCode(rdb *redis.Client, user UserRegister) string {
	code := strconv.Itoa(rand.Int())[:6]
	userByte, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err, "error marhshalling user to json")
	}
	_, err = rdb.Set(context.Background(), code, userByte, time.Minute*3).Result()
	if err != nil {
		fmt.Println(err, "error saving code to redis")
		return code
	}
	return code
}

func SendCode(email, code string) {
	from := "dostonshernazarov2001@gmail.com"
	password := "yzri faon zuix pldt"

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(code)

	auth := smtp.PlainAuth("Verification Code for join Golang Community", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
}
