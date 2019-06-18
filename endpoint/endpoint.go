package endpoint

import "net/http"

const (
	HEADER_FORWARDED_HOST = "X-Forwarded-Host"
	HEADER_ORIGIN_HOST    = "X-Origin-Host"
)

type Endpoint struct {
	Name string
	Host string
	Path string
	Scheme string
}

func NewEndpoint(name string, host string, path string) *Endpoint {
	return &Endpoint{Name: name, Host: host, Path: path}
}

func (e *Endpoint) TransformRequest(request *http.Request) {
	e.setHeaders(request)

	request.Host = e.Host
	request.URL.Path = e.Path
	request.URL.Host = e.Host
	request.URL.Scheme = e.Scheme
}

func (e *Endpoint) setHeaders(request *http.Request) {
	request.Header.Set(HEADER_FORWARDED_HOST, request.Host)
	request.Header.Set(HEADER_ORIGIN_HOST, e.Host)
}
