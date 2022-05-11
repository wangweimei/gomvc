package cache

import (
	redigo "github.com/gomodule/redigo/redis"
)

type Redis struct {
	config map[string]interface{}
}

func (r *Redis) Conn() *redigo.Pool {
	c := r.config
	pool := redigo.NewPool(func() (redigo.Conn, error) {
		conn, err := redigo.Dial("tcp", c["host"].(string)+":"+c["port"].(string))
		if err != nil {
			return nil, err
		}
		if c["pwd"].(string) != "" {
			if _, err := conn.Do("AUTH", c["pwd"].(string)); err != nil {
				conn.Close()
				return nil, err
			}
		}
		return conn, nil
	}, int(c["maxIdle"].(float64)))
	return pool
}
