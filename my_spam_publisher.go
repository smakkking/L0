package main

import (
	"github.com/nats-io/nats.go"
)

func spam(js nats.JetStreamContext) {
	js.Publish("ORDERS.usa", []byte("hgfjhgfjfgh"))
	js.Publish("ORDERS.russia", []byte("uyiuyi"))
	js.Publish("ORDERS.russia", []byte("sdkalfjg;lkads"))
	js.Publish("ORDERS.russia", []byte("oitjrypiyjrt"))
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := nc.JetStream()

	defer nc.Close()

	spam(js)
}
