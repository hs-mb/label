//go:build js && wasm

package owner

import (
	"fmt"
	"strings"
	"time"

	"github.com/hs-mb/eplutil"
	"github.com/hs-mb/etikett/webprint"
	"github.com/hs-mb/etikett/webprint/script/label"
	"honnef.co/go/js/dom/v2"
)

var (
	inputField *dom.HTMLDivElement
	printButton *dom.HTMLButtonElement
)

func Index() {
	d := dom.GetWindow().Document()

	inputField = d.GetElementByID("owner-text").(*dom.HTMLDivElement)
	printButton = d.GetElementByID("owner-print").(*dom.HTMLButtonElement)

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
	date := time.Now().Local().Format("02.01.2006")
	b := eplutil.NewEPLBuilder()
	header := fmt.Sprintf("Eigentümer    %s", date)
	b.Label()
	b.FittedText(header, 50, 0, b.Width - 100, 70, eplutil.FittedTextOptions{
		CenterX: false,
		CenterY: true,
		Font: webprint.Font,
	})
	b.FittedText(text, 10, 90, b.Width - 20, b.Height - 110, eplutil.FittedTextOptions{
		CenterX: true,
		CenterY: true,
		Font: webprint.Font,
	})
	b.Print(1)
	return b.String()
}
