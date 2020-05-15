package webserver

import (
	"encoding/json"
	"net/http"
	"os"

	log "github.com/cihub/seelog"
)

// statsHandler returns a JSON array with the data from the various system probes
func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fullStats := getFullStats()

	b, err := json.Marshal(fullStats)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	w.Write(b)
}

// RunServer run the main API to expose server usage
func RunServer() {
	defer log.Flush()

	log.Info("Starting server")
	http.HandleFunc("/api/stats", statsHandler)
	http.ListenAndServe(":5000", nil)
}
