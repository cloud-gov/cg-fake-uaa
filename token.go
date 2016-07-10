package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
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
	// TODO: Ensure 'client_secret' is in POST args.
	// TODO: Ensure 'grant_type' is 'authorization_code'.
	// TODO: Ensure 'response_type' is 'token'.

	errBadRequest := func (message string) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		w.Write([]byte(message))
	}

	email := r.PostFormValue("code")
	clientId := r.PostFormValue("client_id")

	if email == "" {
		errBadRequest("'code' is missing or empty")
		return
	}

	if clientId == "" {
		errBadRequest("'client_id' is missing or empty")
		return
	}

	tokenDuration, err := time.ParseDuration("12h")
	if (err != nil) {
		panic("Unable to parse duration!")
	}
	tokenDurationSeconds := int64(tokenDuration.Seconds())

	authTime := time.Now().Unix()

	// TODO: Some of these fields have been hard-coded. Might be better to use "real" values.

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": []string{
			"openid",
			clientId,
		},
		"auth_time": authTime,
		"azp": clientId,
		"cid": clientId,
		"client_id": clientId,
		"email": email,
		"exp": authTime + tokenDurationSeconds,
		"grant_type": "authorization_code",
		"iat": authTime,
		"iss": "https://uaa.cloud.gov/oauth/token",
		"jti": "fake_jti",
		"origin": "gsa.gov",
		"rev_sig": "9ad72122",
		"scope": []string{"openid"},
		"sub": "12345678-1234-1234-1234-123456789abc",
		"user_id": "12345678-1234-1234-1234-123456789abc",
		"user_name": email,
		"zid": "uaa",
	})

	// The client won't need to verify this because it will be communicating
	// directly with the ID provider (i.e., us) over an intermediary-free
	// trusted channel, using its client secret to authenticate with us.
	//
	// https://developers.google.com/identity/protocols/OpenIDConnect#obtainuserinfo

	accessTokenString, err := accessToken.SignedString([]byte("unused secret key (for verification)"))
	if (err != nil) {
		panic(fmt.Sprintf("Unable to create JSON web token! %v", err))
	}

	str, err := json.Marshal(tokenResponse{
		AccessToken:  accessTokenString,
		ExpiresIn:    tokenDurationSeconds,
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
