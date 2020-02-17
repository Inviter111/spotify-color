package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"spotify-color/colorpicker"
)

type colorResponse struct {
	CoverColor string `json:"hex"`
}

func main() {
	http.HandleFunc("/getColor", func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		url, ok := r.URL.Query()["url"]
		if !ok || len(url[0]) < 1 {
			fmt.Println("URL is missing")
			return
		}

		hex := colorpicker.GetImageColor(url[0])
		resp := colorResponse{
			CoverColor: hex,
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			printError(err)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		printError(err)
	}
}

func printError(err error) {
	log.Fatalln("Error:", err)
	os.Exit(1)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
