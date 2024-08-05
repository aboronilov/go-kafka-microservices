package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("Data recevier started")
	recv := NewDataRecevier()
	http.HandleFunc("/ws", recv.wsHandler)
	http.ListenAndServe(":30000", nil)
}

type DataRecevier struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

func NewDataRecevier() *DataRecevier {
	return &DataRecevier{
		msgch: make(chan types.OBUData, 128),
	}
}

func (dr *DataRecevier) wsHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("Origin") != "http://"+r.Host {
	// 	http.Error(w, "Origin not allowed", http.StatusForbidden)
	// 	return
	// }

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
		fmt.Printf("Received OBU data: %+v\n", data)
		dr.msgch <- data
	}
}
