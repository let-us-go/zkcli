package core

import "github.com/c-bata/go-prompt"

type SuggestCache struct {
	data map[string][]prompt.Suggest
}

func newSuggestCache() *SuggestCache {
	return &SuggestCache{
		data: map[string][]prompt.Suggest{},
	}
}

func (c *SuggestCache) get(key string) (value []prompt.Suggest, ok bool) {
	value, ok = c.data[key]
	return
}

func (c *SuggestCache) set(key string, value []prompt.Suggest) {
	c.data[key] = value
}

func (c *SuggestCache) del(key string) {
	delete(c.data, key)
}
