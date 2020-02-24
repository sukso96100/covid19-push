package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/r3labs/sse"
)

func main() {
	initErr := InitDatabase(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_CHARSET"))
	if initErr != nil {
		fmt.Printf("%w", initErr)
	}
	MigrateDb()
	defer DbConn.Close()

	InitPusher()
	Pusher.CreateStream("stat")
	Pusher.CreateStream("news")

	defer Pusher.Close()

	// Create a new Mux and set the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/updates", Pusher.HTTPHandler)
	mux.HandleFunc("/collect", Collect)

	http.ListenAndServe(":8080", mux)
}
