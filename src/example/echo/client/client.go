package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"binaryTCP/src/example/echo"
	"binaryTCP/src/gotcp"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()

	// creates a server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	client := gotcp.NewClient(config, &Callback{}, &echo.EchoProtocol{})
	// stops service
	defer client.Stop()
	// starts service
	go client.Start(conn)

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	echoPacket := p.(*echo.EchoPacket)
	fmt.Printf("Server response :[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}
