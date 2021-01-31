![yomo-integrate-with-emqx](emqx-yomo.jpg)

# yomo-source-emqx-starter

EMQ X Broker 🙌 YoMo

This project demonstrates how to connect to [EMQX MQTT Broker](https://www.emqx.io) and processing data before [store to a Serverless Database FaunaDB](https://github.com/yomorun/yomo-sink-faunadb-example), the code use MQTT public cloud so you could run directly without installation. Code shows a mathematical SUM operation when every data arrived [in Rx way](https://reactivex.io).

## About EMQX

EMQ X broker is a fully open source, highly scalable, highly available distributed MQTT messaging broker for IoT, M2M and Mobile applications that can handle tens of millions of concurrent clients.

Starting from 3.0 release, EMQ X broker fully supports MQTT V5.0 protocol specifications and backward compatible with MQTT V3.1 and V3.1.1, as well as other communication protocols such as MQTT-SN, CoAP, LwM2M, WebSocket and STOMP. The 3.0 release of the EMQ X broker can scaled to 10+ million concurrent MQTT connections on one cluster.

For more information, please visit [EMQ X homepage](https://www.emqx.io/)

## 1: Connect EMQX to YoMo

```go
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	rb := msg.Payload()
	fmt.Printf("Received message: %v from topic: %s \n", rb, msg.Topic())

	// encode data by Y3 codec
	p, _ := strconv.ParseInt(string(rb), 10, 64)
	codec := y3.NewCodec(0x10)
	buf, _ := codec.Marshal(p)

	// send data via QUIC stream
	_, err = stream.Write(buf)
	if err != nil {
		log.Printf("❌ Emit %s to yomo-zipper failure with err: %v", msg.Payload(), err)
	} else {
		log.Printf("✅ Sending message: %v to yomo-zipper", buf)
	}
}
```

Run `go run main.go` to start

## 2: Create YoMo-Zipper

see `./zipper/workflow.yaml`

```yaml
name: Service
host: localhost
port: 9999
flows:
  - name: Noise_Serverless
    host: localhost
    port: 4242

```

Run `cd zipper && yomo wf run` to start the [yomo-zipper](https://yomo.run/zipper)

## 3: Write your data process logic

see `./flow/app.go`

```go
const observedKey = 0x10

// Decode data by Y3 as Int64 type from YoMo-Zipper
var decode = func(v []byte) (interface{}, error) {
	val, err := y3.ToInt64(v)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return val, nil
}

// Calculate sum every new data arrived, about Scan() Operator, can be read here: http://reactivex.io/documentation/operators/scan.html
var sum = func(_ context.Context, acc interface{}, elem interface{}) (interface{}, error) {
	if acc == nil {
		return elem, nil
	}
	return acc.(int64) + elem.(int64), nil
}

// Handler handle data in Rx way
func Handler(rxstream rx.RxStream) rx.RxStream {
	stream := rxstream.
		Subscribe(observedKey).
		OnObserve(decode).
		Scan(sum).
		StdOut()

	return stream
}
```

run `cd ./flow && yomo run` to start 

## More about YoMo

[YoMo Repository](https://github.com/yomorun/yomo)

