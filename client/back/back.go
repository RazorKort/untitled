package back

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

type Client struct {
	uuid          string
	addr          string
	conn          net.Conn
	reconnectChan chan struct{}
	reconnecting  bool
	sendChan      chan *Packet
	closeChan     chan struct{}
	connMu        sync.RWMutex
	wg            sync.WaitGroup
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
func NewClient(adrr string) *Client {
	uuid := GenerateStableID()
	client := &Client{
		uuid:          uuid,
		addr:          adrr,
		reconnectChan: make(chan struct{}),
		sendChan:      make(chan *Packet, 100),
		closeChan:     make(chan struct{}),
	}

	return client
}

func (c *Client) Connect() {
	for {
		conn, err := net.Dial("tcp", c.addr)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		} else {
			c.conn = conn
			fmt.Println("Connected to server")
			go c.HandleConnection()
			go c.HandleSender()
			go c.reconnect()
			c.sendChan <- &Packet{Type: MsgTypeID, Data: []byte(c.uuid)}
			break
		}
	}

}

func (c *Client) Send(packet *Packet) {
	fmt.Println("Trying to send ", packet)
	c.sendChan <- packet
}

func (c *Client) getConn() net.Conn {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
	return c.conn
}

func (c *Client) setConn(conn net.Conn) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	c.conn = conn
}

func (c *Client) reconnect() {
	for range c.reconnectChan {
		for {
			c.reconnecting = true

			oldConn := c.getConn()
			if oldConn != nil {
				oldConn.Close()
			}

			newConn, err := net.Dial("tcp", c.addr)
			if err != nil {
				fmt.Println(err)
				time.Sleep(5 * time.Second)
				continue
			}
			c.setConn(newConn)
			fmt.Println("Reconnected")
			go c.HandleConnection()
			c.reconnecting = false
			break
		}
	}
}
