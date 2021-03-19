package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k88t76/CodeArchives-server/models"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var user models.User
	json.Unmarshal(body, &user)
	if user.Name == "" {
		return
	}
	if user.Validate() {
		err := user.Create()
		token, err := user.CreateSession()
		fmt.Println(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		output, _ := json.MarshalIndent(&token, "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	} else {
		output, _ := json.MarshalIndent("UsedName", "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(output)
	}
}
