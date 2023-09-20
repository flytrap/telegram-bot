package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config redis配置参数
type Config struct {
	Addr      string // 地址(IP:Port)
	DB        int    // 数据库
	Password  string // 密码
	KeyPrefix string // 存储key的前缀
}

// NewStore 创建基于redis存储实例
func NewStore(cfg *Config) *Store {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	return &Store{
		cli:    cli,
		prefix: cfg.KeyPrefix,
	}
}

// NewStoreWithClient 使用redis客户端创建存储实例
func NewStoreWithClient(cli *redis.Client, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

// NewStoreWithClusterClient 使用redis集群客户端创建存储实例
func NewStoreWithClusterClient(cli *redis.ClusterClient, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

type redisClient interface {
	redis.Cmdable
	Close() error
}

// Store redis存储
type Store struct {
	cli    redisClient
	prefix string
}

func (s *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s:%s", s.prefix, key)
}

func (s *Store) Get(ctx context.Context, key string) interface{} {
	result, err := s.cli.Get(ctx, s.wrapperKey(key)).Result()
	if err != nil {
		return ""
	}
	return result
}

func (s *Store) Set(ctx context.Context, key string, v interface{}, expiration time.Duration) error {
	err := s.cli.Set(ctx, s.wrapperKey(key), v, expiration).Err()
	return err
}

func (s *Store) IsExist(ctx context.Context, key string) bool {
	n, err := s.cli.Exists(ctx, s.wrapperKey(key)).Result()
	return err == nil && n > 0
}

func (s *Store) Delete(ctx context.Context, key string) error {
	cmd := s.cli.Del(ctx, s.wrapperKey(key))
	return cmd.Err()
}

func (s *Store) Check(ctx context.Context, key string) (bool, error) {
	n, err := s.cli.Exists(ctx, s.wrapperKey(key)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (s *Store) ZAdd(ctx context.Context, key string, members ...redis.Z) (bool, error) {
	cmd := s.cli.ZAdd(ctx, s.wrapperKey(key), members...)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) ZIncrBy(ctx context.Context, key string, increment float64, member string) (bool, error) {
	cmd := s.cli.ZIncrBy(ctx, s.wrapperKey(key), increment, member)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) ZRem(ctx context.Context, key string, members ...interface{}) (bool, error) {
	cmd := s.cli.ZRem(ctx, s.wrapperKey(key), members...)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}
func (s *Store) ZRemRangeByScore(ctx context.Context, key, min, max string) (bool, error) {
	cmd := s.cli.ZRemRangeByScore(ctx, s.wrapperKey(key), min, max)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return s.cli.ZRange(ctx, s.wrapperKey(key), start, stop).Result()
}

func (s *Store) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return s.cli.ZRevRange(ctx, s.wrapperKey(key), start, stop).Result()
}

func (s *Store) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return s.cli.ZRangeByScore(ctx, s.wrapperKey(key), opt).Result()
}

func (s *Store) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	return s.cli.ZCount(ctx, s.wrapperKey(key), min, max).Result()
}

func (s *Store) Xadd(ctx context.Context, key string, value interface{}) error {
	return s.cli.XAdd(ctx, &redis.XAddArgs{Stream: s.wrapperKey(key), MaxLen: 3000, Values: value}).Err()
}

func (s *Store) Close() error {
	return s.cli.Close()
}
