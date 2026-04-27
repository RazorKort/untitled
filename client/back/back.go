package back

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
)

type Client struct {
	conn      net.Conn
	sendChan  chan *Packet
	closeChan chan struct{}
	wg        sync.WaitGroup
}

// генерирует стабильный uuid
func GenerateStableID() string {
	hostname, _ := os.Hostname()

	var mac string
	if interfaces, err := net.Interfaces(); err == nil && len(interfaces) > 0 {
		mac = interfaces[0].HardwareAddr.String()
	}
	osInfo := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	data := fmt.Sprintf("%s|%s|%s", hostname, mac, osInfo)
	hash := sha256.Sum256([]byte(data))

	uuid := hex.EncodeToString(hash[:16])
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		uuid[0:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:32])
}

// конструктор лего
func NewClient(conn net.Conn) *Client {
	client := &Client{
		conn:      conn,
		sendChan:  make(chan *Packet, 100),
		closeChan: make(chan struct{}),
	}

	return client
}
