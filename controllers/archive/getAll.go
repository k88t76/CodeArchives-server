package archive

import (
	"encoding/json"
	"net/http"

	"github.com/k88t76/CodeArchives-server/models"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var token string
	json.Unmarshal(body, &token)
	name, _ := models.GetUserNameByToken(token)
	archives, _ := models.GetArchivesByUser(name, 1000)
	output, err := json.MarshalIndent(&archives, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
