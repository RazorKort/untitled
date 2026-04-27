package main

import (
	"fmt"
	"time"
	"untitled/back"
)

func main() {
	id := back.GenerateStableID()
	fmt.Println(id)
	conn := back.CreateConnection()
	go back.HandleConnection(conn)
	back.WritePacket(conn, &back.Packet{Type: back.MsgTypePing, Data: nil})
	time.Sleep(200 * time.Second)
}
