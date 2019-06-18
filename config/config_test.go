package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const tomlString = `
[[Proxy]]
Path = "/images"
Host = "www.google.com"
Scheme= "http"
	[[Proxy.Rules]]
		Type = "host"
		Pattern = ["www.facebook.com", "www.instagram.com"]
	[[Proxy.Rules]]
		Type = "path"
		Pattern = [".hello/."]
`

func TestNewConfig(t *testing.T) {
	config, err := NewConfig(tomlString)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(config.Proxy))

	proxy := config.Proxy[0]

	assert.Equal(t, "/images", proxy.Path)
	assert.Equal(t, "www.google.com", proxy.Host)
	assert.Equal(t, "http", proxy.Scheme)
	assert.Equal(t, 2, len(proxy.Rules))

	rule := proxy.Rules[0]
	assert.Equal(t, "host", rule.Type)

	assert.Equal(t, 2, len(rule.Pattern))
	assert.Equal(t, "www.facebook.com", rule.Pattern[0])
	assert.Equal(t, "www.instagram.com", rule.Pattern[1])

	rule2 := proxy.Rules[1]
	assert.Equal(t, "path", rule2.Type)

	assert.Equal(t, 1, len(rule2.Pattern))
	assert.Equal(t, ".hello/.", rule2.Pattern[0])
}
