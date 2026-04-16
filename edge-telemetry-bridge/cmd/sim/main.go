package main
import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"math/rand"
	mqtt "github.com/eclipse/paho.mqtt.golang" // shorthand assignment
)
type Reading struct {
	GatewayID string `json:"gateway_id"`
	SensorID string `json:"sensor_id"`
	Value float64 `json:"value"`
}
func main() {
	options :=mqtt.NewClientOptions()
	options.AddBroker("tcp://localhost:1883")
	options.SetClientID("go_simulator_client")

	client:=mqtt.NewClient(options)
	if token:=client.Connect(); token.Wait() && token.Error()!=nil {
		log.Fatalf(" Failed to Connect %v ",token.Error())
	}
	defer client.Disconnect(250)

	fmt.Println("Simulator Connected to Client to Publish messages....")

	for {
		reading:=Reading {
			GatewayID: "gw-001",
			SensorID: "sensor-001",
			Value: 20.0+ float64(rand.Intn(100)),
		}
		// similar to json encoding
		// put the go struct into json format
		payload, err :=json.Marshal(reading)
		if err !=nil {
			log.Fatalf(" Failed to Marshal Reading %v ",err)
			continue
		}
		topic:="telemetry/gw-001/readings"
		token:=client.Publish(topic,1,false,payload)
		token.Wait()
		log.Printf(" Published Message to Topic %s ",topic)
		time.Sleep(1*time.Second)
	}
}