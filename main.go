package main

import (
	"log"
	"spotify-color/router"
	"spotify-color/ws"

	"github.com/valyala/fasthttp"
)

func main() {
	PORT := "8080"

	hub := ws.NewHub()
	go hub.Run()

	log.Print("Server is listening on port ", PORT)
	// http.ListenAndServe(":"+PORT, nil)
	log.Fatal(fasthttp.ListenAndServe(":"+PORT, func(ctx *fasthttp.RequestCtx) {
		router.RequestHandler(ctx, hub)
	}))
}
