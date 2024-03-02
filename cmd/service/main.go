package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/JoKr4/gpiod2go/pkg/gpiod"
	"github.com/shamaton/msgpack/v2"
)

// type relaisState uint8

// const (
// 	unknown relaisState = iota
// 	off     relaisState
// 	on      relaisState
// 	erro    relaisState
// )

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

	gpiochip0 := gpiod.NewDevice(useDevice)
	err := gpiochip0.Open()
	if err != nil {
		log.Println(err)
		return
	}
	defer gpiochip0.Close()
	log.Println("successfully opened device")

	route := "/relais/states"

	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		result := make([]bool, len(useGPIOsForRelais))
		for i, offset := range useGPIOsForRelais {
			result[i] = false // TODO unknown
			currentValue, err := gpiochip0.GetLineValue(offset)
			if err != nil {
				log.Println(err)
				continue
			}
			if currentValue == gpiod.LineValueActive {
				result[i] = true
			} else {
				result[i] = false
			}
		}
		w.Header().Set("Content-Type", "application/msgpack")
		d, err := msgpack.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(d)
		if err != nil {
			log.Println(err)
		}
		log.Println("responded state request of", r.RemoteAddr)
	})

	for i, offset := range useGPIOsForRelais {

		// NOTE this is not needed for following call "gpiochip0.GetLineDirection(offset)"
		// but to be able to call set-commands on gpiochip0
		err = gpiochip0.AddLine(offset)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("using line", offset)

		direction, err := gpiochip0.GetLineDirection(offset)
		if err != nil {
			log.Println(err)
			continue
		}

		if direction != gpiod.LineDirectionOutput {
			err = gpiochip0.SetLineDirection(offset, gpiod.LineDirectionOutput)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("successfully set line direction")
		} else {
			log.Println("existing output line", offset)
		}

		route = fmt.Sprintf("/relais%d/on", i)

		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			err := gpiochip0.SetLineValue(offset, gpiod.LineValueActive)
			if err != nil {
				log.Println(err)
			}
			// TODO post current state
		})

		route = fmt.Sprintf("/relais%d/off", i)

		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			err := gpiochip0.SetLineValue(offset, gpiod.LineValueInactive)
			if err != nil {
				log.Println(err)
			}
			// TODO post current state
		})

	}

	http.ListenAndServe(":8090", nil)
}
