package main

import (
	"log"
	"net/http"
	"spotify-color/ws"
)

func main() {
	PORT := "8080"

	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	log.Print("Server is listening on port ", PORT)
	http.ListenAndServe(":"+PORT, nil)
	// log.Fatal(fasthttp.ListenAndServe(":"+PORT, router.RequestHandler))
}
