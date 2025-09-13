package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type GatewayHandler struct{}

func proxyRequest(w http.ResponseWriter, r *http.Request, prefix string, target string) {
	path := strings.TrimPrefix(r.URL.Path, prefix+"/")
	targetURL, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		req.URL.Path = prefix + "/" + path
	}
	proxy.ServeHTTP(w, r)
}

// BLOG servis (8082)
func (h *GatewayHandler) HandleBlog(w http.ResponseWriter, r *http.Request) {
	proxyRequest(w, r, "/blogs", "http://localhost:8082")
}

// STAKEHOLDERS servis (8080)
func (h *GatewayHandler) HandleStakeholders(w http.ResponseWriter, r *http.Request) {
	proxyRequest(w, r, "/users", "http://localhost:8080")
}
func (h *GatewayHandler) HandleFollowers(w http.ResponseWriter, r *http.Request) {
	proxyRequest(w, r, "/followers", "http://localhost:8083")
}
func (h *GatewayHandler) HandleTours(w http.ResponseWriter, r *http.Request) {
	proxyRequest(w, r, "/", "http://localhost:8081")
}
