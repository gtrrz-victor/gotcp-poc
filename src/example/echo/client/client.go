package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

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
	srv := gotcp.NewServer(config, &Callback{}, &echo.EchoProtocol{})
	// stops service
	defer srv.Stop()
	// starts service
	go srv.StartDual(conn, time.Second)

	http.HandleFunc("/hello", hello(conn))
	http.ListenAndServe(":9000", nil)
	/*
		// catchs system signal
		chSig := make(chan os.Signal)
		signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
		fmt.Println("Signal: ", <-chSig)
	*/
}

func hello(conn *net.TCPConn) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("On HTTP request. Request body content:", (string(bodyBytes)))
		fmt.Println("Sending content through tcp...")
		conn.Write(echo.NewEchoPacket(bodyBytes, false).Serialize())
	}
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
