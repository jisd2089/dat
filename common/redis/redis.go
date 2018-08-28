package redis

/**
    Author: luzequan
    Created: 2018-01-13 20:05:28
*/
import (
	"gopkg.in/redis.v5"
	"time"
	"sync"
	"fmt"
)

type (
	RedisClient interface {
		Get(key string) ([]byte, error)
		Set(key string, value []byte, expiration time.Duration) error
		HGet(key, field string) ([]byte, error)
		HGetAll(key string) (map[string]string, error)

		GetString(key string) (string, error)
		SetString(key string, value string, expiration time.Duration) error
		HGetString(key, field string) (string, error)
		HSetString(key, field, value string) error
		HDelString(key string, field []string) error
		HMSetStrings(key string, fields map[string]string) error
		Expire(key string, expiration time.Duration) error
		HExistString(key string) (bool, error)
		HIncrBy(key, field string, incr int64) error
		HGetInt64(key, field string) (int64, error)
		HSetIn64(key, field string, value int64) error
		Keys(pattern string) ([]string, error)
		PipeLineSetString(kvs []PipeKeyValue) error
		ConfigSet(password string) error
		Auth(password string) error
		Ping() error
	}
	redisClient struct {
		isSingle      bool                 // 是否为单节点
		client        *redis.Client        // 单节点cli
		clusterClient *redis.ClusterClient // 集群cli
	}
	PipeKeyValue struct {
		Key        string
		Value      string
		Expiration time.Duration
	}
	ConnectOptions struct {
		AddressList []string
		Password    string
		DBIndex     int
		PoolSize    int
	}
)

const (
	settings_xpath_addr     = "redis.Addr"
	settings_xpath_db       = "redis.DB"
	settings_xpath_poolsize = "redis.PoolSize"
	password_default        = "2d6hwi22oM3KUyhd"
	poolsize_default        = 16
	//settings_xpath_readonly = "redis.ReadOnly"
)

var (
	redisCli RedisClient
	once     sync.Once
	mutex    sync.Mutex
)

func NewRedisClient(opts *redis.Options) RedisClient {
	return &redisClient{
		isSingle: true,
		client:   redis.NewClient(opts),
	}
}

func NewRedisClusterClient(opts *redis.ClusterOptions) RedisClient {
	return &redisClient{
		isSingle:      false,
		clusterClient: redis.NewClusterClient(opts),
	}
}

func GetRedisClient(o *ConnectOptions) (RedisClient, error) {
	mutex.Lock()
	defer mutex.Unlock()
	opts, err := connect(o, false)
	if err != nil {
		return nil, err
	}
	redisClient := NewRedisClient(opts.(*redis.Options))
	password := o.Password
	if password == "" {
		password = password_default
	}
	if err = redisClient.ConfigSet(password); err != nil {
		opts, _ = connect(o, true)
		redisClient = NewRedisClient(opts.(*redis.Options))
	}
	if err = redisClient.Auth(password); err != nil {
		fmt.Println("redis err: ", err.Error())
		return nil, err
	}
	return redisClient, nil
}

func GetRedisClusterClient(o *ConnectOptions) (RedisClient, error) {
	mutex.Lock()
	defer mutex.Unlock()
	opts, err := connect(o, false)
	if err != nil {
		return nil, err
	}
	password := o.Password
	if password == "" {
		password = password_default
	}
	redisClient := NewRedisClusterClient(opts.(*redis.ClusterOptions))
	if err = redisClient.ConfigSet(password); err != nil {
		opts, _ = connect(o, true)
		redisClient = NewRedisClusterClient(opts.(*redis.ClusterOptions))
	}
	if err = redisClient.Auth(password); err != nil {
		return nil, err
	}
	return redisClient, nil
}

func connect(connectOpts *ConnectOptions, needAuth bool) (interface{}, error) {
	if len(connectOpts.AddressList) == 0 {
		return nil, fmt.Errorf("redis addrs is nil")
	}
	poolSize := connectOpts.PoolSize
	if poolSize == 0 {
		poolSize = poolsize_default
	}
	addr := connectOpts.AddressList
	switch len(addr) {
	case 1:
		opts := &redis.Options{
			Addr:         addr[0],
			DB:           connectOpts.DBIndex,
			PoolSize:     poolSize,
			DialTimeout:  time.Second * time.Duration(10),
			WriteTimeout: time.Second * time.Duration(10),
			ReadTimeout:  time.Second * time.Duration(10),
		}
		if needAuth && connectOpts.Password != "" {
			opts.Password = connectOpts.Password
		} else if needAuth && connectOpts.Password == "" {
			opts.Password = password_default
		}
		return opts, nil

	default:
		opts := &redis.ClusterOptions{
			Addrs:    addr,
			PoolSize: poolSize,
		}
		if needAuth && connectOpts.Password != "" {
			opts.Password = connectOpts.Password
		} else if needAuth && connectOpts.Password == "" {
			opts.Password = password_default
		}
		return opts, nil
	}
}

// Get 实现RedisClient接口的Get方法
func (rc *redisClient) Get(key string) ([]byte, error) {
	var value []byte
	var err error

	if rc.isSingle {
		value, err = rc.client.Get(key).Bytes()
	} else {
		value, err = rc.clusterClient.Get(key).Bytes()
	}
	if err == redis.Nil {
		return nil, nil
	}
	return value, err
}

// Set 实现RedisClient接口的Set方法
func (rc *redisClient) Set(key string, value []byte, expiration time.Duration) error {
	if rc.isSingle {
		return rc.client.Set(key, value, expiration).Err()
	}
	return rc.clusterClient.Set(key, value, expiration).Err()
}

func (rc *redisClient) HGet(key, field string) ([]byte, error) {
	var value []byte
	var err error

	if rc.isSingle {
		value, err = rc.client.HGet(key, field).Bytes()
	} else {
		value, err = rc.clusterClient.HGet(key, field).Bytes()
	}
	if err == redis.Nil {
		return nil, nil
	}
	return value, err
}

func (rc *redisClient) HGetAll(key string) (map[string]string, error) {
	var val map[string]string
	var err error
	if rc.isSingle {
		val, err = rc.client.HGetAll(key).Result()
	} else {
		val, err = rc.clusterClient.HGetAll(key).Result()
	}
	return val, err
}

func (rc *redisClient) GetString(key string) (string, error) {
	value, err := rc.Get(key)
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (rc *redisClient) SetString(key string, value string, expiration time.Duration) error {
	return rc.Set(key, []byte(value), expiration)
}

func (rc *redisClient) HGetString(key, field string) (string, error) {
	value, err := rc.HGet(key, field)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (rc *redisClient) HSetString(key, field, value string) error {
	var err error
	if rc.isSingle {
		_, err = rc.client.HSet(key, field, value).Result()

	} else {
		_, err = rc.clusterClient.HSet(key, field, value).Result()
	}
	if err != redis.Nil {
		return err
	}
	return nil
}

func (rc *redisClient) HDelString(key string, field []string) error {
	var err error
	if rc.isSingle {
		err = rc.client.HDel(key, field...).Err()
	} else {
		err = rc.clusterClient.HDel(key, field...).Err()
	}
	return err
}

func (rc *redisClient) HMSetStrings(key string, fields map[string]string) error {
	var err error
	if rc.isSingle {
		_, err = rc.client.HMSet(key, fields).Result()

	} else {
		_, err = rc.clusterClient.HMSet(key, fields).Result()
	}
	if err != redis.Nil {
		return err
	}
	return nil
}

func (rc *redisClient) Expire(key string, expiration time.Duration) error {
	var err error
	if rc.isSingle {
		_, err = rc.client.Expire(key, expiration).Result()

	} else {
		_, err = rc.clusterClient.Expire(key, expiration).Result()
	}
	if err != redis.Nil {
		return err
	}
	return nil
}

func (rc *redisClient) HExistString(key string) (bool, error) {
	var val map[string]string
	var err error
	if rc.isSingle {
		val, err = rc.client.HGetAll(key).Result()
	} else {
		val, err = rc.clusterClient.HGetAll(key).Result()
	}
	if err != nil {
		return false, err
	}
	if len(val) == 0 {
		return false, nil
	}
	return true, nil
}

func (rc *redisClient) HIncrBy(key, field string, incr int64) error {
	var err error
	if rc.isSingle {
		err = rc.client.HIncrBy(key, field, incr).Err()
	} else {
		err = rc.clusterClient.HIncrBy(key, field, incr).Err()
	}
	return err
}

func (rc *redisClient) HGetInt64(key, field string) (int64, error) {
	var value int64
	var err error
	if rc.isSingle {
		value, err = rc.client.HGet(key, field).Int64()
	} else {
		value, err = rc.clusterClient.HGet(key, field).Int64()
	}
	if err == redis.Nil {
		return 0, nil
	}
	return value, err
}

func (rc *redisClient) HSetIn64(key, field string, value int64) error {
	var err error
	if rc.isSingle {
		err = rc.client.HSet(key, field, value).Err()

	} else {
		err = rc.clusterClient.HSet(key, field, value).Err()
	}
	return err
}

func (rc *redisClient) Keys(pattern string) ([]string, error) {
	var err error
	var keys []string
	if rc.isSingle {
		keys, err = rc.client.Keys(pattern).Result()

	} else {
		keys, err = rc.clusterClient.Keys(pattern).Result()
	}
	if err != nil || err != redis.Nil {
		return keys, err
	}
	return keys, nil
}

func (rc *redisClient) PipeLineSetString(kvs []PipeKeyValue) error {
	if rc.isSingle {
		pipeline := rc.client.Pipeline()
		for _, kv := range kvs {
			cmd := pipeline.Set(kv.Key, kv.Value, kv.Expiration)
			if cmd.Err() != nil {
				return cmd.Err()
			}
		}
		_, err := pipeline.Exec()
		if err != nil {
			return err
		}
	} else {
		//keys, err = rc.clusterClient.Keys(pattern).Result()
	}
	return nil
}

func (rc *redisClient) ConfigSet(password string) error {
	var err error
	if rc.isSingle {
		err = rc.client.ConfigSet("requirepass", password).Err()

	} else {
		err = rc.clusterClient.ConfigSet("requirepass", password).Err()
	}
	return err
}

func (rc *redisClient) Auth(password string) error {
	var err error
	if rc.isSingle {
		_, err = rc.client.Pipelined(func(pipe *redis.Pipeline) error {
			pipe.Auth(password)
			return nil
		})
	} else {
		_, err = rc.clusterClient.Pipelined(func(pipe *redis.Pipeline) error {
			pipe.Auth(password)
			return nil
		})
	}
	return err
}

func (rc *redisClient) Ping() error {
	var err error
	if rc.isSingle {
		_, err = rc.client.Ping().Result()
	} else {
		_, err = rc.clusterClient.Ping().Result()
	}
	return err
}
