package main

import (
	"api_user_service_booking/api"
	"api_user_service_booking/config"
	"api_user_service_booking/pkg/logger"
	"api_user_service_booking/services"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	// // Connect postgres for casbin
	//psqlCon := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`, "localhost", 5432, "doston", "doston", "user_service")
	//db, err := gormadapter.NewAdapter("postgres", psqlCon, true)
	//if err != nil {
	//	log.Error("error gormadapter", logger.Error(err))
	//	return
	//}

	// //connect to postgres ACL
	//enforcer, err := casbin.NewEnforcer("config/auth.conf", db)
	//if err != nil {
	//	log.Error("error NewEnforcer", logger.Error(err))
	//	return
	//}

	// Casbin connect with .csv file
	fileAdapter := fileadapter.NewAdapter("./config/auth.csv")

	enforcer, err := casbin.NewEnforcer("./config/auth.conf", fileAdapter)
	if err != nil {
		log.Error("NewEnforcer error", logger.Error(err))
		return
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		Enforcer:       enforcer,
		ServiceManager: serviceManager,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}
}
