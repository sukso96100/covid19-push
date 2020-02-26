package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	initErr := InitDatabase(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_CHARSET"))
	if initErr != nil {
		fmt.Println("DB Init Fail")
		fmt.Printf("%w", initErr)
	}
	MigrateDb()
	defer DbConn.Close()

	InitPusher()

	defer Pusher.Shutdown()

	// Create a new Mux and set the handler
	mux := http.NewServeMux()
	mux.Handle("/updates/", Pusher)
	mux.HandleFunc("/collect", Collect)
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("=====Server is now up and running=====")

	http.ListenAndServe(":8080", mux)
}
