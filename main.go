package main

import (
	"bms/web"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	w := web.NewWebserver(os.Getenv("JWT_SIGNKEY"))

	w.Run()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-s
	w.Off()
}
