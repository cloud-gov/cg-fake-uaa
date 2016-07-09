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

func TestRenderLoginPageWorksWithNoQueryArgs(t *testing.T) {
  context := new(LoginPageContext)
  recorder := httptest.NewRecorder()
  RenderLoginPage(recorder, context)
  assertStatus(t, recorder, 200)
  assertHeader(t, recorder, "Content-Type", "text/html")
}

func TestRenderLoginPageWorksWithQueryArgs(t *testing.T) {
  context := &LoginPageContext{QueryArgs: make(map[string]string)}
  context.QueryArgs["boop"] = "hi"

  recorder := httptest.NewRecorder()
  RenderLoginPage(recorder, context)
  assertStatus(t, recorder, 200)
  assertHeader(t, recorder, "Content-Type", "text/html")
}

func TestRedirectToCallbackWorks(t *testing.T) {
  recorder := httptest.NewRecorder()
  RedirectToCallback(recorder, *urlify("http://example.org"),
                     "someCode", "someState")
  assertStatus(t, recorder, 302)
  assertHeader(t, recorder, "Location",
               "http://example.org?code=someCode&state=someState")
}

func TestHandlerGetSvgWorks(t *testing.T) {
  request := &http.Request{
    Method: "GET",
    URL: urlify("/fake-cloud.gov.svg"),
  }
  recorder := httptest.NewRecorder()

  Handler(recorder, request)
  assertStatus(t, recorder, 200)
  assertHeader(t, recorder, "Content-Type", "image/svg+xml")
}

func TestHandler404Works(t *testing.T) {
  request := &http.Request{
    Method: "GET",
    URL: urlify("/blah"),
  }
  recorder := httptest.NewRecorder()

  Handler(recorder, request)
  assertStatus(t, recorder, 404)
  assertHeader(t, recorder, "Content-Type", "text/plain")
}
