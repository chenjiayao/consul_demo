package routes

import (
	"upstream_service/controllers"
	"upstream_service/globals"
)

var (
	indexController = &controllers.IndexController{}
)

func LoadRoute() {

	globals.Engine.GET("/health", indexController.Health)
	globals.Engine.GET("/v1", indexController.TestV1)
	globals.Engine.GET("/v2", indexController.TestV2)
}
