package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

var Nc *nats.Conn

func ConnectNat() {
	var err error

	Nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	if !Nc.IsConnected() {
		log.Fatalf("Failed to connect to NATS server")
	}

	log.Printf("connected to NATS server")
}

func CloseNat() {
	Nc.Close()
}
