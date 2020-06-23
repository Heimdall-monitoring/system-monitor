package webserver

import (
	"encoding/json"
	"fmt"
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
func RunServer(config Configuration) {
	defer log.Flush()

	log.Info("Starting server")
	http.HandleFunc("/api/stats", statsHandler)

	listenFullHost := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Debugf("Server listening on %q", listenFullHost)

	http.ListenAndServe(listenFullHost, nil)
}
