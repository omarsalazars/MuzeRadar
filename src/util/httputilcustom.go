package httputilcustom

import (
	"MuzeRadar/vars"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

func PrintRequest(r *http.Request) {
	req, err := httputil.DumpRequest(r, true)
	HandleError(err)
	fmt.Println(string(req))
}

func PrintResponse(r *http.Response) {
	res, err := httputil.DumpResponse(r, true)
	HandleError(err)
	fmt.Println(string(res))
}

func HandleError(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func ParseJsonResponseAsMap(res *http.Response) map[string]string {
	body, err := ioutil.ReadAll(res.Body)
	HandleError(err)
	res.Body.Close()
	var jsonRes map[string]string
	_ = json.Unmarshal(body, &jsonRes)
	return jsonRes
}

func GetAuthorizedRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	HandleError(err)
	req.Header.Add("Authorization", "Bearer "+vars.Cookie_value)
	return req
}

func GetCookie(r *http.Request, cookie_name string) string {
	cookie, err := r.Cookie(cookie_name)
	HandleError(err)
	if cookie == nil || cookie.Value == "" {
		HandleError(errors.New("error: authorization code not found"))
	}
	return cookie.Value
}
