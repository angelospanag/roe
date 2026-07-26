package feed

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	faviconHTTPTimeout  = 10 * time.Second
	maxFaviconHTMLBytes = 1 << 20 // 1 MiB, enough for any reasonable <head>
)

// discoverFaviconURL fetches homePageURL and looks for a <link rel="icon">
// (or "shortcut icon") tag in its HTML, resolving it to an absolute URL.
// Falls back to the conventional {origin}/favicon.ico path if the page can't
// be fetched or declares no icon — same priority order NetNewsWire uses,
// since sites like blog.jetbrains.com don't serve /favicon.ico directly.
func discoverFaviconURL(
	ctx context.Context,
	httpClient *http.Client,
	homePageURL string,
) (string, error) {
	base, err := url.Parse(homePageURL)
	if err != nil || base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("invalid home page URL %q: %w", homePageURL, err)
	}
	defaultFaviconURL := base.Scheme + "://" + base.Host + "/favicon.ico"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	if err != nil {
		return defaultFaviconURL, nil
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return defaultFaviconURL, nil
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return defaultFaviconURL, nil
	}

	if iconURL := findFaviconLink(
		io.LimitReader(resp.Body, maxFaviconHTMLBytes),
		base,
	); iconURL != "" {
		return iconURL, nil
	}

	return defaultFaviconURL, nil
}

// findFaviconLink scans HTML tokens for the first usable <link rel="icon">
// href, resolved against base. Returns "" if none is found.
func findFaviconLink(r io.Reader, base *url.URL) string {
	tokenizer := html.NewTokenizer(r)

	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return ""
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data != "link" {
				continue
			}

			var rel, href, typ string
			for _, attr := range token.Attr {
				switch attr.Key {
				case "rel":
					rel = strings.ToLower(attr.Val)
				case "href":
					href = attr.Val
				case "type":
					typ = strings.ToLower(attr.Val)
				}
			}

			if href == "" || !strings.Contains(rel, "icon") {
				continue
			}
			// Skip SVG favicons — the frontend renders this into a plain
			// <img>, which not all browsers rasterize reliably from SVG.
			if typ == "image/svg+xml" || strings.HasSuffix(strings.ToLower(href), ".svg") {
				continue
			}

			resolved, err := base.Parse(href)
			if err != nil {
				continue
			}
			return resolved.String()
		}
	}
}
