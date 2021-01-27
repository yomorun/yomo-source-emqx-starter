![yomo-integrate-with-emqx]()

# yomo-source-emqx-starter

EMQ X Broker 🙌 YoMo

## About EMQX

EMQ X broker is a fully open source, highly scalable, highly available distributed MQTT messaging broker for IoT, M2M and Mobile applications that can handle tens of millions of concurrent clients.

Starting from 3.0 release, EMQ X broker fully supports MQTT V5.0 protocol specifications and backward compatible with MQTT V3.1 and V3.1.1, as well as other communication protocols such as MQTT-SN, CoAP, LwM2M, WebSocket and STOMP. The 3.0 release of the EMQ X broker can scaled to 10+ million concurrent MQTT connections on one cluster.

For more information, please visit [EMQ X homepage](https://www.emqx.io/)

## 1/x Installing via EMQX Docker Image

```bash
docker pull emqx/emqx
```

start a single node

```bash
sudo docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8883:8883 -p 8084:8084 -p 18083:18083 emqx/emqx
```

[EMQX officai installation page](https://docs.emqx.io/en/broker/latest/getting-started/install.html)

## 2/
