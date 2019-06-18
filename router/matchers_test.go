package router

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHostMatcher_Matches(t *testing.T) {
	matcher, e := NewHostMatcher("h([a-z]+)e.com")
	request := httptest.NewRequest("GET", "/horse", nil)
	request.Host = "horse.com"

	assert.NoError(t, e)

	assert.True(t, matcher.Matches(request))
}

func TestHostMatcher_DoesntMatch(t *testing.T) {
	matcher, e := NewHostMatcher("h([a-z])e.com")
	request := httptest.NewRequest("GET", "/horse", nil)
	request.Host = "horse.com"

	assert.NoError(t, e)

	assert.False(t, matcher.Matches(request))
}

func TestHeaderMatcher_Matches(t *testing.T) {
	matcher, e := NewHeaderMatcher("X-Args", "aa")
	request := httptest.NewRequest("GET", "/horse", nil)
	request.Header.Set("X-Args", "valueaa")

	assert.NoError(t, e)

	assert.True(t, matcher.Matches(request))
}

func TestHeaderMatcher_MatchesMulti(t *testing.T) {
	matcher, e := NewHeaderMatcher("X-Args", "aa")
	request := httptest.NewRequest("GET", "/horse", nil)
	request.Header.Set("X-Args", "nopenopenope;aa")

	assert.NoError(t, e)

	assert.True(t, matcher.Matches(request))
}

func TestHeaderMatcher_NoMatch(t *testing.T) {
	matcher, e := NewHeaderMatcher("X-Args", "aa")
	request := httptest.NewRequest("GET", "/horse", nil)
	request.Header.Set("X-Nopes", "nopenopenope")

	assert.NoError(t, e)

	assert.False(t, matcher.Matches(request))
}

func TestPathMatcher_Matches(t *testing.T) {
	matcher, e := NewPathMatcher("./before/.")
	request := httptest.NewRequest("GET", "/horse/before/the/cart", nil)

	assert.NoError(t, e)

	assert.True(t, matcher.Matches(request))
}

func TestPathMatcher_DoesntMatch(t *testing.T) {
	matcher, e := NewPathMatcher("./dont/.")
	request := httptest.NewRequest("GET", "/horse/before/the/cart", nil)

	assert.NoError(t, e)

	assert.False(t, matcher.Matches(request))
}