package main

import (
	_ "github.com/mjaliz/deviran/docs"
	"github.com/mjaliz/deviran/internal/config"
	"github.com/mjaliz/deviran/internal/handlers"
	"github.com/mjaliz/deviran/internal/models"
	"github.com/mjaliz/deviran/internal/routes"
	"github.com/mjaliz/deviran/internal/utils"
)

// @title Deviran API
// @version 1.0
// @description This is the server of Deviran platform
// @license.name Apache 2.0

var app config.AppConfig

func main() {
	rds, ctx := utils.RedisInit()
	app.RedisClient = rds
	app.RedisCtx = ctx
	app.DB = models.ConnectDB()
	repo := handlers.NewRepo(&app)
	utilsRepo := utils.NewRepo(&app)
	handlers.NewHandlers(repo)
	utils.NewUtils(utilsRepo)
	routes.Routes()
}
