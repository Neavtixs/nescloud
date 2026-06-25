package main

import (
	"nescloud/backend-app/configs"
	"nescloud/backend-app/internal"
	"nescloud/backend-app/internal/route"
)

func main() {
	configs.LoadEnv()

	db := configs.GetConnection()
	rdb := configs.NewAccess()
	validate := configs.NewValidator()
	log := configs.NewLogger()
	app := configs.NewGin()

	internal.InitDependencies(db, rdb, validate, log)
	route.SetupRoute(app)

	app.Run(":8080")
}
