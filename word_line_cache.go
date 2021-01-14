package cdk

import (
	"fmt"
)

// TODO: Canvas caching!

type CWordPage []WordLine

type WordLineCacheFn = func() []WordLine

type WordPageCache interface {
	Hit(tag string, fn WordLineCacheFn) []WordLine
}

type CWordLineCache struct {
	cache map[string]CWordPage
}

func NewWordPageCache() (wpc *CWordLineCache) {
	return &CWordLineCache{
		cache: make(map[string]CWordPage),
	}
}

func (c *CWordLineCache) Clear() {
	if len(c.cache) > 0 {
		c.cache = make(map[string]CWordPage)
		TraceF("WordLineCache.Clear(): cache purged")
	}
}

func (c *CWordLineCache) Hit(tag string, fn WordLineCacheFn) (lines []WordLine) {
	if v, ok := c.cache[tag]; ok {
		// TraceF("WordLineCache.Hit(): returning cached value for \"%v\"", tag)
		return v
	}
	TraceF("WordLineCache.Hit(): caching new value for \"%v\"", tag)
	c.cache[tag] = fn()
	return c.cache[tag]
}

func MakeTag(argv ...interface{}) (tag string) {
	tag += "{"
	for i, v := range argv {
		if len(tag) > 0 {
			tag += ";"
		}
		tag += fmt.Sprintf("%d=%v", i, v)
	}
	tag += "}"
	return
}
