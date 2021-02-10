package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/k88t76/CodeArchives-server/models"
)

func StartWebServer() {
	http.HandleFunc("/", t)
	http.HandleFunc("/archive/", get)
	http.HandleFunc("/archives", getAll)
	http.HandleFunc("/create", create)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/delete/", delete)
	http.HandleFunc("/search/", search)
	http.HandleFunc("/signin", signIn)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/userbytoken", userByToken)
	http.HandleFunc("/testsignin", testSignIn)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

func t(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
	fmt.Fprint(w, time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02T15:04:05+09:00"))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
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
	return
}

func get(w http.ResponseWriter, r *http.Request) {
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

func search(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
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
	return
}

func signIn(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
	fmt.Printf("r.Method: %v\n", r.Method)
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var user models.User
	json.Unmarshal(body, &user)
	uID, check := models.CheckUser(user)
	if check {
		fmt.Println("check OK")
		token, err := models.UpdateToken(uID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		output, _ := json.MarshalIndent(&token, "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	} else {
		output, _ := json.MarshalIndent("Wrong Password", "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	}
	return
}

func signUp(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var user models.User
	json.Unmarshal(body, &user)
	if user.Name == "" {
		return
	}
	fmt.Println(user)
	fmt.Println(user.Validate())
	if user.Validate() {
		fmt.Println("validation OK")
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
		w.Write(output)
	}
	return
}

func create(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
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

func edit(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
	uuid := path.Base(r.URL.Path)
	fmt.Printf("[Edit] uuid: %v\n", uuid)
	archive := models.GetArchive(uuid)
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	if body == nil {
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

func delete(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
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

func userByToken(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
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
	return
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

func testSignIn(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var user models.User
	json.Unmarshal(body, &user)
	uID, check := models.CheckUser(user)
	if check {
		fmt.Println("check OK")
		token, err := models.UpdateToken(uID)
		models.CreateTestArchives()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		output, _ := json.MarshalIndent(&token, "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	} else {
		output, _ := json.MarshalIndent("Wrong Password", "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	}
	return
}
