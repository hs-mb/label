//go:build js && wasm

package hackspace

import (
	"strings"

	"github.com/hs-mb/eplutil"
	"github.com/hs-mb/label/webprint/script/label"
	"honnef.co/go/js/dom/v2"
)

var (
	inputField *dom.HTMLDivElement
	printButton *dom.HTMLButtonElement
)

func Index() {
	d := dom.GetWindow().Document()

	inputField = d.GetElementByID("hs-text").(*dom.HTMLDivElement)
	printButton = d.GetElementByID("hs-print").(*dom.HTMLButtonElement)

	inputField.AddEventListener("input", false, func(e dom.Event) {
		if strings.Trim(inputField.TextContent(), "\n ") == "" {
			inputField.SetInnerHTML("")
		}
		inputChange()
	})
	inputChange()

	printButton.AddEventListener("click", false, func(e dom.Event) {
		go label.SendPrintServer(makeLabel(inputField.TextContent()))
	})
}

func inputChange() {
	printButton.SetDisabled(inputField.TextContent() == "")
}

func makeLabel(text string) string {
	b := eplutil.NewEPLBuilder()
	b.FittedText(text, 0, 0, b.Width, b.Height, eplutil.FittedTextOptions{
		CenterX: true,
		CenterY: true,
	})
	b.Print(1)
	return b.String()
}
