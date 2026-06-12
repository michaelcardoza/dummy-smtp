package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed dist
var dist embed.FS

type Handler struct {
	broker     *Broker
	fileServer http.Handler
}

func NewHandler(broker *Broker) (*Handler, error) {
	distFS, _ := fs.Sub(dist, "dist")
	return &Handler{
		broker:     broker,
		fileServer: http.FileServer(http.FS(distFS)),
	}, nil
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.serverSPA)
	mux.HandleFunc("GET /events", h.Events)

	return mux
}

func (h *Handler) serverSPA(w http.ResponseWriter, r *http.Request) {
	h.fileServer.ServeHTTP(w, r)
}

func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := h.broker.Subscribe()
	defer h.broker.Unsubscribe(ch)

	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case m := <-ch:
			data, err := json.Marshal(m)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", data)
			flusher.Flush()
		}
	}
}
