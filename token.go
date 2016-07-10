package main

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type tokenResponse struct {
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
	str, err := json.Marshal(tokenResponse{
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
