package main

import (
  "fmt"
  "net/http"
  "net/url"
  "html/template"
)

type IndexPageContext struct {
  QueryArgs map[string]string
}

// TODO: For more docs, see:
//   * https://golang.org/doc/articles/wiki/
//   * https://golang.org/pkg/net/url/
//   * https://golang.org/pkg/net/http/

// TODO: run this through https://golang.org/cmd/gofmt/

func handler(w http.ResponseWriter, r *http.Request) {
  var data []byte
  var err error

  written := false

  if (r.URL.Path == "/oauth/authorize") {
    rq := r.URL.Query()
    email := rq.Get("email")
    if len(email) == 0 {
      data, err = Asset("data/index.html")
      if err != nil {
        panic("Couldn't find index.html!")
      }
      s := string(data)
      t, _ := template.New("index.html").Parse(s)
      w.Header().Set("Content-Type", "text/html")
      queryArgs := make(map[string]string)
      for k, v := range r.URL.Query() {
        queryArgs[k] = v[0]
      }
      t.Execute(w, &IndexPageContext{QueryArgs: queryArgs})
    } else {
      // TODO: Read callback URL from environment or cmdline
      callbackUrl := "http://localhost:8000/callback"

      u, err := url.Parse(callbackUrl)
      if err != nil {
        panic("Couldn't parse callback URL!")
      }
      q := u.Query()
      q.Set("code", email)
      q.Set("state", rq.Get("state"))
      u.RawQuery = q.Encode()
      w.Header().Set("Location", u.String())
      w.WriteHeader(302)
    }
    written = true
  } else if (r.URL.Path == "/oauth/token") {
    // TODO: Finish implementing this based on
    // https://github.com/18F/calc/blob/develop/fake_uaa_provider/views.py
    data = []byte("TODO: Implement this!")
    w.Header().Set("Content-Type", "text/plain")
  } else if (r.URL.Path == "/fake-cloud.gov.svg") {
    data, err = Asset("data/fake-cloud.gov.svg")
    if err != nil {
      panic("Couldn't find fake-cloud.gov.svg!")
    }
    w.Header().Set("Content-Type", "image/svg+xml")
  } else {
    data = []byte("Not Found")
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(404)
  }

  if !written {
    fmt.Fprintf(w, "%s", data)
  }
}

func main() {
  print("Listening on port 8080.\n")
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
