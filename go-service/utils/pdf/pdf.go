package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v3"
)

type reportsConfig map[string]string

var (
	reports  reportsConfig
	tplCache map[string]*template.Template
	loadOnce sync.Once
)

func loadReports() error {
	var err error
	loadOnce.Do(func() {
		p := filepath.Join("resource", "templates", "reports.yaml")
		b, e := os.ReadFile(p)
		if e != nil {
			err = e
			return
		}
		var r reportsConfig
		e = yaml.Unmarshal(b, &r)
		if e != nil {
			err = e
			return
		}
		reports = r
		tplCache = make(map[string]*template.Template)
	})
	return err
}

func parseTemplate(name string) error {
	if _, ok := tplCache[name]; ok {
		return nil
	}
	htmlStr, ok := reports[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}
	t, err := template.New(name).Parse(htmlStr)
	if err != nil {
		return err
	}
	tplCache[name] = t
	return nil
}

func toMap(v interface{}) (map[string]any, error) {
	if m, ok := v.(map[string]any); ok {
		return m, nil
	}
	if b, ok := v.([]byte); ok {
		var out map[string]any
		if err := json.Unmarshal(b, &out); err != nil {
			return nil, err
		}
		return out, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func renderHTML(templateName string, data map[string]any) (string, error) {
	if err := loadReports(); err != nil {
		return "", err
	}
	if err := parseTemplate(templateName); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tplCache[templateName].Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error rendering template %q: %w", templateName, err)
	}
	return buf.String(), nil
}

func htmlToPDF(html string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	dataURL := "data:text/html;base64," + base64.StdEncoding.EncodeToString([]byte(html))

	var pdf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate(dataURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			b, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0.39).
				WithMarginBottom(0.39).
				WithMarginLeft(0.39).
				WithMarginRight(0.39).
				Do(ctx)
			if err != nil {
				return err
			}
			pdf = b
			return nil
		}),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return nil, err
	}
	return pdf, nil
}

func GeneratePDF(ctx context.Context, templateName string, data interface{}) ([]byte, error) {
	m, err := toMap(data)
	if err != nil {
		return nil, err
	}
	html, err := renderHTML(templateName, m)
	if err != nil {
		return nil, err
	}
	b, err := htmlToPDF(html)
	if err != nil {
		return nil, err
	}
	return b, nil
}
