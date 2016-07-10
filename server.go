package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type ServerConfig struct {
	CallbackUrl *url.URL
}

func NewServerHandler(config *ServerConfig) func(http.ResponseWriter, *http.Request) {
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
		data, err := GetAsset("data/fake-cloud.gov.svg")
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
