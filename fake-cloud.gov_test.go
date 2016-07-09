package main

import (
  "testing"
  "fmt"
  "net/url"
  "net/http"
  "net/http/httptest"
)

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder,
                  code int) {
  if recorder.Code != code {
    t.Errorf("Expected code %d, got %d", code, recorder.Code)
  }
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder,
                  header string, value string) {
  actualValue := recorder.HeaderMap.Get(header)
  if actualValue != value {
    t.Errorf("Expected header '%s' to be '%s', but it is '%s'",
             header, value, actualValue)
  }
}

func urlify(uStr string) *url.URL {
  u, err := url.Parse(uStr)

  if (err != nil) {
    panic(fmt.Sprintf("'%s' is not a valid URL!", uStr))
  }

  return u
}

func handle(request *http.Request) *httptest.ResponseRecorder {
  recorder := httptest.NewRecorder()

  Handler(recorder, request)

  return recorder
}

func TestLoginPageWorksWithoutQueryArgs(t *testing.T) {
  recorder := handle(&http.Request{
    Method: "GET",
    URL: urlify("/oauth/authorize"),
  })

  assertStatus(t, recorder, 200)
  assertHeader(t, recorder, "Content-Type", "text/html")
}

func TestLoginPageWorksWithQueryArgs(t *testing.T) {
  recorder := handle(&http.Request{
    Method: "GET",
    URL: urlify("/oauth/authorize?state=blah"),
  })

  assertStatus(t, recorder, 200)
  assertHeader(t, recorder, "Content-Type", "text/html")
}

func TestRedirectToCallbackWorks(t *testing.T) {
  recorder := handle(&http.Request{
    Method: "GET",
    URL: urlify("/oauth/authorize?email=foo&state=bar"),
  })

  assertStatus(t, recorder, 302)
  assertHeader(t, recorder, "Location",
               "http://localhost:8000/callback?code=foo&state=bar")
}

func TestGetSvgWorks(t *testing.T) {
  recorder := handle(&http.Request{
    Method: "GET",
    URL: urlify("/fake-cloud.gov.svg"),
  })

  assertStatus(t, recorder, 200)
  assertHeader(t, recorder, "Content-Type", "image/svg+xml")
}

func Test404Works(t *testing.T) {
  recorder := handle(&http.Request{
    Method: "GET",
    URL: urlify("/blah"),
  })

  assertStatus(t, recorder, 404)
  assertHeader(t, recorder, "Content-Type", "text/plain")
}
