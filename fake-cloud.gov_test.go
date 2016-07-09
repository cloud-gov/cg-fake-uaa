package main

import (
  "testing"
  "net/http/httptest"
)

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder,
                  code int) {
  if recorder.Code != code {
    t.Errorf("Expected code %d, got %d", code, recorder.Code)
  }
}

func TestRenderIndexWorksWithNoQueryArgs(t *testing.T) {
  context := new(IndexPageContext)
  recorder := httptest.NewRecorder()
  RenderIndex(recorder, context)
  assertStatus(t, recorder, 200)
}

func TestRenderIndexWorksWithQueryArgs(t *testing.T) {
  context := &IndexPageContext{QueryArgs: make(map[string]string)}
  context.QueryArgs["boop"] = "hi"

  recorder := httptest.NewRecorder()
  RenderIndex(recorder, context)
  assertStatus(t, recorder, 200)
}
