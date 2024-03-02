package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/shamaton/msgpack/v2"

	"github.com/JoKr4/gpiod2go/pkg/gpiod"
)

func main() {

	var mu sync.Mutex

	var useGPIOsForRelais [6]uint = [...]uint{
		17,
		22,
		23,
		24,
		25,
		27,
	}

	useDevice := "/dev/gpiochip0"

	d := gpiod.NewDevice(useDevice)
	err := d.Open()
	if err != nil {
		log.Println(err)
		return
	}
	defer d.Close()
	log.Println("successfully opened device")

	for i, offset := range useGPIOsForRelais {

		err = d.AddLine(offset, gpiod.LineDirectionOutput)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("successfully added line", offset)

		route := fmt.Sprintf("/relais%d/state", i)

		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			currentValue, err := d.GetLineValue(offset)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/msgpack")
			state := true
			if currentValue == gpiod.LineValueInactive {
				state = false
			}
			d, err := msgpack.Marshal(state)
			if err != nil {
				log.Println(err)
			}
			_, err = w.Write(d)
			if err != nil {
				log.Println(err)
			}
		})

		route = fmt.Sprintf("/relais%d/on", i)

		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			err := d.SetLineValue(offset, gpiod.LineValueActive)
			if err != nil {
				log.Println(err)
			}
			// TODO post current state
		})

		route = fmt.Sprintf("/relais%d/off", i)

		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			err := d.SetLineValue(offset, gpiod.LineValueInactive)
			if err != nil {
				log.Println(err)
			}
			// TODO post current state
		})

	}

	http.ListenAndServe(":8090", nil)
}
