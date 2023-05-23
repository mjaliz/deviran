package main

import (
	_ "github.com/mjaliz/deviran/docs"
	"github.com/mjaliz/deviran/internal/config"
)

// @title Deviran API
// @version 1.0
// @description This is the server of Deviran platform
// @license.name Apache 2.0

var app config.AppConfig

func main() {
	app.DB = nil
	routes(&app)
}
