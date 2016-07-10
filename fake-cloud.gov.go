package main

import (
	"fmt"
	"flag"
	"html/template"
	"encoding/json"
	"net/http"
	"net/url"
	"math/rand"
)

var taglines = [...]string{
	`Welcome to your fake <abbr title="User Account and Authentication">UAA</abbr> provider!`,
	`The convenience of zero-factor authentication is here.`,
	`We're like <a href="https://cloud.gov">cloud.gov</a>, but without the security.`,
}

type UrlMap map[string]string

var urls = UrlMap{
	"authorize": "/oauth/authorize",
	"token": "/oauth/token",
	"svgLogo": "/fake-cloud.gov.svg",
}

type LoginPageContext struct {
	Tagline template.HTML
	QueryArgs map[string]string
}

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

func GetRandomTagline() string {
	n := rand.Intn(len(taglines))
	return taglines[n]
}

func (u UrlMap) Reverse(name string) string {
	result := u[name]
	if result == "" {
		panic(fmt.Sprintf("No URL named '%s'!", name))
	}
	return result
}

func RenderLoginPage(w http.ResponseWriter, context *LoginPageContext) {
	data, err := Asset("data/login.html")
	if err != nil {
		panic("Couldn't find login.html!")
	}
	s := string(data)
	t, _ := template.New("login.html").Funcs(template.FuncMap{
		"reverse": urls.Reverse,
	}).Parse(s)
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, context)
}

func RedirectToCallback(w http.ResponseWriter, u url.URL, code string, state string) {
	q := u.Query()
	q.Set("code", code)
	q.Set("state", state)
	u.RawQuery = q.Encode()
	w.Header().Set("Location", u.String())
	w.WriteHeader(302)
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
	if r.URL.Path == urls.Reverse("authorize") && r.Method == "GET" {
		rq := r.URL.Query()
		email := rq.Get("email")
		// TODO: Ensure 'client_id' is in GET params.
		// TODO: Ensure 'state' is in GET params.
		// TODO: Ensure 'response_type' is 'code'.
		if len(email) == 0 {
			queryArgs := make(map[string]string)
			for k, v := range r.URL.Query() {
				queryArgs[k] = v[0]
			}
			RenderLoginPage(w, &LoginPageContext{
				Tagline: template.HTML(GetRandomTagline()),
				QueryArgs: queryArgs,
			})
		} else {
			RedirectToCallback(w, *config.CallbackUrl, email, rq.Get("state"))
		}
	} else if r.URL.Path == urls.Reverse("token") && r.Method == "POST" {
		ExchangeCodeForAccessToken(w, r)
	} else if r.URL.Path == urls.Reverse("svgLogo") {
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

func Urlify(uStr string) *url.URL {
	u, err := url.Parse(uStr)

	if err != nil {
		panic(fmt.Sprintf("'%s' is not a valid URL!", uStr))
	}

	return u
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
	fmt.Printf("OAuth2 authorize URL: http://localhost:%d%s\n", *portPtr, urls.Reverse("authorize"))
	fmt.Printf("OAuth2 token URL: http://localhost:%d%s\n", *portPtr, urls.Reverse("token"))
	fmt.Printf("\nListening on port %d.\n", *portPtr)

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), nil)
}
