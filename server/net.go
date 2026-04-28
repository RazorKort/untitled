package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	MsgTypePing    byte = 1
	MsgTypePong    byte = 2
	MsgTypeText    byte = 3
	MsgTypeID      byte = 4
	MsgTypeExecute byte = 4
)

type Packet struct {
	Type byte
	Data []byte
}
type Client struct {
	uuid string
}

func readPacket(conn net.Conn) (*Packet, error) {
	lenBuff := make([]byte, 4)
	_, err := io.ReadFull(conn, lenBuff)
	if err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint32(lenBuff)

	typeBuff := make([]byte, 1)
	_, err = io.ReadFull(conn, typeBuff)
	if err != nil {
		return nil, err
	}

	data := make([]byte, dataLen)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	return &Packet{
		Type: typeBuff[0],
		Data: data,
	}, nil
}

func writePacket(conn net.Conn, packet *Packet) error {
	dataLen := uint32(len(packet.Data))

	buff := make([]byte, 4+1+dataLen)
	binary.BigEndian.PutUint32(buff[0:4], dataLen)
	buff[4] = packet.Type
	copy(buff[5:], packet.Data)

	_, err := conn.Write(buff)
	return err
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Клиент подключился %s\n", conn.RemoteAddr())

	for {
		packet, err := readPacket(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Ошибка чтения %s\n", err)
			} else {
				fmt.Printf("EOF\n")
			}
			return

		}

		switch packet.Type {
		case MsgTypePing:
			fmt.Printf("PING\n")
			writePacket(conn, &Packet{Type: MsgTypePong, Data: nil})

		case MsgTypePong:
			fmt.Printf("PONG\n")

		case MsgTypeText:
			fmt.Printf("Получен текст: %s\n", string(packet.Data))
			writePacket(conn, &Packet{Type: MsgTypeText, Data: packet.Data})

		case MsgTypeID:
			fmt.Printf("Set id of %s to %s\n", conn.RemoteAddr(), packet.Data)
			//ну когда я придумаю как это сторить, то он реально будет менять
		default:
			fmt.Printf("Message type %d\n", packet.Type)
		}
	}
}
