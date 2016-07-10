package main

import (
	"fmt"
	"net/url"
)

type UrlMap map[string]string

var Urls = UrlMap{
	"authorize":  "/oauth/authorize",
	"token":      "/oauth/token",
	"svgLogo":    "/fake-cloud.gov.svg",
	"stylesheet": "/style.css",
}

func (u UrlMap) Reverse(name string) string {
	result := u[name]
	if result == "" {
		panic(fmt.Sprintf("No URL named '%s'!", name))
	}
	return result
}

func Urlify(uStr string) *url.URL {
	u, err := url.Parse(uStr)

	if err != nil {
		panic(fmt.Sprintf("'%s' is not a valid URL!", uStr))
	}

	return u
}
