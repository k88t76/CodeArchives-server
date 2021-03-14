package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k88t76/CodeArchives-server/models"
)

func UserByToken(w http.ResponseWriter, r *http.Request) {
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
	name, err := models.GetUserNameByToken(token)
	if err != nil {
		return
	}
	fmt.Println(name)
	output, _ := json.MarshalIndent(&name, "", "\t\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
