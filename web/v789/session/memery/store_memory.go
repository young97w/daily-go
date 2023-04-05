package memery

import (
	"context"
	"errors"
	"geektime/web/v789/session"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	errSessionNotFound    = errors.New("session Not Found")
	errSessionKeyNotFound = errors.New("session key Not Found")
)

type Store struct {
	//加锁
	mutex sync.RWMutex
	//使用cache包
	c *cache.Cache
	//默认过期时间
	duration time.Duration
}

//type myStore interface {
//	Generate(ctx context.Context,id string) (session.Session , error)
//	Get(ctx context.Context,id string) error
//	Refresh(ctx context.Context,id string) error
//	Remove(ctx context.Context,id string) error
//}

func NewStore(d time.Duration) *Store {
	return &Store{
		c:        cache.New(d, time.Second),
		duration: d,
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sess := &memorySession{
		id:   id,
		data: make(map[string]string),
	}

	s.c.Set(id, sess, s.duration)
	return sess, nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	//从cache中拿session
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess, ok := s.c.Get(id)
	if !ok {
		return nil, errSessionNotFound
	}

	return sess.(session.Session), nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess, ok := s.c.Get(id)
	if !ok {
		return errSessionNotFound
	}

	s.c.Set(id, sess, s.duration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.c.Get(id)
	if !ok {
		return errSessionNotFound
	}

	s.c.Delete(id)
	return nil
}

type memorySession struct {
	mutex sync.RWMutex
	id    string
	data  map[string]string
}

func (m *memorySession) Get(ctx context.Context, key string) (string, error) {
	//m.mutex.Lock()
	//defer m.mutex.Unlock()
	res, ok := m.data[key]
	if !ok {
		return "", errSessionKeyNotFound
	}

	return res, nil
}

func (m *memorySession) Set(ctx context.Context, key string, val string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = val
	return nil
}

func (m *memorySession) ID() string {
	return m.id
}
