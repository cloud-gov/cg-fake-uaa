package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
)

const CLIENT_ID_WARNING = `The <code>client_id</code> querystring parameter is missing.`
const RESPONSE_TYPE_WARNING = `The <code>response_type</code> querystring parameter is expected to be <code>code</code>.`
const STATE_WARNING = `The <code>state</code> querystring parameter is missing. See <a href="http://www.twobotechnologies.com/blog/2014/02/importance-of-state-in-oauth2.html">The Importance of the state parameter in OAuth2</a> for more details.`

var taglines = [...]string{
	`Welcome to your fake <abbr title="User Account and Authentication">UAA</abbr> provider!`,
	`The convenience of zero-factor authentication is here.`,
	`We're like <a href="https://cloud.gov">cloud.gov</a>, but without the security.`,
}

type loginPageContext struct {
	Tagline   template.HTML
	QueryArgs map[string]string
	Version   string
	Warnings  []template.HTML
}

func getRandomTagline() string {
	n := rand.Intn(len(taglines))
	return taglines[n]
}

func renderLoginPage(w http.ResponseWriter, context *loginPageContext) {
	data := GetAsset("data/login.html")
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
	var warnings []template.HTML

	rq := r.URL.Query()
	email := rq.Get("email")

	addWarning := func(warning string) {
		warnings = append(warnings, template.HTML(warning))
	}

	if rq.Get("client_id") == "" {
		addWarning(CLIENT_ID_WARNING)
	}
	if rq.Get("response_type") != "code" {
		addWarning(RESPONSE_TYPE_WARNING)
	}
	if rq.Get("state") == "" {
		addWarning(STATE_WARNING)
	}

	if email == "" {
		queryArgs := make(map[string]string)
		for k, v := range r.URL.Query() {
			queryArgs[k] = v[0]
		}
		renderLoginPage(w, &loginPageContext{
			Tagline:   template.HTML(getRandomTagline()),
			Warnings:  warnings,
			Version:   GetVersion(),
			QueryArgs: queryArgs,
		})
	} else {
		redirectToCallback(w, *config.CallbackUrl, email, rq.Get("state"))
	}
}
