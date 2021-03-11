package archive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/k88t76/CodeArchives-server/models"
)

func Edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	uuid := path.Base(r.URL.Path)
	fmt.Printf("Edit uuid: %v\n", uuid)
	archive := models.GetArchive(uuid)
	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)
	if len(body) == 0 {
		return
	}
	fmt.Printf("body: %v\n", body)
	json.Unmarshal(body, &archive)
	fmt.Printf("archive: %v", archive)
	err := archive.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	return
}
