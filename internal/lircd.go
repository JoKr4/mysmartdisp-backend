package internal

import (
	"log"
	"net/http"
	"os/exec"
	"sync"
)

var mu sync.Mutex

func ServeLIRCD() error {

	http.HandleFunc("/irremote/amp/volup", func(w http.ResponseWriter, r *http.Request) {

		mu.Lock()
		defer mu.Unlock()

		cmd := exec.Command("irsend", "SEND_ONCE", "cambr_volup_burst", "KEY_VOLUMEUP_BURST")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
		log.Println("did irsend KEY_VOLUMEUP_BURST request of", r.RemoteAddr)
	})

	return nil
}
