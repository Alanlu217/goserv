package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type DashRequest struct {
	Id string
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	rpi := flag.Bool("rpi", false, "Run on Rpi or not")
	flag.Parse()

	var io_handler Io

	if *rpi {
		log.Info().Msg("Using RPiIo")
		io_handler = RPiIo{}
	} else {
		log.Info().Msg("Using LogIo")
		io_handler = LogIo{}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/index.html")
	})
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/favicon.ico")
	})
	mux.HandleFunc("GET /dash", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/dash.html")
	})
	mux.HandleFunc("POST /dash", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		d := DashRequest{}
		decoder.Decode(&d)
		io_handler.Handle(d)
	})

	server := http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: Logging(mux),
	}

	log.Info().Msg(fmt.Sprintf("Server starting on %s", server.Addr))
	defer log.Printf("Server stopping\n")
	log.Fatal().Err((server.ListenAndServe()))
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info().Msg(fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, time.Since(start)))
	})
}
