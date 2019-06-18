package endpoint

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestEndpoint_TransformRequest(t *testing.T) {
	newHost := "www.google.com"
	newPath := "/images"
	endpoint := NewEndpoint("Test", newHost, newPath)
	request := httptest.NewRequest("GET", "http://www.facebook.com", nil)

	endpoint.TransformRequest(request)

	assert.Equal(t, request.Host, newHost, "Host was not updated")
	assert.Equal(t, request.URL.Host, newHost, "Host was not updated")
	assert.Equal(t, request.URL.Path, newPath, "Path was not updated")
}

func TestEndpoint_setHeaders(t *testing.T) {
	newHost := "www.google.com"
	newPath := "/images"
	endpoint := NewEndpoint("Test", newHost, newPath)
	request := httptest.NewRequest("GET", "http://www.facebook.com", nil)

	endpoint.TransformRequest(request)

	assert.Equal(t, len(request.Header), 2, "Headers were not updated")
	assert.Equal(t, request.Header.Get(HEADER_FORWARDED_HOST), "www.facebook.com", "Forwared host wrong")
	assert.Equal(t, request.Header.Get(HEADER_ORIGIN_HOST), "www.google.com", "Origin host wrong")
}