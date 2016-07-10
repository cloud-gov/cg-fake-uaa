package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder, code int) {
	if recorder.Code != code {
		t.Errorf("Expected code %d, got %d", code, recorder.Code)
	}
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder, header string, value string) {
	actualValue := recorder.HeaderMap.Get(header)
	if actualValue != value {
		t.Errorf("Expected header '%s' to be '%s', but it is '%s'",
			header, value, actualValue)
	}
}

func assertBody(t *testing.T, recorder *httptest.ResponseRecorder, value string) {
	actualValue := recorder.Body.String()
	if (actualValue != value) {
		t.Errorf("Expected body '%s', got '%s'", value, actualValue)
	}
}

func handle(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	handler := NewServerHandler(&ServerConfig{
		CallbackUrl: Urlify("http://client/callback"),
	})
	handler(recorder, request)

	return recorder
}

func TestLoginPageWorksWithoutQueryArgs(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "GET",
		URL:    Urlify("/oauth/authorize"),
	})

	assertStatus(t, recorder, 200)
	assertHeader(t, recorder, "Content-Type", "text/html")
}

func TestLoginPageWorksWithQueryArgs(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "GET",
		URL:    Urlify("/oauth/authorize?state=blah"),
	})

	assertStatus(t, recorder, 200)
	assertHeader(t, recorder, "Content-Type", "text/html")
}

func TestRedirectToCallbackWorks(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "GET",
		URL:    Urlify("/oauth/authorize?email=foo&state=bar"),
	})

	assertStatus(t, recorder, 302)
	assertHeader(t, recorder, "Location",
		"http://client/callback?code=foo&state=bar")
}

func TestTokenErrorsWhenCodeIsEmpty(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "POST",
		URL:    Urlify("/oauth/token"),
	})

	assertStatus(t, recorder, 400)
	assertHeader(t, recorder, "Content-Type", "text/plain")
	assertBody(t, recorder, "'code' is missing or empty")
}

func TestTokenWorks(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "POST",
		URL:    Urlify("/oauth/token"),
		PostForm: url.Values{
			"code": []string{"foo@bar.gov",},
		},
	})

	assertStatus(t, recorder, 200)
	assertHeader(t, recorder, "Content-Type", "application/json")
}

func TestGetSvgLogoWorks(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "GET",
		URL:    Urlify("/fake-cloud.gov.svg"),
	})

	assertStatus(t, recorder, 200)
	assertHeader(t, recorder, "Content-Type", "image/svg+xml")
}

func Test404Works(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "GET",
		URL:    Urlify("/blah"),
	})

	assertStatus(t, recorder, 404)
	assertHeader(t, recorder, "Content-Type", "text/plain")
}
