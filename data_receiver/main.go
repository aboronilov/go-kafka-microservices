package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
)

var kafkaTopic = "obudata"

func main() {
	fmt.Println("Data recevier started")

	recv, err := NewDataRecevier()
	if err != nil {
		log.Fatal("Error creating data receiver:", err)
	}

	// defer p.Close()

	http.HandleFunc("/ws", recv.wsHandler)
	http.ListenAndServe(":30000", nil)
}

type DataRecevier struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  *kafka.Producer
}

func NewDataRecevier() (*DataRecevier, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &DataRecevier{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataRecevier) produceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaTopic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)

	return err
}

func (dr *DataRecevier) wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error connecting to WS:", err)
	}
	// defer conn.Close()

	dr.conn = conn
	go dr.wsRecieveLoop()
}

func (dr *DataRecevier) wsRecieveLoop() {
	fmt.Println("New OBU client connected")
	for {
		var data types.OBUData
		err := dr.conn.ReadJSON(&data)
		if err != nil {
			log.Printf("Error reading from WS: %v\n", err)
			continue
		}
		if err = dr.produceData(data); err != nil {
			log.Printf("Error producing data: %v\n", err)
			continue
		}
		// dr.prod.Flush(15 * 1000)
	}
}
