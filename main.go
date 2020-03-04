package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sukso96100/covid19-push/database"
	"github.com/sukso96100/covid19-push/fcm"
	"github.com/sukso96100/covid19-push/tgbot"
)

func main() {
	var dbHost string

	if os.Getenv("DB_HOST") == "" {
		dbHost = os.Getenv("CLOUD_SQL_CONNECTION_NAME")
	} else {
		dbHost = os.Getenv("DB_HOST")
	}
	initErr := database.InitDatabase(
		os.Getenv("DB_PROTOCOL"),
		dbHost,
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_CHARSET"))
	if initErr != nil {
		fmt.Println("DB Init Fail")
		fmt.Printf("%w", initErr)
	}
	database.MigrateDb()
	defer database.DbConn.Close()

	fcm.InitFCMApp(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	err := tgbot.InitTgBot(
		os.Getenv("TELEGRAM_TOKEN"),
		os.Getenv("TELEGRAM_CHANNEL"))
	if err != nil {
		fmt.Println("Telegram Init Fail")
		fmt.Printf("%w", initErr)
	}

	e := echo.New()
	e.GET("/collect", Collect)
	e.POST("/subscribe/:topic", Subscribe)
	e.POST("/unsubscribe/:topic", Unubscribe)
	e.GET("/stat", CurrentStat)
	e.GET("/news", RecentNews)
	e.File("/redirect/*", "static/index.html")
	e.File("/", "static/index.html")
	//
	e.Static("/", "static")
	e.Logger.Fatal(e.Start(":8080"))
}
