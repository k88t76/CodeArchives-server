package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/k88t76/CodeArchives-server/models"
)

var cookie http.Cookie

func StartWebServer() {
	http.HandleFunc("/archive/", handleRequest)
	http.HandleFunc("/archives", handleGetAll)
	http.HandleFunc("/archive/c", handleCreate)
	http.HandleFunc("/edit/", handleEdit)
	http.HandleFunc("/delete/", handleDelete)
	http.HandleFunc("/search/", handleSearch)
	http.HandleFunc("/signin", handleSignIn)
	http.HandleFunc("/signup", handleSignUp)
	http.HandleFunc("/signout", handleSignOut)

	http.HandleFunc("/testsignin", handleTestSignIn)

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

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = Get(w, r)
	case "POST":
		//err = handlePost(w, r)
	case "PUT":
		//err = handlePut(w, r)
	case "DELETE":
		//err = handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGetAll(w http.ResponseWriter, r *http.Request) {

	err := GetAll(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	err := Create(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleEdit(w http.ResponseWriter, r *http.Request) {
	err := Edit(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	err := Delete(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	err := Search(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	err := SignIn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	err := SignUp(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSignOut(w http.ResponseWriter, r *http.Request) {
	err := SignOut(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleTestSignIn(w http.ResponseWriter, r *http.Request) {
	err := TestSignIn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("cookie: %v\n", cookie)
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	if cookie.Value != "" {
		fmt.Println("cookie_detect!")
		userName := models.GetUserNameBySessionID(cookie.Value)
		archives, _ := models.GetArchivesByUser(userName, 100)
		output, _ := json.MarshalIndent(&archives, "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	} else {
		output, _ := json.MarshalIndent("UnLogin", "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	}

	return nil
}

func Get(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	uuid := path.Base(r.URL.Path)

	archive := models.GetArchive(uuid)
	output, err := json.MarshalIndent(&archive, "", "\t\t")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return nil
}

func Search(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	userName := models.GetUserNameBySessionID(cookie.Value)
	word := path.Base(r.URL.Path)
	fmt.Printf("serch: %v\n", word)
	archives, _ := models.GetMatchArchive(word, userName)
	output, err := json.MarshalIndent(&archives, "", "\t\t")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return nil
}

func SignIn(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var user models.User
	json.Unmarshal(body, &user)
	fmt.Println("userName:")
	fmt.Println(user.Name)
	u, check := models.CheckUser(user)
	if check {
		fmt.Println("check OK")
		err := u.CreateSession()
		fmt.Printf("u : %v\n", u)
		if err != nil {
			return err
		}
		session := models.GetSession(u.UUID)
		fmt.Printf("session: %v\n", session)
		cookie = http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(200)
		fmt.Printf("cookie Value: %v\n", cookie.Value)
	} else {
		output, _ := json.MarshalIndent("Failed SignIn", "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	}
	return nil
}

func SignUp(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Printf("Post body: %v\n", body)
	fmt.Printf("r.Method: %v\n", r.Method)
	fmt.Printf("r.Body: %v\n", r.Body)
	var user models.User
	json.Unmarshal(body, &user)
	fmt.Println(user.Name)
	fmt.Printf("user: %v\n", user)
	err := user.Create()
	u, check := models.CheckUser(user)
	if check && err == nil {
		fmt.Println("check OK")
		err := u.CreateSession()
		if err != nil {
			return err
		}
		session := models.GetSession(u.UUID)
		fmt.Printf("session: %v\n", session)
		cookie = http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(200)
		fmt.Printf("cookie Value: %v\n", cookie.Value)
	} else {
		output, _ := json.MarshalIndent("UnLogin", "", "\t\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	}
	return nil
}

func SignOut(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	cookie = http.Cookie{
		Name:     "_cookie",
		Value:    "",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(200)
	fmt.Printf("cookie Value: %v\n", cookie.Value)
	return nil
}

func Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Printf("Post body: %v\n", body)
	var archive models.Archive
	json.Unmarshal(body, &archive)
	userName := models.GetUserNameBySessionID(cookie.Value)
	archive.Author = userName
	err := archive.Create()
	if err != nil {
		return err
	}
	w.WriteHeader(200)
	return nil
}

func Edit(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	uuid := path.Base(r.URL.Path)
	fmt.Printf("[Edit] uuid: %v\n", uuid)
	archive := models.GetArchive(uuid)
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Printf("body: %v\n", body)
	json.Unmarshal(body, &archive)
	fmt.Printf("archive: %v", archive)
	err := archive.Update()
	if err != nil {
		return err
	}
	w.WriteHeader(200)
	return nil
}

func Delete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	uuid := path.Base(r.URL.Path)
	fmt.Printf("uuid: %v\n", uuid)
	archive := models.GetArchive(uuid)
	fmt.Println(archive)
	err := archive.Delete()
	if err != nil {
		return err
	}
	w.WriteHeader(200)
	return nil
}

func TestSignIn(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	archives, _ := models.GetTestArchives()
	output, _ := json.MarshalIndent(&archives, "", "\t\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return nil
}
