package provider

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/cocotyty/summer"
	"gopkg.in/ini.v1"
	"gopkg.in/redis.v2"
	"strings"
	"time"
)

type Cacher struct {
	ID string
	c  *redis.Client
}

func (this *Cacher) Get(key string) string {
	r := this.c.Get(this.ID + ":" + key)
	v, _ := r.Result()
	return v
}

func (this *Cacher) Set(key string, value string) (err error) {
	r := this.c.Set(this.ID + ":" + key, value)
	return r.Err()
}

func init() {
	summer.Add("CacherFactory", &CacherFactory{})
}

type CacherFactory struct {
	RedisClient *redis.Client
	RedisConf   string `sm:"#.redis.conf"`
}

func (this *CacherFactory) InitClient() error {

	cfg, err := ini.Load([]byte(strings.Replace(this.RedisConf, ",", "\n", -1)))
	if err != nil {
		return err
	}

	opt := &redis.Options{
		Network: "tcp",
	}
	for k, v := range cfg.Section("").KeysHash() {
		switch k {
		case "network":
			opt.Network = v
		case "addr":
			opt.Addr = v
		case "password":
			opt.Password = v
		case "db":
			opt.DB = com.StrTo(v).MustInt64()
		case "pool_size":
			opt.PoolSize = com.StrTo(v).MustInt()
		case "idle_timeout":
			opt.IdleTimeout, err = time.ParseDuration(v + "s")
			if err != nil {
				return fmt.Errorf("error parsing idle timeout: %v", err)
			}
		case "hset_name":
		case "prefix":
		default:
			return fmt.Errorf("session/redis: unsupported option '%s'", k)
		}
	}

	this.RedisClient = redis.NewClient(opt)
	if err = this.RedisClient.Ping().Err(); err != nil {
		return err
	}

	return nil
}

func (this *CacherFactory) Ready() {
	this.InitClient()
}

func (this *CacherFactory) NewCacher(ID string) *Cacher {
	return &Cacher{
		ID: ID,
		c:  this.RedisClient,
	}
}
