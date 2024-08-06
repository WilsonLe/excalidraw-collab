package main

import (
	"log"
	"net/http"

	"github.com/wilsonle/excalidraw-collab/controllers"
	"github.com/wilsonle/excalidraw-collab/middleware"
	"github.com/wilsonle/excalidraw-collab/pkg/database"
	"github.com/wilsonle/excalidraw-collab/views"
)

func main() {
	database.InitDB("./excalidraw-collab.db")

	http.Handle("/", middleware.BasicAuth(http.HandlerFunc(views.RootHandler)))
	http.HandleFunc("/login", controllers.LoginHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
