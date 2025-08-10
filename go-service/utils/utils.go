package utils

import (
	"net/http"
	"strings"
)

func ExpandURL(u string, vars map[string]string) string {
	for k, v := range vars {
		u = strings.ReplaceAll(u, "${"+k+"}", v)
	}
	return u
}

func CookieValue(cookies []*http.Cookie, name string) string {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}
