package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const upnp_request = "M-SEARCH * HTTP/1.1\r\n" +
	"HOST: 239.255.255.250:1900\r\n" +
	`MAN: "ssdp:discover"` + "\r\n" +
	"MX: 10\r\n" +
	"ST: ssdp:all\r\n" +
	"\r\n"

func main() {
	conn, err := net.ListenPacket("udp4", ":1900")
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()
	log.Printf("Listening for upnp packets on %v\n", conn.LocalAddr())

	go func() {
		buf := make([]byte, 8192)
		for {
			n, _, _ := conn.ReadFrom(buf)
			if n > 0 {
				fmt.Print(string(buf[0:n]))
			}
		}
	}()

	upnp_broadcast := &net.UDPAddr{
		IP:   net.IPv4(239, 255, 255, 250),
		Port: 1900}
	for {
		// send upnp request every minute
		conn.WriteTo([]byte(upnp_request), upnp_broadcast)
		log.Println("Sent M-SEARCH request")
		time.Sleep(1 * time.Minute)
	}
}
