package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/k88t76/CodeArchives-server/server/models"
)

var cookie http.Cookie

func StartWebServer() {
	fmt.Println("Server 8080")
	/*server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	*/
	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/archive/", handleRequest)
	http.HandleFunc("/archives", hR)
	http.HandleFunc("/archive/c", hRc)
	http.HandleFunc("/edit/", hRe)
	http.HandleFunc("/delete/", hRd)
	http.HandleFunc("/search/", hRs)
	http.HandleFunc("/signin", hRsignIn)
	http.HandleFunc("/signup", hRsignUp)
	http.HandleFunc("/signout", hRsignOut)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

func hR(w http.ResponseWriter, r *http.Request) {
	err := handleGetAll(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func hRc(w http.ResponseWriter, r *http.Request) {
	err := handlePost(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func hRe(w http.ResponseWriter, r *http.Request) {
	err := handlePut(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func hRd(w http.ResponseWriter, r *http.Request) {
	err := handleDelete(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func hRs(w http.ResponseWriter, r *http.Request) {
	err := handleSearch(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func hRsignIn(w http.ResponseWriter, r *http.Request) {
	err := handleSignIn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func hRsignUp(w http.ResponseWriter, r *http.Request) {
	err := handleSignUp(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func hRsignOut(w http.ResponseWriter, r *http.Request) {
	err := handleSignOut(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGetAll(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("cookie: %v\n", cookie)
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

func handleGet(w http.ResponseWriter, r *http.Request) error {
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

func handleSearch(w http.ResponseWriter, r *http.Request) error {
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

func handleSignIn(w http.ResponseWriter, r *http.Request) error {
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
	}
	return nil
}

func handleSignUp(w http.ResponseWriter, r *http.Request) error {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var user models.User
	json.Unmarshal(body, &user)
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
	}
	return nil
}

func handleSignOut(w http.ResponseWriter, r *http.Request) error {
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

func handlePost(w http.ResponseWriter, r *http.Request) error {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
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

func handlePut(w http.ResponseWriter, r *http.Request) error {
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

func handleDelete(w http.ResponseWriter, r *http.Request) error {
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
