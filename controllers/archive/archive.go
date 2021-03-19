package archive

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/k88t76/CodeArchives-server/models"
)

func HandleArchive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	switch r.Method {
	case "GET":
		getArchive(w, r)
	case "POST":
		postArchive(w, r)
	case "PUT":
		putArchive(w, r)
	case "DELETE":
		deleteArchive(w, r)
	default:
		return
	}

}

func getArchive(w http.ResponseWriter, r *http.Request) {
	uuid := path.Base(r.URL.Path)
	archive, err := models.GetArchive(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	output, err := json.MarshalIndent(&archive, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func postArchive(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var archive models.Archive
	json.Unmarshal(body, &archive)
	uuid, err := archive.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	output, err := json.MarshalIndent(&uuid, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(output)

}

func putArchive(w http.ResponseWriter, r *http.Request) {
	uuid := path.Base(r.URL.Path)
	archive, err := models.GetArchive(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)
	if len(body) == 0 {
		return
	}
	json.Unmarshal(body, &archive)
	uuid, err = archive.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	output, err := json.MarshalIndent(&uuid, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(output)
}

func deleteArchive(w http.ResponseWriter, r *http.Request) {
	uuid := path.Base(r.URL.Path)
	if uuid == "" {
		return
	}
	archive, err := models.GetArchive(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = archive.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
