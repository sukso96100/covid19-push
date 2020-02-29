package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sukso96100/covid19-push/database"
	"github.com/sukso96100/covid19-push/fcm"
)

func main() {
	initErr := database.InitDatabase(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_CHARSET"))
	if initErr != nil {
		fmt.Println("DB Init Fail")
		fmt.Printf("%w", initErr)
	}
	database.MigrateDb()
	defer database.DbConn.Close()

	fcm.InitFCMApp(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	e := echo.New()
	e.GET("/collect", Collect)
	e.POST("/subscribe/:topic", Subscribe)
	e.POST("/unsubscribe/:topic", Unubscribe)
	e.GET("/stat", CurrentStat)
	e.GET("/news", RecentNews)
	e.Logger.Fatal(e.Start(":8080"))
}
