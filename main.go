package main

import (
	"fmt"
	"net/http"
	"os"

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

	// Create a new Mux and set the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/collect", Collect)
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("=====Server is now up and running=====")

	http.ListenAndServe(":8080", mux)
}
