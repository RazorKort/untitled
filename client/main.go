package main

import (
	"time"
	"untitled/back"
)

func main() {
	client := back.NewClient("skam.su:10001")
	client.Connect()
	for range 10 {
		client.Send(&back.Packet{Type: back.MsgTypePing, Data: nil})
		time.Sleep(5 * time.Second)
	}

}
