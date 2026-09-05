package gateway

import (
	"net/http/httputil"
	"net/url"
)

// NewMovieProxy creates a reverse proxy for the Movie service
func NewMovieProxy(targetURL string) (*httputil.ReverseProxy, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	return proxy, nil
}

// NewUserProxy creates a reverse proxy for the User service
func NewUserProxy(targetURL string) (*httputil.ReverseProxy, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	return proxy, nil
}
