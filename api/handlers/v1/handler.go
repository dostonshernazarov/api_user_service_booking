package v1

import (
	"api_user_service_booking/api/handlers/models"
	"api_user_service_booking/api/tokens"
	"api_user_service_booking/config"
	"api_user_service_booking/pkg/logger"
	"api_user_service_booking/queue/kafka/producer"
	"api_user_service_booking/services"
	"github.com/casbin/casbin/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	jwtHandler     tokens.JwtHandler
	enforcer       *casbin.Enforcer
	writer         *producer.KafkaProducer
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	jwtHandler     tokens.JwtHandler
	Enforcer       *casbin.Enforcer
	Writer         *producer.KafkaProducer
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		jwtHandler:     c.jwtHandler,
		enforcer:       c.Enforcer,
		writer:         c.Writer,
	}
}

func handleBadRequestWithErrorMessage(c *gin.Context, l logger.Logger, err error, message string) bool {
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardErrorModel{
			Error: models.Error{
				Message: "Incorrect data supplied",
			},
		})
		l.Error(message, logger.Error(err))
		return true
	}
	return false
}
