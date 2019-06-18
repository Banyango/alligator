package reverseProxy

import (
	"net/http"
	"regexp"
)

const (
	Matcher_Path_Type = "path"
	Matcher_Host_Type = "host"
	Matcher_Header_Type = "header"
)

type Matcher interface {
	Matches(request *http.Request) bool
}

type PathMatcher struct {
	URL *regexp.Regexp
}

type HostMatcher struct {
	Host *regexp.Regexp
}

type HeaderMatcher struct {
	Header *regexp.Regexp
	Value  *regexp.Regexp
}

func NewPathMatcher(pattern string) (*PathMatcher, error) {
	matcher := new(PathMatcher)
	compile, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	matcher.URL = compile
	return matcher, nil
}

func NewHostMatcher(pattern string) (*HostMatcher, error) {
	matcher := new(HostMatcher)
	compile, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	matcher.Host = compile
	return matcher, nil
}

func NewHeaderMatcher(headerPattern string, valuePattern string) (*HeaderMatcher, error) {
	matcher := new(HeaderMatcher)

	compiledHeader, errHeader := regexp.Compile(headerPattern)
	if errHeader != nil {
		return nil, errHeader
	}
	matcher.Header = compiledHeader

	compileValue, errValue := regexp.Compile(valuePattern)
	if errValue != nil {
		return nil, errValue
	}
	matcher.Value = compileValue

	return matcher, nil
}

func (p *HostMatcher) Matches(request *http.Request) bool {
	return p.Host.MatchString(request.Host)
}

func (p *HeaderMatcher) Matches(request *http.Request) bool {
	for name, value := range request.Header {
		if p.Header.MatchString(name) {
			for _, val := range value {
				if p.Value.MatchString(val) {
					return true
				}
			}
		}
	}

	return false
}

func (p *PathMatcher) Matches(request *http.Request) bool {
	return p.URL.MatchString(request.URL.Path)
}
