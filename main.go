package main

import (
	"fmt"

	"github.com/hs-mb/eplutil"
)

func main() {
	b := eplutil.NewEPLBuilder()
	b.Label()
	err := b.FittedText("PLA\nMüll", 0, 0, b.Width, b.Height, eplutil.FittedTextOptions{
		LineSpace: 0,
		CenterX: true,
		CenterY: true,
	})
	if err != nil {
		panic(err)
	}
	b.Print(1)
	fmt.Print(b.String())
}
