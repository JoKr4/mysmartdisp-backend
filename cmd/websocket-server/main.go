package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/JoKr4/gpiod2go/pkg/gpiod"
)

func main() {
	log.SetFlags(0)

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("please provide an address to listen on as the first argument")
	}

	useGPIOsForRelais := make(map[int]string)
	useGPIOsForRelais[17] = ""
	useGPIOsForRelais[22] = ""
	useGPIOsForRelais[23] = ""
	useGPIOsForRelais[24] = ""
	useGPIOsForRelais[25] = ""
	useGPIOsForRelais[27] = ""

	useDevice := "/dev/gpiochip0"

	gpiochip0 := gpiod.NewDevice(useDevice)
	err := gpiochip0.Open()
	if err != nil {
		return err
	}
	defer gpiochip0.Close()
	log.Println("successfully opened device")

	l, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	cs := newChatServer()
	s := &http.Server{
		Handler:      cs,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	toggleIndex := 22
	toggleState := false

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	for {
		select {
		case <-tick.C:
			err = gpiochip0.Toogle(uint(toggleIndex))
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("toggled")
			toggleState = !toggleState
			msg := make([]message, 1)
			msg[0] = message{index: toggleIndex, state: toggleState}
			cs.publish(msg)
		case err := <-errc:
			log.Printf("failed to serve: %v", err)
			break
		case sig := <-sigs:
			log.Printf("terminating: %v", sig)
			break
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
