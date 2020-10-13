package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aHugues/system-monitor/monitor/utils"

	log "github.com/sirupsen/logrus"
)

// statsHandler returns a JSON array with the data from the various system probes
func statsHandler(config utils.ProbesConfig, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fullStats := getFullStats(config)

	b, err := json.Marshal(fullStats)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	w.Write(b)
}

// displayConfig returns a JSON object representing the current monitoring configuration
func displayConfig(config utils.ProbesConfig, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(config)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(b)
	}
}

// RunServer run the main API to expose server usage
func RunServer(config utils.FullConfiguration) {

	log.Info("Starting server")
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		statsHandler(config.Probes, w, r)
	})
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		displayConfig(config.Probes, w, r)
	})

	listenFullHost := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Debugf("Server listening on %q", listenFullHost)

	http.ListenAndServe(listenFullHost, nil)
}
