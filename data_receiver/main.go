package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/gorilla/websocket"
)

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
	prod  DataProducer
}

func NewDataRecevier() (*DataRecevier, error) {
	p, err := NewKafkaProducer("obudata")
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleware(p)

	return &DataRecevier{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
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
		data.RequestID = rand.Intn(math.MaxInt)
		if err = dr.prod.ProduceData(data); err != nil {
			log.Printf("Error producing data: %v\n", err)
			continue
		}
		// dr.prod.Flush(15 * 1000)
	}
}
