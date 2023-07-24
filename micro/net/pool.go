package net

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Pool struct {
	idlesConn chan *conn
	reqQueue  []chan *conn
	maxIdles  int
	maxCnt    int
	cnt       int
	timeout   time.Duration
	factory   func() (*conn, error)
	mutex     sync.Mutex
}

func NewPool(initCnt, maxCnt, maxIdles int, timeout time.Duration, factory func() (*conn, error)) (*Pool, error) {
	if initCnt > maxCnt {
		return nil, errors.New("pool: initCnt不得大于maxCnt")
	}
	idlesConn := make(chan *conn, maxCnt)
	for i := 0; i < maxCnt; i++ {
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		idlesConn <- conn
	}
	return &Pool{
		idlesConn: idlesConn,
		maxIdles:  maxIdles,
		maxCnt:    maxCnt,
		timeout:   timeout,
		cnt:       0,
		factory:   factory,
	}, nil
}

func (p *Pool) Get(ctx context.Context) (*conn, error) {
	// 1.空闲队列有，拿走
	// 2.空闲队列无，创建
	// 3.空闲队列无，等待
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for {
		select {
		case c := <-p.idlesConn:
			//检测超时
			p.mutex.Lock()
			if c.lastTime.Add(p.timeout).Before(time.Now()) {
				p.mutex.Unlock()
				//return nil, c.Close(context.Background()) 不应该返回，应该继续
				continue
			}
			p.cnt++
			p.mutex.Unlock()
			return c, nil
		default:
			// 创建连接
			p.mutex.Lock()
			if p.cnt < p.maxIdles {
				p.cnt++
				p.mutex.Unlock()
				return p.factory()
			}
			// 等待连接，将req放进队列
			req := make(chan *conn, 1)
			p.reqQueue = append(p.reqQueue, req)
			p.mutex.Unlock()
			for {
				select {
				case c := <-req:
					return c, nil
				case <-ctx.Done():
					go func() {
						c := <-req
						p.Put(context.Background(), c)
					}()
					return nil, ctx.Err()
				}
			}
		}
	}
}

func (p *Pool) Put(ctx context.Context, conn *conn) error {
	if p.maxIdles >= len(p.idlesConn) {
		return conn.Close(context.Background())
	}

	p.mutex.Lock()
	// 检测有无等待的请求
	if len(p.reqQueue) > 0 {
		req := p.reqQueue[0]
		p.reqQueue = p.reqQueue[1:]
		p.mutex.Unlock()
		req <- conn
		return nil
	}

	defer p.mutex.Unlock()
	//if p.maxIdles >= len(p.idlesConn) {
	//	return conn.Close(context.Background())
	//}
	//p.mutex.Lock()
	//p.idlesConn <- conn
	//p.mutex.Unlock()
	//这段换成select
	select {
	case p.idlesConn <- conn:
	default:
		return conn.Close(context.Background())
		p.cnt--
	}
	return nil
}

type conn struct {
	lastTime time.Time
}

func (c *conn) Close(ctx context.Context) error {
	return nil
}
