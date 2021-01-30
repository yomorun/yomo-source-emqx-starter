![yomo-integrate-with-emqx]()

# yomo-source-emqx-starter

EMQ X Broker 🙌 YoMo

## About EMQX

EMQ X broker is a fully open source, highly scalable, highly available distributed MQTT messaging broker for IoT, M2M and Mobile applications that can handle tens of millions of concurrent clients.

Starting from 3.0 release, EMQ X broker fully supports MQTT V5.0 protocol specifications and backward compatible with MQTT V3.1 and V3.1.1, as well as other communication protocols such as MQTT-SN, CoAP, LwM2M, WebSocket and STOMP. The 3.0 release of the EMQ X broker can scaled to 10+ million concurrent MQTT connections on one cluster.

For more information, please visit [EMQ X homepage](https://www.emqx.io/)

## 1: Installing via EMQX Docker Image or use EMQX public cloud

```bash
docker pull emqx/emqx
```

start a single node

```bash
sudo docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8883:8883 -p 8084:8084 -p 18083:18083 emqx/emqx
```

[EMQX officai installation page](https://docs.emqx.io/en/broker/latest/getting-started/install.html)

## 2: Connect EMQX to YoMo

```go
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
```

## 3: Create YoMo-Zipper

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

Run `cd zipper && yomo wf run`

## 4: Write your data process logic

see `./flow/app.go`

```go
var printer = func(_ context.Context, i interface{}) (interface{}, error) {
	value := i.(int64)
	fmt.Println("serverless get value:", value)
	return value, nil
}

var callback = func(v []byte) (interface{}, error) {
	res, _, _, err := y3.DecodePrimitivePacket(v)
	if err != nil {
		log.Errorf("y3 err: %v", err)
	}

	val, err := res.ToInt64()
	if err != nil {
		log.Errorf("y3 toint64 err: %v", err)
	}
	return val, nil
}

// Handler will handle data in Rx way
func Handler(rxstream rx.RxStream) rx.RxStream {
	stream := rxstream.
		Subscribe(0x01).
		OnObserve(callback).
		Map(printer).
		StdOut()

	return stream
}
```

run `cd ./flow && yomo run`