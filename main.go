package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
	"synapso/config"
	"synapso/logger"
	"synapso/migration"
	"synapso/repository"
	"synapso/transport/router"
)

func main() {
	e := echo.New()

	config.Load()
	logger.InitLogger(e, config.GetConfig())
	e.Logger.Info("Loaded this configuration : application." + *config.GetEnv() + ".yml")

	repository.InitDB()
	db := repository.GetDB()

	migration.CreateDatabase(config.GetConfig())
	migration.InitMasterData(config.GetConfig())

	router.Init(e, config.GetConfig())

	port := os.Getenv("PORT")
	fmt.Println("port:" + port)
	if err := e.Start(":" + port); err != nil {
		e.Logger.Error(err)
	}

	defer db.Close()
}
