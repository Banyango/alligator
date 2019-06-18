package router

import (
	"github.com/Banyango/Alligator/config"
	"net/http"
	"net/http/httputil"
)

type Route struct {
	Matchers []Matcher
	Route httputil.ReverseProxy
}

type Router struct {
	Routes []*Route
}




