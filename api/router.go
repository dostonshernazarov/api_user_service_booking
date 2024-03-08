package api

import (
	_ "api_user_service_booking/api/docs" // swag
	v1 "api_user_service_booking/api/handlers/v1"
	casbinC "api_user_service_booking/api/middleware"
	"api_user_service_booking/config"
	"api_user_service_booking/pkg/logger"
	"api_user_service_booking/queue/rabbitmq/producermq"
	"api_user_service_booking/services"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	Enforcer       *casbin.Enforcer
	ServiceManager services.IServiceManager
	Writer         *producermq.RabbitMQProducerImpl
}

// New ...
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Enforcer:       option.Enforcer,
		Writer:         option.Writer,
	})

	api := router.Group("/v1")
	api.Use(casbinC.CheckCasbinPermission(option.Enforcer, option.Conf))
	// register
	api.POST("/users/signup", handlerV1.Registr)
	api.GET("/users/verify", handlerV1.Verification)
	api.GET("/users/login", handlerV1.LogIn)

	api.GET("/users/retoken", handlerV1.RefreshAccessToken)

	//user
	api.GET("/users/:id", handlerV1.GetUser)
	api.PUT("/users/:id", handlerV1.UpdateUser)
	api.DELETE("/users/:id", handlerV1.DeleteUser)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
