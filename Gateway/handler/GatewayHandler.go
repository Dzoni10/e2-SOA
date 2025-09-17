package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type GatewayHandler struct{}

func proxyRequest(w http.ResponseWriter, r *http.Request, prefix string, target string) {
	//path := strings.TrimPrefix(r.URL.Path, prefix+"/")
	targetURL, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		// req.URL.Path = prefix + "/" + path
	}
	proxy.ServeHTTP(w, r)
}

func (h *GatewayHandler) HandleBlog(w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8082",
	})
	proxy.ServeHTTP(w, r)
}

func (h *GatewayHandler) HandleStakeholders(w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	})
	proxy.ServeHTTP(w, r)
}
func (h *GatewayHandler) HandleFollowers(w http.ResponseWriter, r *http.Request) {
	proxyRequest(w, r, "/followers", "http://localhost:8083")
}
func (h *GatewayHandler) HandleTours(w http.ResponseWriter, r *http.Request) {
	proxyRequest(w, r, "/position", "http://localhost:8081")
	// proxy := httputil.NewSingleHostReverseProxy(&url.URL{
	// 	Scheme: "http",
	// 	Host:   "localhost:8081",
	// })
	// proxy.ServeHTTP(w, r)
}
func (h *GatewayHandler) HandlePosition(w http.ResponseWriter, r *http.Request) {
	proxyRequest(w, r, "/position", "http://localhost:8081")
}
func (h *GatewayHandler) HandleReviews(w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8081",
	})
	proxy.ServeHTTP(w, r)
}
func (h *GatewayHandler) HandleKeypoints(w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	})
	proxy.ServeHTTP(w, r)
}

func (h *GatewayHandler) HandleTourExecutions(w http.ResponseWriter, r *http.Request) {
	proxyRequest(w, r, "/position", "http://localhost:8081")
	// proxy := httputil.NewSingleHostReverseProxy(&url.URL{
	// 	Scheme: "http",
	// 	Host:   "localhost:8081",
	// })
	// proxy.ServeHTTP(w, r)
}
