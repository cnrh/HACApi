package cache

import (
	"strings"
	"time"

	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
	"github.com/jellydator/ttlcache/v3"
)

// cache format -
// key: username\npassword\nbase
// val: logged-in colly.Collector
var cache *ttlcache.Cache[string, *colly.Collector] = nil

// InitializeCache initializes the local cache which stores
// logged-in collectors for username/password combinations.
func InitializeCache() {
	//Loader to recache username/password combos if they expired and were requested again
	loader := ttlcache.LoaderFunc[string, *colly.Collector](
		func(cache *ttlcache.Cache[string, *colly.Collector], key string) *ttlcache.Item[string, *colly.Collector] {
			//Get username/password
			splitKey := strings.Split(key, "\n")
			username, password, base := splitKey[0], splitKey[1], splitKey[2]

			//Login
			collector, err := utils.Login(base, username, password)

			if err != nil {
				return nil
			}

			item := cache.Set(key, collector, 10*time.Minute)

			return item
		},
	)

	cache = ttlcache.New(
		ttlcache.WithTTL[string, *colly.Collector](10*time.Minute),
		ttlcache.WithCapacity[string, *colly.Collector](100),
		ttlcache.WithLoader[string, *colly.Collector](loader),
	)
}

// CurrentCache returns the current cache.
func CurrentCache() *ttlcache.Cache[string, *colly.Collector] {
	//Initialize the cache if its nil
	if cache == nil {
		InitializeCache()
	}
	return cache
}
