package spotify

import (
	httputilcustom "MuzeRadar/util"
	"MuzeRadar/vars"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var redirect_uri string = "http://localhost:8080/authorized"

type Image struct {
	Height int
	Width  int
	Url    string
}

type Artist struct {
	Id         string
	Genres     []string
	Images     []Image
	Name       string
	Popularity int
	Uri        string
}

type TopResponse struct {
	Items    []Artist
	Total    int
	Limit    int
	Offset   int
	Href     string
	Next     string
	Previous string
}

type Track struct {
	Artists    []Artist
	Name       string
	Popularity int
	URI        string
}

type RecommendationResponse struct {
	Tracks []Track
}

type CurrentUserPlaylistsResponse struct {
	Items []Playlist
}

type Playlist struct {
	Id   string
	Name string
}

type User struct {
	Id string
}

func SpotifyAuthorization(w http.ResponseWriter, r *http.Request) {

	requestURL := "https://accounts.spotify.com/authorize?"

	scope := "playlist-read-collaborative "
	scope += "playlist-modify-public "
	scope += "user-top-read "
	scope += "user-read-recently-played "
	scope += "playlist-modify-private"

	values := url.Values{}
	values.Add("client_id", vars.Client_id)
	values.Add("response_type", "code")
	values.Add("redirect_uri", redirect_uri)
	values.Add("scope", scope)
	query := values.Encode()
	http.Redirect(w, r, requestURL+query, http.StatusSeeOther)
}

func SpotifyRequestAccessToken(code string) string {
	requestURL := "https://accounts.spotify.com/api/token"

	query := url.Values{}
	query.Add("code", code)
	query.Add("redirect_uri", redirect_uri)
	query.Add("grant_type", "authorization_code")

	req, err := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(query.Encode()))
	httputilcustom.HandleError(err)

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(vars.Client_id+":"+vars.Client_secret)))

	fmt.Println(req.Header.Get("Authorization"))

	// printRequest(req)

	res, err := http.DefaultClient.Do(req)
	httputilcustom.HandleError(err)

	// printResponse(res)

	if res.StatusCode != 200 {
		fmt.Println("Something went wrong, status ", res.StatusCode)
	}

	parsedResponse := httputilcustom.ParseJsonResponseAsMap(res)
	access_token := parsedResponse["access_token"]
	fmt.Println(access_token)
	return access_token
}

func MakeRecommendations(artists []Artist) {
	artistsIds := ""
	for i := 0; i < 5; i++ {
		artistsIds += artists[i].Id
		artistsIds += ","
	}
	artistsIds = artistsIds[:len(artistsIds)-1]

	fmt.Println("ARTIST IDS:", artistsIds)
	requestURL := "https://api.spotify.com/v1/recommendations?"
	values := url.Values{}
	values.Add("seed_artists", artistsIds)
	// values.Add("seed_genres", )
	values.Add("market", "ES")
	values.Add("limit", "50")
	values.Add("min_popularity", "0")
	values.Add("max_popularity", "50")
	values.Add("target_popularity", "25")
	query := values.Encode()

	req := httputilcustom.GetAuthorizedRequest(http.MethodGet, requestURL+query, nil)

	res, err := http.DefaultClient.Do(req)
	httputilcustom.HandleError(err)

	// printResponse(res)
	tracks := SpotifyParseRecommendationsJson(res)
	for i := range tracks {
		fmt.Print(tracks[i].Name, " ")
		for j := range tracks[i].Artists {
			fmt.Print(tracks[i].Artists[j].Name, " ")
		}
		fmt.Println(tracks[i].Popularity)
	}

	playlist := SpotifyCheckIfPlaylistExists(SpotifyGetUserId())
	SpotifyAddTracksToPlaylist(playlist.Id, tracks)
}

func SpotifyGetUserId() string {
	requestURL := "https://api.spotify.com/v1/me"
	var user User

	req := httputilcustom.GetAuthorizedRequest(http.MethodGet, requestURL, nil)

	res, err := http.DefaultClient.Do(req)
	httputilcustom.HandleError(err)

	body, err := ioutil.ReadAll(res.Body)
	httputilcustom.HandleError(err)
	res.Body.Close()

	err = json.Unmarshal(body, &user)
	httputilcustom.HandleError(err)
	return user.Id
}

func SpotifyCreatePlaylist(userID string) Playlist {
	requestURL := "https://api.spotify.com/v1/users/" + userID + "/playlists"

	values := []byte(fmt.Sprintf(`{
		"name" : "%v",
		"public" : "true",
		"collaborative" : "false",
		"description" : "%v"
	}`, vars.Playlist_name, vars.Playlist_description))

	req := httputilcustom.GetAuthorizedRequest(http.MethodPost, requestURL, strings.NewReader(string(values)))
	req.Header.Add("content-type", "application/json")

	httputilcustom.PrintRequest(req)
	res, err := http.DefaultClient.Do(req)
	httputilcustom.HandleError(err)
	httputilcustom.PrintResponse(res)

	body, err := ioutil.ReadAll(res.Body)
	httputilcustom.HandleError(err)
	res.Body.Close()

	var playlist Playlist
	err = json.Unmarshal(body, &playlist)
	httputilcustom.HandleError(err)
	return playlist
}

func SpotifyAddTracksToPlaylist(playlistID string, tracks []Track) {
	requestURL := "https://api.spotify.com/v1/playlists/" + playlistID + "/tracks"

	trackUris := ""
	for i := 0; i < len(tracks); i++ {
		trackUris += "\"" + tracks[i].URI + "\","
	}
	trackUris = trackUris[:len(trackUris)-1]

	values := []byte(fmt.Sprintf(`{
		"uris" : [ %v ]
	}`, trackUris))

	req := httputilcustom.GetAuthorizedRequest(http.MethodPost, requestURL, strings.NewReader(string(values)))
	req.Header.Add("content-type", "application/json")

	httputilcustom.PrintRequest(req)
	res, err := http.DefaultClient.Do(req)
	httputilcustom.HandleError(err)
	httputilcustom.PrintResponse(res)
}

func SpotifyGetUserPlaylists() []Playlist {
	requestURL := "https://api.spotify.com/v1/me/playlists"

	req := httputilcustom.GetAuthorizedRequest(http.MethodGet, requestURL, nil)
	httputilcustom.PrintRequest(req)
	res, err := http.DefaultClient.Do(req)
	httputilcustom.HandleError(err)
	body, err := ioutil.ReadAll(res.Body)
	httputilcustom.HandleError(err)
	res.Body.Close()
	var currentPlaylists CurrentUserPlaylistsResponse
	err = json.Unmarshal(body, &currentPlaylists)
	httputilcustom.HandleError(err)
	return currentPlaylists.Items
}

// Returns playlist Id if exists, creates playlist otherwise
func SpotifyCheckIfPlaylistExists(name string) Playlist {
	playlists := SpotifyGetUserPlaylists()
	for i := range playlists {
		if playlists[i].Name == name {
			return playlists[i]
		}
	}
	return SpotifyCreatePlaylist(SpotifyGetUserId())
}

func SpotifyParseRecommendationsJson(res *http.Response) []Track {
	var recommends RecommendationResponse
	body, err := ioutil.ReadAll(res.Body)
	httputilcustom.HandleError(err)
	res.Body.Close()
	err = json.Unmarshal(body, &recommends)
	httputilcustom.HandleError(err)
	return recommends.Tracks
}

func SpotifyParseArtistJson(res *http.Response) []Artist {
	var artists []Artist
	var top_response TopResponse
	body, err := ioutil.ReadAll(res.Body)
	httputilcustom.HandleError(err)
	res.Body.Close()
	err = json.Unmarshal(body, &top_response)
	httputilcustom.HandleError(err)
	artists = top_response.Items
	return artists
}
