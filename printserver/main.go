package main

import (
	"flag"

	"github.com/hs-mb/etikett"
)

var lprBin string
var printer string

func main() {
	lprBinArg := flag.String("b", "lpr", "lpr binary")

	flag.Parse()

	printer = flag.Arg(0)
	tcpAddr := flag.Arg(1)
	if tcpAddr == "" {
		tcpAddr = "0.0.0.0:6244"
	}
	wsAddr := flag.Arg(2)
	if wsAddr == "" {
		wsAddr = "0.0.0.0:6245"
	}

	lprBin = *lprBinArg

	go TCPServer(tcpAddr)
	go WebSocketServer(wsAddr)

	select {}
}

func makePrint(source string) error {
	return etikett.Print(printer, source, lprBin)
}
