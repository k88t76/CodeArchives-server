package archive

import (
	"encoding/json"
	"net/http"

	"github.com/k88t76/CodeArchives-server/models"
)

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var archive models.Archive
	json.Unmarshal(body, &archive)
	err := archive.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	return
}
