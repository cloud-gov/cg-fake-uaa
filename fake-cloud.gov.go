package main

import (
	"fmt"
	"flag"
	"net/http"
)

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
