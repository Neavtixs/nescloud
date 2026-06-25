package main

import (
	"nescloud/backend-app/configs"
	"nescloud/backend-app/internal"
)

func main() {
	configs.LoadEnv()

	db := configs.GetConnection()
	rdb := configs.NewAccess()
	validate := configs.NewValidator()
	log := configs.NewLogger()
	app := configs.NewGin()

	routeHandler := internal.InitDependencies(db, rdb, validate, log)
	routeHandler.SetupRoute(app)

	app.Run(":8080")
}
