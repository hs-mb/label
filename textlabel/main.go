package main

import (
	"flag"

	"github.com/hs-mb/eplutil"
	"github.com/hs-mb/etikett"
	"golang.org/x/image/font/opentype"

	_ "embed"
)

//go:embed res/univga.ttf
var fontData []byte
var font *opentype.Font

func init() {
	var err error
	font, err = opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}
}

func main() {
	lineSpace := flag.Int("l", 10, "Line space")
	margin := flag.Int("m", 20, "Margin")
	n := flag.Int("n", 1, "Label number")

	flag.Parse()

	text := flag.Arg(0)
	printer := flag.Arg(1)

	b := eplutil.NewEPLBuilder()
	b.Label()
	b.FittedText(text, *margin, *margin, b.Width - (2 * *margin), b.Height - (2 * *margin), eplutil.FittedTextOptions{
		LineSpace: *lineSpace,
		CenterX: true,
		CenterY: true,
		Font: font,
	})
	b.Print(*n)


	err := label.Print(printer, b.String())
	if err != nil {
		panic(err)
	}
}
