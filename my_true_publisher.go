package main

import (
	"io/ioutil"

	"github.com/nats-io/nats.go"
)

func norm_values(js nats.JetStreamContext) {
	content, _ := ioutil.ReadFile("./json_data/model.json")

	// нормальное значение
	js.Publish("ORDERS.usa", content)

	content, _ = ioutil.ReadFile("./json_data/model2.json")

	// нормальное значение
	js.Publish("ORDERS.russia", content)
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := nc.JetStream()

	defer nc.Close()

	norm_values(js)
}
