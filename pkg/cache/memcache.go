package cache

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcache struct {
	config map[string]interface{}
}

func (m *Memcache) Conn() *memcache.Client {
	c := m.config
	dsn := fmt.Sprintf("%s:%s", c["host"].(string), c["port"].(string))
	return memcache.New(dsn)
}
