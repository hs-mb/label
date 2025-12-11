package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
)

func main() {
	flag.Parse()

	printer := flag.Arg(0)
	addr := flag.Arg(1)

	log.Printf("Listening on %s", addr)

	l, err := net.Listen("tcp4", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Err: %v", err)
			continue
		}
		go handle(printer, c)
	}
}

func handle(printer string, c net.Conn) {
	log.Printf("Conn: %v", c.RemoteAddr().String())
	defer c.Close()
	packet := make([]byte, 0, 4096)
	for {
		recv := make([]byte, 4096)
		n, err := c.Read(recv)
		if err != nil {
			if err != io.EOF {
				log.Printf("Conn err: %v", err)
			}
			break
		}
		packet = append(packet, recv[:n]...)
	}
	err := print(printer, string(packet))
	if err != nil {
		log.Printf("Print err: %v", err)
	}
}

func print(printer string, data string) error {
	cmd := exec.Command("lpr", "-P", printer, "-o", "raw")
	cmd.Stdin = strings.NewReader(data)
	return cmd.Run()
}

