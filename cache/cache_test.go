package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache_GetSet(t *testing.T) {
	c := NewCache(DEFAULT_CACHE_SIZE)

	c.ExpiryTime = 5

	key := []byte("gatekeeper")
	value := []byte("keymaster")

	err := c.Set(key, value)
	assert.Nil(t, err)

	result, err := c.Get(key)
	assert.Nil(t, err)

	assert.Equal(t, value, result)
}

func TestCache_Expire(t *testing.T) {
	c := NewCache(DEFAULT_CACHE_SIZE)

	c.ExpiryTime = 1

	key := []byte("gatekeeper")
	value := []byte("keymaster")

	err := c.Set(key, value)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	_, err = c.Get(key)
	assert.NotNil(t, err)
}
