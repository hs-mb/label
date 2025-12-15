//go:build js && wasm

package raw

import (
	"syscall/js"

	"github.com/hs-mb/label/webprint/script/label"
	"honnef.co/go/js/dom/v2"
)

var (
	inputField *dom.HTMLTextAreaElement
	printButton *dom.HTMLButtonElement
	errorText dom.Element
	fileUpload *dom.HTMLInputElement
)

func Index() {
	d := dom.GetWindow().Document()

	inputField = d.GetElementByID("raw-input").(*dom.HTMLTextAreaElement)
	printButton = d.GetElementByID("raw-print").(*dom.HTMLButtonElement)
	errorText = d.GetElementByID("raw-error")
	fileUpload = d.GetElementByID("raw-file").(*dom.HTMLInputElement)

	inputField.AddEventListener("change", false, inputChange)
	inputField.AddEventListener("input", false, inputChange)
	inputChange(nil)

	printButton.AddEventListener("click", false, func(e dom.Event) {
		go label.SendPrintServer(inputField.Value())
	})

	fileUpload.AddEventListener("change", false, func(e dom.Event) {
		go changeFileUpload()
	})
}

func inputChange(_ dom.Event) {
	printButton.SetDisabled(inputField.Value() == "")
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
	inputField.SetValue(<-text)
	inputChange(nil)
}
