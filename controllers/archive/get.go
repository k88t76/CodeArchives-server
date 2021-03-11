package archive

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/k88t76/CodeArchives-server/models"
)

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	uuid := path.Base(r.URL.Path)
	archive := models.GetArchive(uuid)
	output, err := json.MarshalIndent(&archive, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}
