package router

import (
	"github.com/Banyango/Alligator/config"
	"github.com/Banyango/Alligator/endpoint"
	"log"
	"net/http"
	"net/http/httputil"
)

type Route struct {
	Matchers []Matcher
	Upstream endpoint.Endpoint
}

type Router struct {
	Routes []*Route
}

func (r *Router) FindRoute(request *http.Request) *Route {
	for _, route := range r.Routes {
		matchesAll := true

		for _, matcher := range route.Matchers {
			matchesAll = matchesAll && matcher.Matches(request)
		}

		if matchesAll {
			return route
		}
	}
	return nil
}

func (r *Router) BuildRoutes(config config.Config) {
	for _,configRoute := range config.Proxy {
		var matchers []Matcher

		for _, m := range configRoute.Rules {
			var newMatcher Matcher
			var err error
			switch m.Type {
			case "host":
				newMatcher, err = NewHostMatcher(m.Pattern[0])
			case "path":
				newMatcher, err = NewPathMatcher(m.Pattern[0])
			case "header":
				if len(m.Pattern) == 2 {
					newMatcher, err = NewHeaderMatcher(m.Pattern[0], m.Pattern[1])
				} else {
					log.Println("len 2 required Pattern ignored ", m.Type, m.Pattern)
				}
			}

			if err == nil {
				matchers = append(matchers, newMatcher)
				log.Println("Added matcher ", m.Type, m.Pattern)
			} else {
				log.Fatal(err)
			}
		}

		route := Route{
			Upstream: endpoint.Endpoint{Host: configRoute.Host, Path: configRoute.Path, Name: configRoute.Name, Scheme:configRoute.Scheme},
			Matchers: matchers,
		}

		r.Routes = append(r.Routes, &route)
	}

}

func (r *Router) BuildProxy() http.Handler {
	return &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			route := r.FindRoute(request)
			if route == nil {
				//return errors.New("No Route was found")
				return
			}
			route.Upstream.TransformRequest(request)
		}, ModifyResponse: func(response *http.Response) error {
			// request := response.Request
			// cache the response here.
			return nil
		},
	}
}
