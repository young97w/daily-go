package withRedis

import (
	"context"
	"errors"
	"fmt"
	"geektime/web/v789/session"
	"github.com/redis/go-redis/v9"
	"time"
)

var errSessionNotFound = errors.New("session not found")

type Store struct {
	prefix   string
	client   redis.Cmdable
	duration time.Duration
}

type StoreOption func(s *Store)

func NewStore(prefix string, client redis.Cmdable, opts ...StoreOption) *Store {
	res := &Store{
		prefix:   prefix,
		client:   client,
		duration: 15 * time.Minute,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (s *Store) key(id string) string {
	return fmt.Sprintf("%s_%s", s.prefix, id)
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	const lua = `
redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
return redis.call("pexpire", KEYS[1], ARGV[3])
`
	key := s.key(id)
	_, err := s.client.Eval(ctx, lua, []string{key}, id, s.duration.Milliseconds()).Result()
	if err != nil {
		return nil, err
	}
	res := &Session{
		key:    key,
		id:     id,
		client: s.client,
	}
	return res, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	key := s.key(id)
	affected, err := s.client.Expire(ctx, key, s.duration).Result()
	if err != nil {
		return err
	}
	if !affected {
		return errSessionNotFound
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	key := s.key(id)
	_, err := s.client.Del(ctx, key).Result()
	return err
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	key := s.key(id)
	i, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if i < 0 {
		return nil, errSessionNotFound
	}
	res := &Session{
		key:    key,
		id:     id,
		client: s.client,
	}
	return res, nil
}

type Session struct {
	key    string
	id     string
	client redis.Cmdable
}

func (s *Session) Get(ctx context.Context, key string) (string, error) {
	return s.client.HGet(ctx, s.key, key).Result()
}

func (s *Session) Set(ctx context.Context, key string, val string) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else
	return -1
end
`
	res, err := s.client.Eval(ctx, lua, []string{s.key}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errSessionNotFound
	}
	return nil
}

func (s *Session) ID() string {
	return s.id
}
