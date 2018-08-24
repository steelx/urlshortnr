// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/steelx/urlshortnr/storages"
)

type response struct {
	Data    interface{} `json:"response"`
	Success bool        `json:"success"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// New returns Http Handlers for url shortner app
func New(prefix string, storage storages.IFStorage) http.Handler {
	mux := http.NewServeMux()
	h := handler{prefix, storage}

	mux.HandleFunc("/encode/", responseHandler(h.EncodeHandler))
	mux.HandleFunc("/info/", responseHandler(h.DecodeHandler))
	mux.HandleFunc("/", h.RedirectHandler)

	return mux
}

type handler struct {
	prefix  string
	storage storages.IFStorage
}

func responseHandler(h func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}
		w.WriteHeader(status)

		err = json.NewEncoder(w).Encode(response{data, err == nil})
		if err != nil {
			fmt.Printf("Could not encode response: %v", err)
		}
	}
}

func (h handler) EncodeHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	enableCors(&w)

	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method %v not allowed", r.Method)
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct{ URL string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to decode JSON request body: %v", err.Error())
	}

	URL := strings.TrimSpace(body.URL)

	if URL == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("URL is Empty")
	}

	c, err := h.storage.Save(URL)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Could not store in database: %v", err.Error())
	}

	return h.prefix + c, http.StatusCreated, nil
}

func (h handler) DecodeHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method %v not allowed", r.Method)
	}

	w.Header().Set("Content-Type", "application/json")
	code := r.URL.Path[len("/info/"):]

	model, err := h.storage.LoadInfo(code)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("URL Not Found")
	}

	return model, http.StatusOK, nil
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
