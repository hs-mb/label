//go:build js && wasm

package img

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/gif"
	"image/png"
	"log"
	"syscall/js"

	"github.com/hs-mb/etikett/webprint/script/label"
	"honnef.co/go/js/dom/v2"
)

var (
	previewImage *dom.HTMLImageElement
	printButton *dom.HTMLButtonElement
	fileUpload *dom.HTMLInputElement
)

var img image.Image

func Index() {
	d := dom.GetWindow().Document()

	previewImage = d.GetElementByID("img-preview").(*dom.HTMLImageElement)
	printButton = d.GetElementByID("img-print").(*dom.HTMLButtonElement)
	fileUpload = d.GetElementByID("img-file").(*dom.HTMLInputElement)

	printButton.AddEventListener("click", false, func(e dom.Event) {
		go label.SendPrintServer("")
	})

	fileUpload.AddEventListener("change", false, func(e dom.Event) {
		go changeFileUpload()
	})
}

func changeFileUpload() {
	files := fileUpload.Files()
	if len(files) == 0 {
		return
	}
	finished := make(chan bool)
	buf := new(bytes.Buffer)
	files[0].Call("bytes").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		blockSize := 4096
		offset := 0
		for {
			block := make([]byte, blockSize)
			n := js.CopyBytesToGo(block, args[0].Call("subarray", offset))
			buf.Write(block[:n])
			offset += n
			if n == 0 {
				break
			}
		}
		finished <- true
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
		finished <- false
		log.Printf("Catch\n")
		return nil
	}))
	if !<-finished {
		return
	}
	var err error
	img, _, err = image.Decode(buf)
	if err != nil {
		log.Println(err)
		return
	}
	showImage(img)
}

func showImage(img image.Image) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		log.Println(err)
		return
	}
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", b64)
	previewImage.SetSrc(dataUrl)
}
