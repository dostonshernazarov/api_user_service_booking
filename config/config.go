package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment string // develop, staging, production

	UserServiceHost string
	UserServicePort int

	//PostServiceHost string
	//PostServicePort int

	//CommentServiceHost string
	//CommentServicePort int
	//
	//LikeServiceHost string
	//LikeServicePort int

	// context timeout in seconds
	AccessTokenTimout int // MINUTES
	SignInKey         string

	CtxTimeout int

	LogLevel string
	HTTPPort string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "localhost"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 1111))

	//c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "127.0.0.1"))
	//c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_PORT", 2222))
	//
	//c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "127.0.0.1"))
	//c.CommentServicePort = cast.ToInt(getOrReturnDefault("COMMENT_SERVICE_PORT", 3333))
	//
	//c.LikeServiceHost = cast.ToString(getOrReturnDefault("LIKE_SERVICE_HOST", "127.0.0.1"))
	//c.LikeServicePort = cast.ToInt(getOrReturnDefault("LIKE_SERVICE_PORT", 4444))

	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))
	c.AccessTokenTimout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 300))
	c.SignInKey = cast.ToString(getOrReturnDefault("SIGN_IN_KEY", "test_key"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
