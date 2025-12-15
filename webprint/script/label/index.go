//go:build js && wasm

package label

import (
	"fmt"
	"net/http"
	"strings"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

var (
	inputField *dom.HTMLTextAreaElement
	printButton *dom.HTMLButtonElement
	errorText dom.Element
	fileUpload *dom.HTMLInputElement
)

func Raw() {
	d := dom.GetWindow().Document()

	inputField = d.GetElementByID("raw-input").(*dom.HTMLTextAreaElement)
	printButton = d.GetElementByID("raw-print").(*dom.HTMLButtonElement)
	errorText = d.GetElementByID("raw-error")
	fileUpload = d.GetElementByID("raw-file").(*dom.HTMLInputElement)

	inputFieldChange := func(e dom.Event) {
		printButton.SetDisabled(inputField.TextLength() == 0)
	}
	inputField.AddEventListener("change", false, inputFieldChange)
	inputFieldChange(nil)

	printButton.AddEventListener("click", false, func(e dom.Event) {
		text := inputField.Value()
		go sendPrintServer(text)
	})

	fileUpload.AddEventListener("change", false, func(e dom.Event) {
		go changeFileUpload()
	})

	select {}
}

func changeFileUpload() {
	files := fileUpload.Files()
	if len(files) == 0 {
		return
	}
	text := make(chan string)
	files[0].Call("text").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		text <- args[0].String()
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
		text <- ""
		return nil
	}))
	inputField.SetTextContent(<-text)
}

func sendPrintServer(source string) {
	location := dom.GetWindow().Location()
	url := fmt.Sprintf("%s//%s/label/print", location.Protocol(), location.Host())
	_, err := http.Post(url, "application/x.epl2", strings.NewReader(source))
	if err != nil {
		errorText.SetTextContent(err.Error())
	}
}
