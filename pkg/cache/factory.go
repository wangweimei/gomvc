package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	redigo "github.com/gomodule/redigo/redis"
)

type Factory struct {
	Config map[string]interface{}
}

func (f *Factory) ConnRedis() *redigo.Pool {
	redis := Redis{config: f.Config}
	return redis.Conn()
}

func (f *Factory) ConnMemcache() *memcache.Client {
	memcache := Memcache{config: f.Config}
	return memcache.Conn()
}
