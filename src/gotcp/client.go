package gotcp

import (
	"fmt"
	"net"
	"sync"
)

type Client struct {
	Server
}

// NewServer creates a server
func NewClient(config *Config, callback ConnCallback, protocol Protocol) *Client {
	return &Client{Server: Server{
		config:    config,
		callback:  callback,
		protocol:  protocol,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}}
}

// Start starts service
func (c *Client) Start(conn *net.TCPConn) {
	defer conn.Close()
	go newConn(conn, &c.Server).Do()
	exitSignal := <-c.exitChan
	fmt.Println("Exit signal:", exitSignal)
}

// Stop stops service
func (c *Client) Stop() {
	close(c.exitChan)
}
