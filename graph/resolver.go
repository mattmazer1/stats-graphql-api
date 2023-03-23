package graph

import "github.com/nats-io/nats.go"

type Resolver struct {
	Nats *nats.Conn
}
