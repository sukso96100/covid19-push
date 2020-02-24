package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/r3labs/sse"
)

var Pusher *sse.Server

func InitPusher() {
	Pusher = sse.New()
}
