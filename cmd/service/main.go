package main

import (
	"log"
	"net/http"

	"github.com/JoKr4/gpiod2go/pkg/gpiod"
)

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

	http.HandleFunc("/relais1/on", func(w http.ResponseWriter, r *http.Request) {
		err := d.SetLineValue(uint(useOffset), gpiod.LineValueActive)
		if err != nil {
			log.Println(err)
		}
	})
	http.HandleFunc("/relais1/off", func(w http.ResponseWriter, r *http.Request) {
		err := d.SetLineValue(uint(useOffset), gpiod.LineValueInactive)
		if err != nil {
			log.Println(err)
		}
	})
	http.ListenAndServe(":8090", nil)
}
