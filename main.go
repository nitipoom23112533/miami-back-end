package main

import (
	"log"
	"miami-back-end/api"
	"miami-back-end/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"miami-back-end/mg"
)

func main()  {
	err := db.InitDB(db.Config{
		Username:     "Username",
		Password:     "Password",
		Server:       "Server",
		DatabaseName: "DatabaseName",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer db.DB.Close()

 	mg.InitMailGunClient()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())

	api.InitRoutes(e)

	// เริ่ม server
	// e.Logger.Fatal(e.Start(":1323"))
	e.Logger.Fatal(e.Start(":8080"))



}