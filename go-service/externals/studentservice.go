package externals

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go-service/constant"
	"go-service/utils"
	config "go-service/utils/configs"
)

func GetStudentReport(ctx context.Context, id string) (interface{}, error) {
	cfg := config.Load()
	t, err := login(ctx)
	if err != nil {
		return nil, err
	}
	entry := cfg.API["getStudent"]
	url := utils.ExpandURL(entry.URL, map[string]string{"BACKEND_BASE_URL": constant.BackendBaseURL, "STUDENT_ID": id})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if t.CSRF != "" {
		req.Header.Set("x-csrf-token", t.CSRF)
	}
	var cookies []string
	if t.Access != "" {
		cookies = append(cookies, "accessToken="+t.Access)
	}
	if t.CSRF != "" {
		cookies = append(cookies, "csrfToken="+t.CSRF)
	}
	if t.Refresh != "" {
		cookies = append(cookies, "refreshToken="+t.Refresh)
	}
	if len(cookies) > 0 {
		req.Header.Set("Cookie", strings.Join(cookies, "; "))
	}
	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
