package archive

import (
	"fmt"
	"net/http"
	"path"

	"github.com/k88t76/CodeArchives-server/models"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	uuid := path.Base(r.URL.Path)
	if uuid == "" {
		return
	}
	archive := models.GetArchive(uuid)
	fmt.Println(archive)
	err := archive.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	return
}
