package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/hs-mb/etikett"
)

var lprBin string

func main() {
	lprBinArg := flag.String("b", "lpr", "lpr binary")

	flag.Parse()

	printer := flag.Arg(0)
	addr := flag.Arg(1)

	lprBin = *lprBinArg

	log.Printf("Listening on %s", addr)

	l, err := net.Listen("tcp", addr)
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
	err := label.Print(printer, string(packet), lprBin)
	if err != nil {
		log.Printf("Print err: %v", err)
	}
}

