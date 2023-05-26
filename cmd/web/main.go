package main

import (
	_ "github.com/mjaliz/deviran/docs"
	"github.com/mjaliz/deviran/internal/config"
	"github.com/mjaliz/deviran/internal/handlers"
	"github.com/mjaliz/deviran/internal/models"
	"github.com/mjaliz/deviran/internal/routes"
)

// @title Deviran API
// @version 1.0
// @description This is the server of Deviran platform
// @license.name Apache 2.0

var app config.AppConfig

func main() {
	app.DB = models.ConnectDB()
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	routes.Routes()
}
