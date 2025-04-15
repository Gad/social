package cache

import (
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

func NewMemcachedClient(host string, startingPort int, endingPort int) *memcache.Client {

	var addr []string
	for port := range(endingPort - startingPort +1) {
		addr = append(addr, fmt.Sprintf("%s:%d", host, startingPort+port))
	}
	log.Printf("Memcached address: %v", addr)
	client := memcache.New(addr...)

	return client
}
