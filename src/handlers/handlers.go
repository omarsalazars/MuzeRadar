package handlers

import (
	"MuzeRadar/spotify"
	httputilcustom "MuzeRadar/util"
	"MuzeRadar/vars"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func SpotifyRequestAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET params were: ", r.URL.Query())
	err := r.URL.Query().Get("error")
	if err != "" {
		httputilcustom.HandleError(errors.New("error encontrado: " + err))
	}

	code := r.URL.Query().Get("code")

	fmt.Printf("code = %s\n", code)

	access_token := spotify.SpotifyRequestAccessToken(code)

	http.SetCookie(w, &http.Cookie{
		Name:     vars.Token_cookie_name,
		Value:    access_token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
}

func SpotifyGetUsersTopArtistsHandler(w http.ResponseWriter, r *http.Request) {
	requestURL := "https://api.spotify.com/v1/me/top/artists?"

	values := url.Values{}
	values.Add("time_range", "medium_term")
	values.Add("limit", "25")
	values.Add("offset", "0")

	query := values.Encode()

	req, err := http.NewRequest(http.MethodGet, requestURL+query, nil)
	httputilcustom.HandleError(err)

	vars.Cookie_value = httputilcustom.GetCookie(r, vars.Token_cookie_name)
	req.Header.Add("Authorization", "Bearer "+vars.Cookie_value)

	httputilcustom.PrintRequest(req)

	res, err := http.DefaultClient.Do(req)
	httputilcustom.HandleError(err)

	// printResponse(res)
	artists := spotify.SpotifyParseArtistJson(res)
	for i := range artists {
		fmt.Println(artists[i].Name)
	}

	spotify.MakeRecommendations(artists)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "This is root page %s", r.URL.Path[1:])
	spotify.SpotifyAuthorization(w, r)
}

func ValidateAuthorizationToken(token string) bool {
	return true
	// ToDo
}

func ValidateAuthorizationTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := httputilcustom.GetCookie(r, vars.Token_cookie_name)
		valid_token := ValidateAuthorizationToken(token)
		if valid_token {
			http.Redirect(w, r, "http://localhost:8080/top", http.StatusAccepted)
		} else {
			next(w, r)
		}
	}
}
