package proxy

import (
	"fmt"

	"github.com/ONSdigital/log.go/v2/log"

	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type URL struct {
	url.URL
}

// New creates a fresh instance of httputil.ReverseProxy configured with
// the prefix-stripping director function.
func New(routerPath string, targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}
	return &httputil.ReverseProxy{Director: director(routerPath, target)}, nil
}

// director returns a proxy director function that is identical to the default one
// from httputil.NewSingleHostReverseProxy, except that it strips the path prefix
// from the forwarded URL.
// For example, if pathPrefix=/api, it forwards http://proxy-host/api/v3/chart
// to http://backend-host/v3/chart. This allows setting up multiple backends
// under a single proxy host, each backend behind a different path on the proxy.
func director(pathPrefix string, target *url.URL) func(req *http.Request) {
	singleHostProxy := httputil.NewSingleHostReverseProxy(target)
	director := func(req *http.Request) {
		originalURL := req.URL.String()

		path := strings.TrimPrefix(req.URL.Path, pathPrefix)
		rawPath := strings.TrimPrefix(req.URL.RawPath, pathPrefix)

		singleHostProxy.Director(req)

		req.URL.Path = path
		req.URL.RawPath = rawPath
		log.Info(req.Context(), fmt.Sprintf("forwarding request at: %v to: %v", originalURL, req.URL.String()))
	}
	return director
}
