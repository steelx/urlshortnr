// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/steelx/urlshortnr/storages"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// New returns Http Handlers for url shortner app
func New(prefix string, storage storages.IFStorage) http.Handler {
	mux := http.NewServeMux()
	h := handler{prefix, storage}

	mux.HandleFunc("/encode/", h.EncodeHandler)
	mux.HandleFunc("/", h.RedirectHandler)
	mux.HandleFunc("/info/", h.DecodeHandler)

	return mux
}

type handler struct {
	prefix  string
	storage storages.IFStorage
}

func (h handler) EncodeHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var b struct{ URL string }

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e := response{Data: "Unable to decode JSON request body: " + err.Error(), Success: false}
		createResponse(w, e)
		return
	}

	b.URL = strings.TrimSpace(b.URL)

	if b.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		e := response{Data: "URL is Empty", Success: false}
		createResponse(w, e)
		return
	}

	c, err := h.storage.Save(b.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e := response{Data: err.Error(), Success: false}
		createResponse(w, e)
		return
	}

	response := response{Data: h.prefix + c, Success: true}
	createResponse(w, response)
}

func (h handler) DecodeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	code := r.URL.Path[len("/info/"):]

	model, err := h.storage.LoadInfo(code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		e := response{Data: "URL Not Found", Success: false}
		createResponse(w, e)
		return
	}

	response := response{Data: model, Success: true}
	createResponse(w, response)
}

func (h handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := r.URL.Path[len("/"):]

	model, err := h.storage.Load(code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL Not Found"))
		return
	}

	http.Redirect(w, r, string(model.Url), 301)
}
