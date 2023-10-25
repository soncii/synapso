package test

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/ybkuroki/go-webapp-project-template/config"
	"github.com/ybkuroki/go-webapp-project-template/logger"
	"github.com/ybkuroki/go-webapp-project-template/migration"
	"github.com/ybkuroki/go-webapp-project-template/repository"
	"github.com/ybkuroki/go-webapp-project-template/session"
)

// Prepare func is to prepare for unit test.
func Prepare() *echo.Echo {
	e := echo.New()

	conf := &config.Config{}
	conf.Database.Dialect = "sqlite3"
	conf.Database.Host = "file::memory:?cache=shared"
	conf.Database.Migration = true
	conf.Extension.MasterGenerator = true
	conf.Extension.SecurityEnabled = false
	conf.Log.Format = "${time_rfc3339} [${level}] ${remote_ip} ${method} ${uri} ${status}"
	conf.Log.Level = 1
	config.SetConfig(conf)

	logger.InitLogger(e, config.GetConfig())

	repository.InitDB()

	migration.CreateDatabase(config.GetConfig())
	migration.InitMasterData(config.GetConfig())

	session.Init(e, config.GetConfig())
	return e
}

// ConvertToString func is convert model to string.
func ConvertToString(model interface{}) string {
	bytes, _ := json.Marshal(model)
	return string(bytes)
}
