package reverseProxy

import (
	"github.com/Banyango/Alligator/cache"
	"github.com/Banyango/Alligator/config"
	"github.com/Banyango/Alligator/endpoint"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
)

type Route struct {
	Matchers []Matcher
	Upstream endpoint.Endpoint
}

type ReverseProxyServer struct {
	Routes []*Route
	Cache *cache.Cache
}

func New(config config.Config) *ReverseProxyServer {
	server := ReverseProxyServer{
		Cache: cache.NewCache(config.CacheSize),
	}
	server.buildRoutes(config)
	return &server
}

func (r *ReverseProxyServer) Build() http.Handler {
	return r.cacheMiddleware(r.buildProxy())
}

func (r *ReverseProxyServer) cacheMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if bytes, err := r.Cache.Get([]byte(request.URL.String())); err == nil {
			responseCache, err := cache.ResponseCacheFromBytes(bytes)
			if err != nil {
				log.Println(err)
			}
			for k, v := range responseCache.HeaderMap {
				w.Header().Set(k, strings.Join(v, ","))
			}
			w.WriteHeader(responseCache.Code)
			_, err = w.Write(responseCache.Body)
			if err != nil {
				log.Println(err)
			}
			log.Println("Writing from cache: ", request)
			return
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		if recorder.Code < 400 {
			key := []byte(request.URL.String())

			response := cache.ResponseCache{
				Body:recorder.Body.Bytes(),
				HeaderMap:recorder.Result().Header,
				Code:recorder.Code,
			}

			if bytes, err := cache.ResponseCacheToBytes(&response); err == nil {
				if err := r.Cache.Set(key,bytes); err != nil {
					log.Println(err)
				}
			} else {
				log.Println(err)
			}
		}

		for k, v := range recorder.Result().Header {
			w.Header().Set(k, strings.Join(v, ","))
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
		return
	})
}

func (r *ReverseProxyServer) buildRoutes(config config.Config) {

	log.Println("Building reverse proxies")
	for _,configRoute := range config.Proxy {
		var matchers []Matcher

		log.Println("Reroute => ", configRoute.Scheme, configRoute.Host, configRoute.Path)
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
				log.Println("Added matcher ", m.String())
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

func (r *ReverseProxyServer) buildProxy() http.Handler {
	return &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			route := r.findRoute(request)
			if route != nil {
				route.Upstream.TransformRequest(request)
			}
		},
	}
}

func (r *ReverseProxyServer) findRoute(request *http.Request) *Route {
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
