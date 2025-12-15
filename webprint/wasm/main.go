//go:build js && wasm

package main

import (
	"strings"

	"github.com/hs-mb/label/webprint/script"
	"github.com/hs-mb/label/webprint/script/label/hackspace"
	"github.com/hs-mb/label/webprint/script/label/raw"
	"honnef.co/go/js/dom/v2"
)

func main() {
	path := dom.GetWindow().Location().Pathname()
	path = strings.TrimSuffix(path, "/")

	switch path {
	case "":
		script.Index()
	case "/label/hackspace":
		hackspace.Index()
	case "/label/raw":
		raw.Index()
	}

	select {}
}
