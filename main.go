package main

import (
	"log"
	"net/http"

	"github.com/wilsonle/excalidraw-collab/controllers"
	"github.com/wilsonle/excalidraw-collab/middleware"
	"github.com/wilsonle/excalidraw-collab/pkg/database"
)

func main() {
	database.InitDB("./excalidraw-collab.db")

	http.HandleFunc("/login", controllers.LoginHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", middleware.BasicAuth(fs))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
