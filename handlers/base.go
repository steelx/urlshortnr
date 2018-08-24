package handlers

import (
	"encoding/json"
	"net/http"
)

func createResponse(w http.ResponseWriter, r Response) {
	d, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	w.Write(d)
}
