package main

import (
	"log"
	"spotify-color/router"

	"github.com/valyala/fasthttp"
)

func main() {
	PORT := "8080"

	log.Print("Server is listening on port ", PORT)
	log.Fatal(fasthttp.ListenAndServe(":"+PORT, router.RequestHandler))
}
