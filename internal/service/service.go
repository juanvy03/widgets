package service

import (
	"log"
	"net/http"
	"spike/internal/api"
)

func InitHttpServer() {
	http.HandleFunc("/add", api.AddWidget)
	http.HandleFunc("/remove", api.RemoveWidget)
	http.HandleFunc("/link", api.LinkWidget)

	http.ListenAndServe(":8090", nil)

	log.Println("HTTP Server listening to :8090")
}
