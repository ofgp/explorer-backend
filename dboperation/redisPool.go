package dboperation

import (
	"errors"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

//NewRedisPool initial a new redis pool
func NewRedisPool() *redis.Pool {
	sntnl := &sentinel.Sentinel{
		Addrs:      beego.AppConfig.Strings("redisAddr"),
		MasterName: beego.AppConfig.String("redisMaster"),
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   5,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			beego.Info("masterAddr", masterAddr)
			if err != nil {
				return nil, err
			}

			dbName, _ := beego.AppConfig.Int("redisDatabase")
			c, err := redis.Dial("tcp",
				masterAddr,
				redis.DialDatabase(dbName),
				redis.DialPassword(beego.AppConfig.String("redisAuth")),
			)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}
}
