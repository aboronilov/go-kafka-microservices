package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/aboronilov/go-kafka-microservices/types"
	"github.com/gorilla/websocket"
)

const wsEndpoint = "ws://127.0.0.1:30000/ws"

var sendInterval = time.Second

func genCoord() float64 {
	return (float64(rand.Intn(10000)+1)/10000.0)*180.0 - 90.0
}

func genLocation() (float64, float64) {
	return genCoord(), genCoord()
}

func main() {
	obuIDs := generateOBUIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := genLocation()
			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal("write data error:", err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDs(n int) []int {
	obuIDs := []int{}
	for i := 0; i < n; i++ {
		id := rand.Intn(1000000) + 1
		obuIDs = append(obuIDs, id)
	}

	return obuIDs
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}
