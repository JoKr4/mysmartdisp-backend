package main

import (
	"log"
	"net/http"

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

	http.ListenAndServe(":8090", nil)
}
