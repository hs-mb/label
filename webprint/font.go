package webprint

import (
	"golang.org/x/image/font/opentype"
	_ "embed"
)

//go:embed res/univga.ttf
var fontData []byte
var Font *opentype.Font

func init() {
	var err error
	Font, err = opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}
}

