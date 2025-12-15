//go:build js && wasm

package label

import (
	"fmt"
	"net/http"
	"strings"

	"honnef.co/go/js/dom/v2"
)

func SendPrintServer(source string) error {
	location := dom.GetWindow().Location()
	url := fmt.Sprintf("%s//%s/label/print", location.Protocol(), location.Host())
	_, err := http.Post(url, "application/x.epl2", strings.NewReader(source))
	return err
}
