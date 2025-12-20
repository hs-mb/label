package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/hs-mb/etikett/webprint"
	"github.com/hs-mb/etikett/webprint/views"
	"github.com/hs-mb/etikett/webprint/views/label/hackspace"
	"github.com/hs-mb/etikett/webprint/views/label/img"
	"github.com/hs-mb/etikett/webprint/views/label/raw"
)

var StaticDir = "./static"
var ServeAddr = ":8080"

func main() {
	flag.Parse()

	printServerAddr := flag.Arg(0)
	ctx := context.Background()
	ctx = context.WithValue(ctx, webprint.PrintAddrKey, printServerAddr)

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir(StaticDir))))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(ctx, w)
	})
	mux.HandleFunc("GET /label/raw", func(w http.ResponseWriter, r *http.Request) {
		raw.Index().Render(ctx, w)
	})
	mux.HandleFunc("GET /label/hackspace", func(w http.ResponseWriter, r *http.Request) {
		hackspace.Index().Render(ctx, w)
	})
	mux.HandleFunc("GET /label/img", func(w http.ResponseWriter, r *http.Request) {
		img.Index().Render(ctx, w)
	})

	log.Printf("Listening on %s", ServeAddr)
	http.ListenAndServe(ServeAddr, mux)
}
