package server

import (
	"log"
	"net/http"
	"os"

	"github.com/k88t76/CodeArchives-server/controllers/archive"
	"github.com/k88t76/CodeArchives-server/controllers/user"
)

func StartWebServer() {
	http.HandleFunc("/archive/", archive.Get)
	http.HandleFunc("/archives", archive.GetAll)
	http.HandleFunc("/create", archive.Create)
	http.HandleFunc("/edit/", archive.Edit)
	http.HandleFunc("/delete/", archive.Delete)
	http.HandleFunc("/search/", archive.Search)
	http.HandleFunc("/signin", user.SignIn)
	http.HandleFunc("/signup", user.SignUp)
	http.HandleFunc("/userbytoken", user.UserByToken)
	http.HandleFunc("/guestsignin", user.GuestSignIn)

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
