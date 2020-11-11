package main

import (
	"github.com/sajja/golearn/statemachine"
)

func hello(i int) string {
	println("dakk")
	return "world"
}

type MyType struct {
	Field1 string
	Field2 string
}

func main() {
	// pointers.TestPointer()
	statemachine.TestMultipleStateMachines()
	// channels.TestChannels()
	// mqtt.Mqtt_pub()
}
