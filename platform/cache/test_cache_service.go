package cache

import (
	"fmt"

	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// cache format -
// key: username\npassword\nbase
// val: logged-in colly.Collector

// TestCache is a cache meant to be used
// during testing.
type TestCache struct {
}

func (TestCache) GetOrLogin(key string) (*colly.Collector, error) {
	// Confirm the credentials match the fake ones.
	fakeCredentials := fmt.Sprintf("%s\n%s\n%s", repository.FakeUsername, repository.FakePassword, repository.FakeBase)
	if key == fakeCredentials {
		return colly.NewCollector(), nil
	}
	return nil, repository.ErrorInvalidAuthentication
}

// NewTestCache makes a new Test Cache.
func NewTestCache() TestCache {
	return TestCache{}
}
