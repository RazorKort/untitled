package back

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	MsgTypePing byte = 1
	MsgTypePong byte = 2
	MsgTypeText byte = 3
)

type Packet struct {
	Type byte
	Data []byte
}

func readPacket(conn net.Conn) (*Packet, error) {
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint32(lenBuf)

	typeBuf := make([]byte, 1)
	if _, err := io.ReadFull(conn, typeBuf); err != nil {
		return nil, err
	}

	data := make([]byte, dataLen)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	return &Packet{
		Type: typeBuf[0],
		Data: data,
	}, nil
}

func WritePacket(conn net.Conn, packet *Packet) error {
	dataLen := uint32(len(packet.Data))

	buf := make([]byte, 4+1+dataLen)
	binary.BigEndian.PutUint32(buf[0:4], dataLen)
	buf[4] = packet.Type
	copy(buf[5:], packet.Data)

	_, err := conn.Write(buf)
	return err
}

func CreateConnection() net.Conn {
	conn, err := net.Dial("tcp", "skam.su:10001")
	if err != nil {
		panic(err)
	}
	return conn
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(15 * time.Second)
	fmt.Println("Подключено к серверу")

	for {
		packet, err := readPacket(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			} else {
				fmt.Println("EOF")
			}
			return
		}
		switch packet.Type {
		case MsgTypePing:
			fmt.Printf("<- PING\n")
			WritePacket(conn, &Packet{Type: MsgTypePong, Data: nil})
			fmt.Printf("-> PONG\n")

		case MsgTypePong:
			fmt.Printf("<- PONG\n")

		case MsgTypeText:
			fmt.Printf("<- Data: %s\n", string(packet.Data))

		default:
			fmt.Printf("<- Unrecognized type %d, Data: %s\n", packet.Type, string(packet.Data))
		}
	}
}
