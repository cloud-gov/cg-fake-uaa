package main

import (
  "testing"
  "net/url"
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

func TestRenderIndexWorksWithNoQueryArgs(t *testing.T) {
  context := new(IndexPageContext)
  recorder := httptest.NewRecorder()
  RenderIndex(recorder, context)
  assertStatus(t, recorder, 200)
  assertHeader(t, recorder, "Content-Type", "text/html")
}

func TestRenderIndexWorksWithQueryArgs(t *testing.T) {
  context := &IndexPageContext{QueryArgs: make(map[string]string)}
  context.QueryArgs["boop"] = "hi"

  recorder := httptest.NewRecorder()
  RenderIndex(recorder, context)
  assertStatus(t, recorder, 200)
  assertHeader(t, recorder, "Content-Type", "text/html")
}

func TestRedirectToCallbackWorks(t *testing.T) {
  recorder := httptest.NewRecorder()
  u, _ := url.Parse("http://example.org")
  RedirectToCallback(recorder, *u, "someCode", "someState")
  assertStatus(t, recorder, 302)
  assertHeader(t, recorder, "Location",
               "http://example.org?code=someCode&state=someState")
}
