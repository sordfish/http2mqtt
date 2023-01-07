# http2mqtt

Does what it says on the tin, takes json and transforms it to an mqtt message - no auth no tls!

config
-----
- MQTT Broker address -
```
-broker="tcp://127.0.0.1:1883"
```


json format
-----
```
{
	"topic": "yourTopicHere",
	"message": "addYourDataHere"
}
```

Helm Chart
-----
Nothing special, just built using the create chart feature of Helm.

Deployment
-----
Set broker via Values.mqtt.broker

Service
-----
Listening on port 80