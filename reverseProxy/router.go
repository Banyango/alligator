package reverseProxy

import (
	"github.com/Banyango/Alligator/cache"
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

type ReverseProxyServer struct {
	Routes []*Route
	Cache *cache.Cache
	ErrorHandler http.Handler
}

func New(config config.Config) *ReverseProxyServer {
	server := ReverseProxyServer{
		Cache: cache.NewCache(10 * 1024),
	}
	server.buildRoutes(config)
	return &server
}

func (r *ReverseProxyServer) FindRoute(request *http.Request) *Route {
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

func (r *ReverseProxyServer) Build() http.Handler {
	return r.cacheMiddleware(r.buildProxy())
}

func (r *ReverseProxyServer) buildRoutes(config config.Config) {
	for _,configRoute := range config.Proxy {
		var matchers []Matcher

		for _, m := range configRoute.Rules {
			var newMatcher Matcher
			var err error

			switch m.Type {
			case Matcher_Host_Type:
				newMatcher, err = NewHostMatcher(m.Pattern[0])
			case Matcher_Path_Type:
				newMatcher, err = NewPathMatcher(m.Pattern[0])
			case Matcher_Header_Type:
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

func (r *ReverseProxyServer) cacheMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {

		requestBytes, err := httputil.DumpRequest(request, true)
		if err != nil {
			log.Println(err)
		} else {
			if bytes, err := r.Cache.Get(requestBytes); err == nil {
				_, err := w.Write(bytes)
				if err != nil {
					log.Println(err)
				}
				return
			}
		}

		handler.ServeHTTP(w, request)
	})
}

func (r *ReverseProxyServer) buildProxy() http.Handler {
	return &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			route := r.FindRoute(request)
			if route != nil {
				route.Upstream.TransformRequest(request)
			}
		}, ModifyResponse: func(response *http.Response) error {
			requestBytes, errRequest := httputil.DumpRequest(response.Request, true)
			if errRequest != nil {
				log.Println(errRequest)
			}

			responseBytes, errResponse := httputil.DumpResponse(response, true)
			if errResponse != nil {
				log.Println(errResponse)
			}

			if errRequest == nil && errResponse == nil {
				if err := r.Cache.Set(requestBytes,responseBytes); err != nil {
					log.Println(err)
				}
			}

			return nil
		},
	}
}
