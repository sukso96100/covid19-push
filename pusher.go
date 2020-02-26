package main

import (
	"log"
	"os"

	"github.com/alexandrevicenzi/go-sse"
)

var Pusher *sse.Server

func InitPusher() {
	Pusher = sse.NewServer(&sse.Options{
		// CORS headers
		Headers: map[string]string{
			"Connection":                   "keep-alive",
			"Cache-Control":                "no-transform, no-cache",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Keep-Alive,X-Requested-With,Cache-Control,Content-Type,Last-Event-ID",
		},
		Logger: log.New(os.Stdout, "go-sse: ", log.Ldate|log.Ltime|log.Lshortfile),
	})
}
