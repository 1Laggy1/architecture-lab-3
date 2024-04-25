package main

import (
	"net/http"
	"path/filepath"

	"github.com/1Laggy1/architecture-lab-3/painter"
	"github.com/1Laggy1/architecture-lab-3/painter/lang"
	"github.com/1Laggy1/architecture-lab-3/ui"
)

func main() {
	var (
		pv ui.Visualizer // Візуалізатор створює вікно та малює у ньому.

		// Потрібні для частини 2.
		opLoop painter.Loop // Цикл обробки команд.
		parser lang.Parser  // Парсер команд.
	)

	//pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	// Handle requests for the root ("/") endpoint with lang.HttpHandler
	http.Handle("/", lang.HttpHandler(&opLoop, &parser))

	// Serve the index.html file
	http.HandleFunc("/buttons", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	})

	// Serve the main.js file
	http.HandleFunc("/scripts/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join("../../scripts", "main.js"))
	})

	http.HandleFunc("/scripts/green.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join("../../scripts", "green.js"))
	})

	// Serve the white.js file
	http.HandleFunc("/scripts/white.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join("../../scripts", "white.js"))
	})
	http.HandleFunc("/scripts/update.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join("../../scripts", "update.js"))
	})
	http.HandleFunc("/scripts/move.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join("../../scripts", "move.js"))
	})
	http.HandleFunc("/scripts/figure.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join("../../scripts", "figure.js"))
	})

	// Start the HTTP server
	go func() {
		if err := http.ListenAndServe("localhost:17000", nil); err != nil {
			panic(err)
		}
	}()

	// Start the ui.Visualizer
	pv.Main()

	// Stop and wait for the painter.Loop to finish
	opLoop.StopAndWait()
}
