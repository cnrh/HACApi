package cache

import (
	"errors"
	"strings"
	"time"

	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
	"github.com/jellydator/ttlcache/v3"
)

// cache format -
// key: username\npassword\nbase
// val: logged-in colly.Collector
type TTLCache struct {
	Cache *ttlcache.Cache[string, *colly.Collector]
}

// NewCache creates a new TTL cache which stores
// logged-in collectors for username/password combinations.
func NewCache(scraper repository.ScraperProvider) *TTLCache {
	// Loader to recache username/password combos if they expired and were requested again
	loader := ttlcache.LoaderFunc[string, *colly.Collector](
		func(cache *ttlcache.Cache[string, *colly.Collector], key string) *ttlcache.Item[string, *colly.Collector] {
			// Get username/password
			splitKey := strings.Split(key, "\n")
			username, password, base := splitKey[0], splitKey[1], splitKey[2]

			// Login
			collector, err := scraper.Login(base, username, password)

			if err != nil {
				return nil
			}

			item := cache.Set(key, collector, 10*time.Minute)

			return item
		},
	)

	cache := ttlcache.New(
		ttlcache.WithTTL[string, *colly.Collector](10*time.Minute),
		ttlcache.WithCapacity[string, *colly.Collector](100),
		ttlcache.WithLoader[string, *colly.Collector](loader),
	)

	return &TTLCache{Cache: cache}
}

func (cache TTLCache) GetOrLogin(key string) (*colly.Collector, error) {
	res := cache.Cache.Get(key)
	if res == nil {
		return nil, errors.New("not found")
	}
	return res.Value(), nil
}
