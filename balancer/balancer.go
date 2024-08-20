package balancer

import (
	"net/http"
	"net/url"
)

type LoadBalancer struct {
	Path       string
	serverPool *ServerPool
}

func New(path string) *LoadBalancer {
	pool := NewServerPool()
	return &LoadBalancer{
		Path:       path,
		serverPool: pool,
	}
}

func (lb *LoadBalancer) Register(backendURL *url.URL) {
	if backendURL == nil {
		return
	}

	lb.serverPool.RegisterBackend(backendURL)
}
func (lb *LoadBalancer) GetPeer() *Backend {
	// Use strategy pattern here.
	peer := lb.serverPool.NextPeer()
	return peer
}

func (lb *LoadBalancer) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peer := lb.GetPeer()
		if peer != nil {
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}
}
