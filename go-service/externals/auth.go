package externals

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"go-service/constant"
	"go-service/model"
	"go-service/utils"
	config "go-service/utils/configs"
)

func mergeCookies(a, b []*http.Cookie) []*http.Cookie {
	m := map[string]*http.Cookie{}
	for _, c := range a {
		m[c.Name] = c
	}
	for _, c := range b {
		m[c.Name] = c
	}
	out := make([]*http.Cookie, 0, len(m))
	for _, c := range m {
		out = append(out, c)
	}
	return out
}

func login(ctx context.Context) (model.Tokens, error) {
	cfg := config.Load()
	entry := cfg.API["login"]
	loginURL := utils.ExpandURL(entry.URL, map[string]string{"BACKEND_BASE_URL": constant.BackendBaseURL})

	body := map[string]string{
		"username": "admin@school-admin.com",
		"password": "3OU4zn3q6Zh9",
	}
	payload, _ := json.Marshal(body)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 15 * time.Second,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, bytes.NewReader(payload))
	if err != nil {
		return model.Tokens{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return model.Tokens{}, err
	}
	defer res.Body.Close()

	uLogin, err := url.Parse(loginURL)
	if err != nil {
		return model.Tokens{}, err
	}

	cookies := mergeCookies(res.Cookies(), jar.Cookies(uLogin))

	access := utils.CookieValue(cookies, "accessToken")
	refresh := utils.CookieValue(cookies, "refreshToken")
	csrf := utils.CookieValue(cookies, "csrfToken")

	if (access == "" || refresh == "" || csrf == "") && res.StatusCode >= 300 && res.StatusCode < 400 {
		loc := res.Header.Get("Location")
		if loc != "" {
			nextURL := loc
			if !strings.HasPrefix(loc, "http://") && !strings.HasPrefix(loc, "https://") {
				ref, _ := url.Parse(loginURL)
				nextURL = ref.ResolveReference(&url.URL{Path: loc}).String()
			}
			req2, err := http.NewRequestWithContext(ctx, http.MethodGet, nextURL, nil)
			if err != nil {
				return model.Tokens{}, err
			}
			req2.Header.Set("Accept", "application/json")
			res2, err := client.Do(req2)
			if err != nil {
				return model.Tokens{}, err
			}
			defer res2.Body.Close()
			u2, _ := url.Parse(nextURL)
			cookies = mergeCookies(mergeCookies(cookies, res2.Cookies()), jar.Cookies(u2))
			if access == "" {
				access = utils.CookieValue(cookies, "accessToken")
			}
			if refresh == "" {
				refresh = utils.CookieValue(cookies, "refreshToken")
			}
			if csrf == "" {
				csrf = utils.CookieValue(cookies, "csrfToken")
			}
		}
	}

	if access == "" || refresh == "" || csrf == "" {
		return model.Tokens{}, errors.New("login failed: tokens not set")
	}

	return model.Tokens{Access: access, Refresh: refresh, CSRF: csrf}, nil
}
