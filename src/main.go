package main

import (
	log "github.com/sirupsen/logrus"
	"portfolio-cms-server/config"
	"portfolio-cms-server/database"
	"portfolio-cms-server/server"
	utils2 "portfolio-cms-server/utils"
)

func init() {
	app, err := config.Init()
	if err != nil {
		utils2.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error on config initialization")
		return
	}
	if app.AppEnv == "LOC" {
		utils2.PrettyPrint(app)
	}
	database.Init(app.DBHosts, app.DBUsername, app.DBPassword, app.DBPort, app.DBName)
}

func main() {
	server.Run()
}
