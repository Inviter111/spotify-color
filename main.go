package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

type colorResponse struct {
	URL        string `json:"url"`
	CoverColor string `json:"hex"`
}

const redirectURL = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate, spotify.ScopeUserReadCurrentlyPlaying)
	ch    = make(chan *spotify.Client)
	state = "test"
)

func main() {
	var client *spotify.Client
	clientID := "077c4f0ef278439cbc51d698b2fbd7de"
	secret := "e692f290d38d4f16b7def9e9f1b196b2"

	http.HandleFunc("/player", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		url := auth.AuthURL(state)
		http.Redirect(w, r, url, 301)
	})
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/imageColor", func(w http.ResponseWriter, r *http.Request) {
		curr, err := client.PlayerCurrentlyPlaying()
		if err != nil {
			printError(err)
		}
		imgURL := curr.Item.Album.Images[0].URL
		smallImgURL := curr.Item.Album.Images[1].URL
		hex := getImageColor(smallImgURL)
		resp := colorResponse{
			URL:        imgURL,
			CoverColor: hex,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	go func() {
		auth.SetAuthInfo(clientID, secret)
		url := auth.AuthURL(state)
		fmt.Println("Please login to Spotify on this page:", url)

		client = <-ch

		user, err := client.CurrentUser()
		if err != nil {
			printError(err)
		}

		fmt.Println("You logged in as:", user.ID)
	}()

	http.ListenAndServe(":8080", nil)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)
	if err != nil {
		printError(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := auth.NewClient(token)
	// fmt.Fprint(w, "Login success")
	http.Redirect(w, r, "http://localhost:8080/player", 301)
	ch <- &client
}
