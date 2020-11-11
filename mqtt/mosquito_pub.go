package mqtt

import (
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func connLostHandler(c mqtt.Client, err error) {
	fmt.Printf("Connection lost, reason: %v\n", err)

	//Perform additional action...
}

func Mqtt_pub() {
	opts := MQTT.NewClientOptions().
		AddBroker("tcp://localhost:1883").
		SetConnectionLostHandler(connLostHandler)
	opts.SetClientID("go-simple")
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	j := 0
	for {
		text := fmt.Sprintf("this is msg #%d!", j)
		token := c.Publish("go-mqtt/sample", 2, true, text)
		token.WaitTimeout(1 * time.Second)
		fmt.Printf("Published....%s, Connected %b\n", text, c.IsConnected())
		time.Sleep(2 * time.Second)
		j = j + 1
	}

}

func Mqtt_pubsub() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("go-simple")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("go-mqtt/sample", 0, false, text)
		token.Wait()
	}

	time.Sleep(3 * time.Second)

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
