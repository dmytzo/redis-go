package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	CommandPING = "PING"

	SignString = "+"
)

func strb(s string) []byte {
	return []byte(fmt.Sprintf("%s%s\r\n", SignString, s))
}

func main() {
	fmt.Println("redis")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatalf("bind to port: %s", err.Error())
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("accept connection: %s", err.Error())
		}

		c.SetReadDeadline(time.Now().Add(10 * time.Second))

		content := make([]byte, 1024)
		if _, err := c.Read(content); err != nil {
			log.Fatalf("read: %s", err.Error())

		}

		if _, err := c.Write(strb("PONG")); err != nil {
			log.Fatalf("write: %s", err.Error())
		}

		//if err := c.Close(); err != nil {
		//	log.Fatalf("close: %s", err.Error())
		//}
		//switch string(content) {
		//case CommandPING:
		//	fmt.Println("herre")
		//	if _, err := c.Write(strb("PONG")); err != nil {
		//		log.Fatalf("write: %s", err.Error())
		//	}
		//}
		//
		//if err := c.Close(); err != nil {
		//	log.Fatalf("close: %s", err.Error())
		//}
	}
}
