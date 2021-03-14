package user

import (
	"encoding/json"
	"net/http"

	"github.com/k88t76/CodeArchives-server/models"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	if r.Method == "OPTIONS" {
		return
	}
	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)
	var user models.User
	json.Unmarshal(body, &user)
	uID, checkUser, checkPassword := models.CheckUser(user)
	if checkUser && checkPassword {
		token, err := models.UpdateToken(uID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		output, _ := json.MarshalIndent(&token, "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	} else if !checkUser {
		output, _ := json.MarshalIndent("Unknown User", "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(output)
	} else {
		output, _ := json.MarshalIndent("Wrong Password", "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(output)
	}
}
