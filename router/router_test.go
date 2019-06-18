package router

import (
	"github.com/Banyango/Alligator/config"
	"github.com/Banyango/Alligator/endpoint"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

func TestRouter_BuildRoutes(t *testing.T) {

	mockConfig := config.Config{Proxy: []config.Proxy{
		{Name: "Test", Host: "www.google.com", Scheme: "http", Path: "/example", Rules: []config.Rule{{Type: "path", Pattern: []string{"home"}}}},
		{Name: "Test", Host: "www.reddit.com", Scheme: "http", Path: "/", Rules: []config.Rule{{Type: "host", Pattern: []string{"www.example.com"}}}},
	}}

	router := Router{}

	router.BuildRoutes(mockConfig)

	assert.Equal(t, len(router.Routes), 2)
	assert.Equal(t, router.Routes[0].Upstream.Path, "/example")
	assert.Equal(t, router.Routes[0].Upstream.Host, "www.google.com")
	assert.IsType(t, router.Routes[0].Matchers[0], &PathMatcher{})
	assert.IsType(t, router.Routes[1].Matchers[0], &HostMatcher{})

}

func TestRouter_BuildProxy(t *testing.T) {

	var (
		mu      sync.Mutex
		header1 string
		header2 string
	)

	backendServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		header1 = r.Header.Get(endpoint.HEADER_FORWARDED_HOST)
		header2 = r.Header.Get(endpoint.HEADER_ORIGIN_HOST)
		w.WriteHeader(http.StatusOK)
	}))
	l, _ := net.Listen("tcp", "localhost:5185")
	backendServer.Listener = l
	backendServer.Start()
	defer backendServer.Close()

	backendUrl, e := url.Parse(backendServer.URL)
	if e != nil {
		t.Fatal(e)
	}

	mockConfig := config.Config{Proxy: []config.Proxy{
		{Name: "Test", Host: backendUrl.Host, Scheme: "http", Path: "", Rules: []config.Rule{{Type: "host", Pattern: []string{"[localhost:5186]"}}}},
	}}

	router := Router{}
	router.BuildRoutes(mockConfig)

	proxy := router.BuildProxy()

	server := httptest.NewUnstartedServer(proxy)
	l1, _ := net.Listen("tcp", "localhost:5186")
	server.Listener = l1
	server.Start()
	defer server.Close()

	_, err := http.Get("http://localhost:5186")
	assert.Nil(t, err)

	mu.Lock()
	h1 := header1
	h2 := header2
	mu.Unlock()

	assert.Equal(t, "localhost:5186", h1)
	assert.Equal(t, "127.0.0.1:5185", h2)

}
