package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Jti          string `json:"jti"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func SendBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(400)
	w.Write([]byte(message))
}

func HandleTokenRequest(w http.ResponseWriter, r *http.Request) {
	clientId := r.PostFormValue("client_id")
	grant_type := r.PostFormValue("grant_type")

	if clientId == "" {
		SendBadRequest(w, "'client_id' is missing or empty")
		return
	}

	if r.PostFormValue("client_secret") == "" {
		SendBadRequest(w, "'client_secret' is missing or empty")
		return
	}

	if grant_type == "authorization_code" {
		ExchangeCodeForAccessToken(w, r, clientId)
	} else if grant_type == "refresh_token" {
		RefreshAccessToken(w, r, clientId)
	} else {
		SendBadRequest(w, "'grant_type' must be 'authorization_code' or 'refresh_token'")
	}
}

func RefreshAccessToken(w http.ResponseWriter, r *http.Request, clientId string) {
	refresh_token := r.PostFormValue("refresh_token")
	parts := strings.SplitN(refresh_token, ":", 2)

	if (parts[0] != "fake_oauth2_refresh_token") {
		SendBadRequest(w, "'refresh_token' is missing or malformed")
		return
	}

	email := parts[1]

	SendAccessToken(w, clientId, email)
}

func ExchangeCodeForAccessToken(w http.ResponseWriter, r *http.Request, clientId string) {
	email := r.PostFormValue("code")

	if email == "" {
		SendBadRequest(w, "'code' is missing or empty")
		return
	}

	if r.PostFormValue("response_type") != "token" {
		SendBadRequest(w, "'response_type' is expected to be 'token'")
		return
	}

	SendAccessToken(w, clientId, email)
}

func SendAccessToken(w http.ResponseWriter, clientId string, email string) {
	tokenDuration, err := time.ParseDuration("10m")
	if err != nil {
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
		"auth_time":  authTime,
		"azp":        clientId,
		"cid":        clientId,
		"client_id":  clientId,
		"email":      email,
		"exp":        authTime + tokenDurationSeconds,
		"grant_type": "authorization_code",
		"iat":        authTime,
		"iss":        "https://uaa.cloud.gov/oauth/token",
		"jti":        "fake_jti",
		"origin":     "gsa.gov",
		"rev_sig":    "9ad72122",
		"scope":      []string{"openid"},
		"sub":        "12345678-1234-1234-1234-123456789abc",
		"user_id":    "12345678-1234-1234-1234-123456789abc",
		"user_name":  email,
		"zid":        "uaa",
	})

	// The client won't need to verify this because it will be communicating
	// directly with the ID provider (i.e., us) over an intermediary-free
	// trusted channel, using its client secret to authenticate with us.
	//
	// https://developers.google.com/identity/protocols/OpenIDConnect#obtainuserinfo

	accessTokenString, err := accessToken.SignedString([]byte("unused secret key (for verification)"))
	if err != nil {
		panic(fmt.Sprintf("Unable to create JSON web token! %v", err))
	}

	str, err := json.Marshal(tokenResponse{
		AccessToken:  accessTokenString,
		ExpiresIn:    tokenDurationSeconds,
		Jti:          "fake_jti",
		RefreshToken: fmt.Sprintf("fake_oauth2_refresh_token:%s", email),
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
