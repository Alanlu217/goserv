package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	led_pin = 18

	on_1_pin  = 25
	off_1_pin = 24

	on_2_pin  = 23
	off_2_pin = 22

	on_3_pin  = 18
	off_3_pin = 17

	on_4_pin  = 8
	off_4_pin = 7
)

type DashRequest struct {
	Id int
}

type Button struct {
	Name    string
	Handler func()
}

func main() {
	setup_cleanup()

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

	buttons := []Button{
		{
			Name:    "1 On",
			Handler: func() { io_handler.TogglePin(on_1_pin) },
		},
		{
			Name:    "2 On",
			Handler: func() { io_handler.TogglePin(on_2_pin) },
		},
		{
			Name:    "3 On",
			Handler: func() { io_handler.TogglePin(on_3_pin) },
		},
		{
			Name:    "4 On",
			Handler: func() { io_handler.TogglePin(on_4_pin) },
		},
		{
			Name:    "1 Off",
			Handler: func() { io_handler.TogglePin(off_1_pin) },
		},
		{
			Name:    "2 Off",
			Handler: func() { io_handler.TogglePin(off_2_pin) },
		},
		{
			Name:    "3 Off",
			Handler: func() { io_handler.TogglePin(off_3_pin) },
		},
		{
			Name:    "4 Off",
			Handler: func() { io_handler.TogglePin(off_4_pin) },
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/index.html")
	})

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/favicon.ico")
	})

	mux.HandleFunc("GET /uptime", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("uptime")
		data, err := cmd.Output()
		if err != nil {
			log.Err(err)
			return
		}

		w.Write(data)
	})

	dash_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("site/dash.html")
		if err != nil {
			log.Err(err)
			fmt.Fprint(w, "Internal Error")
			return
		}

		data := map[string]any{
			"Buttons": buttons,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Err(err)
		}
	})

	mux.Handle("GET /dash", Auth(dash_handler))

	mux.HandleFunc("GET /internal/dash.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/dash.css")
	})

	mux.HandleFunc("GET /internal/dash.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/dash.js")
	})

	mux.HandleFunc("POST /dash", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		d := DashRequest{}
		decoder.Decode(&d)
		go buttons[d.Id].Handler()
	})

	server := http.Server{
		Addr:    "0.0.0.0:8888",
		Handler: Logging(mux),
	}

	log.Info().Msgf("Server starting on %s", server.Addr)
	defer log.Printf("Server stopping\n")
	log.Fatal().Err((server.ListenAndServe()))
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info().Msgf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte("aln"))
			expectedPasswordHash := sha256.Sum256([]byte("000"))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func setup_cleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Warn().Msg("Closing server")
		os.Exit(1)
	}()
}
