package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/hs-mb/label/webprint/views"
	"github.com/hs-mb/label/webprint/views/label/hackspace"
	"github.com/hs-mb/label/webprint/views/label/raw"
)

var StaticDir = "./static"
var ServeAddr = ":8080"

func main() {
	flag.Parse()

	printServerAddr := flag.Arg(0)
	ctx := context.Background()

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

	mux.HandleFunc("POST /label/print", func(w http.ResponseWriter, r *http.Request) {
		conn, err := net.Dial("tcp", printServerAddr)
		if err != nil {
			log.Print(err)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			return
		}
		_, err = conn.Write(body)
		if err != nil {
			log.Print(err)
			return
		}
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	})

	log.Printf("Listening on %s", ServeAddr)
	http.ListenAndServe(ServeAddr, mux)
}
