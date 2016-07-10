package main

import (
	"fmt"
	"flag"
	"net/http"
	"net/url"
)

type ServerConfig struct {
	CallbackUrl *url.URL
}

func NewHandler(config *ServerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		baseHandler(config, w, r)
	}
}

func baseHandler(config *ServerConfig, w http.ResponseWriter, r *http.Request) {
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
