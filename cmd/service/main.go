package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/JoKr4/gpiod2go/pkg/gpiod"
)

type response struct {
	State bool `json:"state"`
}

var mu sync.Mutex

func main() {

	useDevice := "/dev/gpiochip0"
	useOffset := 22

	d := gpiod.NewDevice(useDevice)
	err := d.Open()
	if err != nil {
		log.Println(err)
		return
	}
	defer d.Close()
	log.Println("successfully opened device")

	err = d.AddLine(uint(useOffset), gpiod.LineDirectionOutput)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("successfully added line")

	http.HandleFunc("/relais1/state", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		currentValue, err := d.GetLineValue(uint(useOffset))
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		jenc := json.NewEncoder(w)
		state := true
		if currentValue == gpiod.LineValueInactive {
			state = true
		}
		jenc.Encode(response{state})
		fmt.Fprintf(w, "Hello, %q")
	})

	http.HandleFunc("/relais1/on", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		err := d.SetLineValue(uint(useOffset), gpiod.LineValueActive)
		if err != nil {
			log.Println(err)
		}
		// TODO post current state
	})

	http.HandleFunc("/relais1/off", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		err := d.SetLineValue(uint(useOffset), gpiod.LineValueInactive)
		if err != nil {
			log.Println(err)
		}
		// TODO post current state
	})

	http.ListenAndServe(":8090", nil)
}
