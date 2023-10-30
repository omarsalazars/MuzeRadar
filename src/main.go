package main

import (
	"MuzeRadar/handlers"
	"MuzeRadar/vars"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	vars.Client_id = os.Getenv("SPOTIFY_CLIENT_ID")
	vars.Client_secret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	fmt.Println("Running! localhost:8080")
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/authorized", handlers.SpotifyRequestAccessTokenHandler)
	http.HandleFunc("/top", handlers.SpotifyGetUsersTopArtistsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
