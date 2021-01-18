package main

import (
	"github.com/k88t76/CodeArchives-server/config"
	"github.com/k88t76/CodeArchives-server/controllers"
	"github.com/k88t76/CodeArchives-server/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	controllers.StartWebServer()
}

/*http.HandleFunc("/", indexHandler)

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

//indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

*/
