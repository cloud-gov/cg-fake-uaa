package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
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
	if actualValue != value {
		t.Errorf("Expected body '%s', got '%s'", value, actualValue)
	}
}

func assertBodyMatches(t *testing.T, recorder *httptest.ResponseRecorder, restr string) {
	actualValue := recorder.Body.String()
	re := regexp.MustCompile(restr)
	if !re.MatchString(actualValue) {
		t.Errorf("Expected body '%s' to match '%s'", actualValue, restr)
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
	assertBodyMatches(t, recorder, `type="hidden" name="state" value="blah"`)
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

func assertTokenError(t *testing.T, postForm url.Values, body string) {
	recorder := handle(&http.Request{
		Method: "POST",
		URL:    Urlify("/oauth/token"),
		PostForm: postForm,
	})

	assertStatus(t, recorder, 400)
	assertHeader(t, recorder, "Content-Type", "text/plain")
	assertBody(t, recorder, body)
}

func TestTokenErrorsWhenCodeIsEmpty(t *testing.T) {
	assertTokenError(t, url.Values{}, "'code' is missing or empty")
}

func TestTokenErrorsWhenClientIdIsEmpty(t *testing.T) {
	assertTokenError(t, url.Values{
		"code": []string{"foo@bar.gov",},
	}, "'client_id' is missing or empty")
}

func TestTokenWorks(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "POST",
		URL:    Urlify("/oauth/token"),
		PostForm: url.Values{
			"code": []string{"foo@bar.gov",},
			"client_id": []string{"baz"},
		},
	})

	assertStatus(t, recorder, 200)
	assertHeader(t, recorder, "Content-Type", "application/json")

	// TODO: Examine the response, decode the access token and ensure it's what we expect.
}

func TestGetSvgLogoWorks(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "GET",
		URL:    Urlify("/fake-cloud.gov.svg"),
	})

	assertStatus(t, recorder, 200)
	assertHeader(t, recorder, "Content-Type", "image/svg+xml")
}

func TestGetStylesheetWorks(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "GET",
		URL:    Urlify("/style.css"),
	})

	assertStatus(t, recorder, 200)
	assertHeader(t, recorder, "Content-Type", "text/css")
}

func Test404Works(t *testing.T) {
	recorder := handle(&http.Request{
		Method: "GET",
		URL:    Urlify("/blah"),
	})

	assertStatus(t, recorder, 404)
	assertHeader(t, recorder, "Content-Type", "text/plain")
}
