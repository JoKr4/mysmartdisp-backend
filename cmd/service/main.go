package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/JoKr4/mysmartdisp-backend/internal"
)

func main() {

	err := internal.ServeGPIOD()
	if err != nil {
		log.Println(err)
		return
	}

	err = internal.ServeLIRCD()
	if err != nil {
		log.Println(err)
		return
	}

	go http.ListenAndServe(":8090", nil)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	internal.CleanupGPIOD()
}
