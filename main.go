package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTPacket struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func livez(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

func healthz(w http.ResponseWriter, req *http.Request, c mqtt.Client) {
	if c.IsConnected() {
		fmt.Fprintf(w, "OK\n")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Not connected to mqtt"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error with JSON marshal - Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

}

func getData(rw http.ResponseWriter, req *http.Request, c mqtt.Client) {
	decoder := json.NewDecoder(req.Body)
	var p MQTTPacket
	err := decoder.Decode(&p)
	if err != nil {
		log.Fatalf("Error with JSON decoder - Err: %s", err)
	}

	log.Println(p.Topic)
	log.Println(p.Message)

	c.Publish(p.Topic, 0, false, p.Message)
}

func main() {

	brokerPtr := flag.String("broker", "tcp://127.0.0.1:1883", "MQTT broker address")
	flag.Parse()

	// Connect to the MQTT server
	opts := mqtt.NewClientOptions().AddBroker(*brokerPtr)
	opts.SetClientID("http2mqtt")
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error with MQTT - Err: %s", token.Error())
	}

	http.HandleFunc("/livez", livez)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthz(w, r, c)
	})
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		getData(w, r, c)
	})
	http.ListenAndServe(":45678", nil)

	c.Disconnect(250)

}
