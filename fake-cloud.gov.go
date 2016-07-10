package main

import (
	"fmt"
	"flag"
	"encoding/json"
	"net/http"
	"net/url"
)

type ServerConfig struct {
	CallbackUrl *url.URL
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Jti          string `json:"jti"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func ExchangeCodeForAccessToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Ensure 'code' is in POST args.
	// TODO: Ensure 'client_id' is in POST args.
	// TODO: Ensure 'client_secret' is in POST args.
	// TODO: Ensure 'grant_type' is 'authorization_code'.
	// TODO: Ensure 'reponse_type' is 'token'.

	email := r.FormValue("code")
	accessToken := fmt.Sprintf("TODO: build jwt access token for %s", email)
	str, err := json.Marshal(TokenResponse{
		AccessToken:  accessToken,
		ExpiresIn:    1, // TODO: Actually provide a useful value here.
		Jti:          "fake_jti",
		RefreshToken: "fake_oauth2_refresh_token",
		Scope:        "openid",
		TokenType:    "bearer",
	})
	if err != nil {
		panic("Unable to encode JSON!")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(str)
}

func NewHandler(config *ServerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		BaseHandler(config, w, r)
	}
}

func BaseHandler(config *ServerConfig, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == Urls.Reverse("authorize") && r.Method == "GET" {
		Authorize(config, w, r)
	} else if r.URL.Path == Urls.Reverse("token") && r.Method == "POST" {
		ExchangeCodeForAccessToken(w, r)
	} else if r.URL.Path == Urls.Reverse("svgLogo") {
		data, err := Asset("data/fake-cloud.gov.svg")
		if err != nil {
			panic("Couldn't find fake-cloud.gov.svg!")
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		fmt.Fprintf(w, "%s", data)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(404)
		fmt.Fprintf(w, "Not Found")
	}
}

func main() {
	var callbackUrl string

	portPtr := flag.Int("port", 8080, "Port to listen on")
	flag.StringVar(&callbackUrl, "callback-url", "http://localhost:8000/callback", "OAuth2 Callback URL")

	flag.Parse()

	handler := NewHandler(&ServerConfig{
		CallbackUrl: Urlify(callbackUrl),
	})

	fmt.Printf("OAuth2 callback URL: %s\n", callbackUrl)
	fmt.Printf("OAuth2 authorize URL: http://localhost:%d%s\n", *portPtr, Urls.Reverse("authorize"))
	fmt.Printf("OAuth2 token URL: http://localhost:%d%s\n", *portPtr, Urls.Reverse("token"))
	fmt.Printf("\nListening on port %d.\n", *portPtr)

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), nil)
}
