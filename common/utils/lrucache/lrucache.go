package lrucache

import (
	"fmt"

	setting "../../../settings"
	freecache "github.com/coocood/freecache"
)

const CacheMemSize = 10 * 1024 * 1024

type CacheManager struct {
	cache *freecache.Cache
}

var v *CacheManager

//GetCacheManager single pattern lazy mod
func GetCacheManager() *CacheManager {
	if v == nil {
		cache := freecache.NewCache(CacheMemSize)
		v = &CacheManager{
			cache: cache,
		}
	}
	return v
}

func (c *CacheManager) Info() string {
	return fmt.Sprintf(`HitCount:%d,MissCount:%d`, c.cache.HitCount(), c.cache.MissCount())
}

func (c *CacheManager) MissCount() int64 {
	return c.cache.MissCount()
}
func (c *CacheManager) HitCount() int64 {
	return c.cache.HitCount()
}

//Cacheable interface for generating key for object catch
type Cacheable interface {
	CacheKey() []byte
	JSON() ([]byte, error)
	LoadJSON([]byte) error
	CacheExpireTime() int
}

func (c *CacheManager) Set(cacheObj Cacheable) error {
	key := cacheObj.CacheKey()
	data, err := cacheObj.JSON()
	if err != nil {
		return err
	}
	c.cache.Set(key, data, cacheObj.CacheExpireTime())
	return nil
}

func (c *CacheManager) Get(cacheObj Cacheable) (found bool) {
	key := cacheObj.CacheKey()
	data, err := c.cache.Get(key)
	if err != nil {
		return false
	}
	err = cacheObj.LoadJSON(data)
	if err != nil {
		setting.GetLogger(nil).Println("gocache:Get", err)
		return false
	}
	return true
}
