package main

import (

	"github.com/r3labs/sse"
)

var Pusher *sse.Server

func InitPusher() {
	Pusher = sse.New()
	Pusher.CreateStream("stat")
	Pusher.CreateStream("news")
}
