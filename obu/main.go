package main

import (
	"log"
	"math/rand"
	"time"
	"tolling/types"

	"github.com/gorilla/websocket"
)

var (
	sendInterval = time.Second
	wsEndponint = "ws://127.0.0.1:30000/ws"
)

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(100000000)
	}
	return ids
}

func main() {
	obuIDS := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndponint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
