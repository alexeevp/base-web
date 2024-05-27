package main

import (
	"bms/webserver"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Instance struct {
	id      int
	name    string
	status  string
	created time.Time
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	webserver.Run()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-s
	webserver.Shutdown()
}
