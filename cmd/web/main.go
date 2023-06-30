package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mjaliz/deviran/docs"
	"github.com/mjaliz/deviran/internal/initializers"
	customMiddleware "github.com/mjaliz/deviran/internal/middleware"
	"github.com/mjaliz/deviran/internal/routes"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
)

// @title Deviran API
// @version 1.0
// @description This is the server of Deviran platform
// @license.name Apache 2.0

func init() {
	err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(initializers.Config)
	initializers.ConnectRedis(initializers.Config)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	userGroup := e.Group("/user")
	routes.UserSubRoutes(userGroup)

	courseGroup := e.Group("/course", customMiddleware.DeserializeUser)
	routes.CourseSubRoutes(courseGroup)

	e.Logger.Info(e.Start(":1323"))
}
