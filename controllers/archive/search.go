package archive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/k88t76/CodeArchives-server/models"
)

func Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var token string
	json.Unmarshal(body, &token)
	if token == "" {
		return
	}
	name, err1 := models.GetUserNameByToken(token)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
	word := path.Base(r.URL.Path)
	fmt.Println(word)
	var archives []models.Archive
	if word == "" {
		archives, _ = models.GetArchivesByUser(name, 1000)
	} else {
		archives, _ = models.GetMatchArchive(word, name)
	}
	output, err2 := json.MarshalIndent(&archives, "", "\t\t")
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
