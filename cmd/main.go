package main

import (
	"api_user_service_booking/api"
	"api_user_service_booking/config"
	"api_user_service_booking/pkg/logger"
	rbmq "api_user_service_booking/queue/rabbitmq/producermq"
	"api_user_service_booking/services"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	writer, err := rbmq.NewRabbitMQProducer("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Error("error rabbit mq", logger.Error(err))
		//	return
	}
	defer writer.Close()

	//// Connect postgres for casbin
	//psqlCon := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`, "localhost", 5432, "doston", "doston", "user_service")
	//db, err := gormadapter.NewAdapter("postgres", psqlCon, true)
	//if err != nil {
	//	log.Error("error gormadapter", logger.Error(err))
	//	return
	//}
	//
	////connect to postgres casbin
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

	// casbin with CSV -------------------------------------------
	casbinEnforcer, err := casbin.NewEnforcer(cfg.AuthConfigPath, cfg.CSVFilePath)
	if err != nil {
		log.Fatal("casbin enforcer error", logger.Error(err))
		return
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil { // gormadapter "github.com/casbin/gorm-adapter/v3"
		log.Fatal("casbin error load policy", logger.Error(err))
		return
	}

	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl).AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl).AddMatchingFunc("keyMatch3", util.KeyMatch3)

	//writer, err := producer.NewKafkaProducerInit([]string{"localhost:9092"})
	//if err != nil {
	//	log.Error("NewKafkaProducerInit: %v", logger.Error(err))
	//}
	//
	//err = writer.ProduceMessage("test-topic", []byte("\nthis message has come from produce"))
	//if err != nil {
	//	log.Fatal("failed to run  ProduceMessage", logger.Error(err))
	//}
	//
	//defer writer.Close()

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		Enforcer:       enforcer,
		ServiceManager: serviceManager,
		Writer:         writer,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}
}
