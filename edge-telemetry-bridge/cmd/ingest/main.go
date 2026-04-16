package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	telemetryv1 "example.com/edge-telemetry-bridge/gen/telemetry/v1"
	grpcserverpkg "example.com/edge-telemetry-bridge/internal/grpcserver"
)
// this is a call back to use when we receive a message
var messageHandler mqtt.MessageHandler=func(client mqtt.Client, message mqtt.Message) {
	fmt.Printf(" Received Message on Topic %s : %s\n",message.Topic(),string(message.Payload()))
}
func main() {
	options:=mqtt.NewClientOptions()
	options.AddBroker("tcp://localhost:1883")
	options.SetClientID("go_ingest_client")
	options.SetDefaultPublishHandler(messageHandler)

	client:=mqtt.NewClient(options)
	if token:=client.Connect(); token.Wait() && token.Error()!=nil {
		log.Fatalf(" Failed to Connect %v ",token.Error())
	}
	defer client.Disconnect(250)

	fmt.Println(" Ingest Client Connected to Broker ")

	topic:="telemetry/+/readings"
	if token:=client.Subscribe(topic,1,nil); token.Wait() && token.Error()!=nil {
		log.Fatalf(" Failed to Subscribe %v ",token.Error())
	}
	sig:=make(chan os.Signal,1)
	signal.Notify(sig,os.Interrupt,syscall.SIGTERM)
	lis,err:=net.Listen("tcp",":50052")
	if err!=nil{
		log.Fatalf(" Failed to Listen %v ",err)
	}
	grpcserver:=grpc.NewServer()
	telemetryv1.RegisterTelemetryServiceServer(grpcserver,&grpcserverpkg.Server{})
	reflection.Register(grpcserver)

	go func(){
		log.Println(" Starting gRPC Server on port 50052 ")
		if err:=grpcserver.Serve(lis); err!=nil{
			log.Fatalf(" Failed to Serve %v ",err)
		}
	}()
	<-sig
	fmt.Println(" Ingest Client Exiting ")

}