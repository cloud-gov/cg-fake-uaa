package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
)

var taglines = [...]string{
	`Welcome to your fake <abbr title="User Account and Authentication">UAA</abbr> provider!`,
	`The convenience of zero-factor authentication is here.`,
	`We're like <a href="https://cloud.gov">cloud.gov</a>, but without the security.`,
}

type loginPageContext struct {
	Tagline template.HTML
	QueryArgs map[string]string
}

func getRandomTagline() string {
	n := rand.Intn(len(taglines))
	return taglines[n]
}

func renderLoginPage(w http.ResponseWriter, context *loginPageContext) {
	data, err := GetAsset("data/login.html")
	if err != nil {
		panic("Couldn't find login.html!")
	}
	s := string(data)
	t, _ := template.New("login.html").Funcs(template.FuncMap{
		"reverse": Urls.Reverse,
	}).Parse(s)
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, context)
}

func redirectToCallback(w http.ResponseWriter, u url.URL, code string, state string) {
	q := u.Query()
	q.Set("code", code)
	q.Set("state", state)
	u.RawQuery = q.Encode()
	w.Header().Set("Location", u.String())
	w.WriteHeader(302)
}

func Authorize(config *ServerConfig, w http.ResponseWriter, r *http.Request) {
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
		renderLoginPage(w, &loginPageContext{
			Tagline: template.HTML(getRandomTagline()),
			QueryArgs: queryArgs,
		})
	} else {
		redirectToCallback(w, *config.CallbackUrl, email, rq.Get("state"))
	}
}
