package main

import (
	"follooow-be/configs"
	"follooow-be/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// run database
	configs.ConnectDB()

	// initialize Cloudinary
	configs.InitCloudinary()

	// routes
	routes.InfluencerRoute(e)
	routes.NewsRoute(e)
	routes.GalleriesRoute(e)
	routes.UserRoute(e)
	routes.MediaRoute(e)

	e.Logger.Fatal(e.Start(":20223"))
}
