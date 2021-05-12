package main

import (
	"log"
	"time"

	"github.com/bh90210/healthz"
)

func main() {
	// setting those values is optional, could be used with "" and the bellow default values.
	h := healthz.NewCheck("live", "ready", "8080")

	go func() {
		if err := h.Start(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		log.Println("ready")
		h.Ready()
		time.Sleep(time.Second * 60)
		log.Println("not ready")
		h.NotReady()
	}()

	if term := h.Terminating(); term == true {
		// do something
	}
}
