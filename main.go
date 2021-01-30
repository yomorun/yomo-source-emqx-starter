package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	y3 "github.com/yomorun/y3-codec-golang"
	quicy "github.com/yomorun/yomo/pkg/quic"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	// Create connection to YoMo-Zipper
	var addr = "localhost:9999"
	qc, err := quicy.NewClient(addr)
	if err != nil {
		log.Printf("NewClient addr=%s error:%s", addr, err.Error())
		return
	}
	stream, err := qc.CreateStream(context.Background())
	if err != nil {
		log.Printf("CreateStream addr=%s error:%s", addr, err.Error())
		return
	}

	// Initial MQTT
	var broker = "broker.emqx.io"
	// var broker = "172.16.0.191"
	var port = 1883

	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		rb := msg.Payload()
		fmt.Printf("Received message: %v from topic: %s \n", rb, msg.Topic())

		p, _ := strconv.ParseInt(string(rb), 10, 64)
		var packet = y3.NewPrimitivePacketEncoder(0x01)
		packet.SetInt64Value(p)

		// send data via QUIC stream.
		buf := packet.Encode()
		_, err = stream.Write(buf)
		if err != nil {
			log.Printf("❌ Emit %s to yomo-zipper failure with err: %v", msg.Payload(), err)
		} else {
			log.Printf("✅ Sending message: %v to yomo-zipper", buf)
		}
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)
	publish(client)

	client.Disconnect(250)
}

func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("%d", i*100)
		token := client.Publish("topic/yomo", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func sub(client mqtt.Client) {
	topic := "topic/yomo"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
