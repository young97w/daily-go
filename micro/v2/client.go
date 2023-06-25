package rpc

import (
	"context"
	"errors"
	"geektime/micro/v2/message"
	"geektime/micro/v2/serialize"
	"geektime/micro/v2/serialize/json"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"time"
)

func (c *Client) InitService(service Service) error {
	// init proxy
	return setFuncField(service, c, c.serializer)
}

func setFuncField(service Service, p Proxy, s serialize.Serializer) error {
	if service == nil {
		return errors.New("rpc: 不支持nil")
	}
	val := reflect.ValueOf(service)
	typ := val.Type()
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("rpc: 只支持结构体的一级指针")
	}
	val = val.Elem()
	typ = typ.Elem()
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		if fieldVal.CanSet() {
			//create function
			fn := func(args []reflect.Value) (results []reflect.Value) {
				retVal := reflect.New(fieldTyp.Type.Out(0).Elem())
				ctx := args[0].Interface().(context.Context)
				reqData, err := s.Encode(args[1].Interface())
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				var meta map[string]string
				if isOneway(ctx) {
					meta = map[string]string{"one-way": "true"}
				}
				req := &message.Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Data:        reqData,
					Serializer:  s.Code(),
					MetaData:    meta,
				}
				//call remote function
				resp, err := p.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				err = s.Decode(resp.Data, retVal.Interface())
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				return []reflect.Value{retVal, reflect.Zero(reflect.TypeOf(new(error)).Elem())}
			}
			//use make func
			fnVal := reflect.MakeFunc(fieldTyp.Type, fn)
			fieldVal.Set(fnVal)
		}
	}
	return nil
}

type Client struct {
	pool       pool.Pool
	serializer serialize.Serializer
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	data := message.EncodeReq(req)
	// send 之后拿到响应
	resp, err := c.Send(ctx, data)
	if err != nil {
		return nil, err
	}
	return message.DecodeResp(resp), nil
}

type ClientOption func(client *Client)

func ClientWithSerializerOpt(sl serialize.Serializer) ClientOption {
	return func(c *Client) {
		c.serializer = sl
	}
}

func NewClient(addr string, opts ...ClientOption) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap:  1,
		MaxCap:      20,
		MaxIdle:     10,
		IdleTimeout: time.Minute,
		Factory: func() (interface{}, error) {
			return net.DialTimeout("tcp", addr, time.Second*3)
		},
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
		Ping: nil,
	})
	if err != nil {
		return nil, err
	}
	c := &Client{
		pool:       p,
		serializer: &json.Serializer{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (c *Client) Send(ctx context.Context, data []byte) ([]byte, error) {
	val, err := c.pool.Get()
	if err != nil {
		return nil, err
	}
	conn := val.(net.Conn)
	defer func() {
		c.pool.Put(val)
	}()
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	if isOneway(ctx) {
		return nil, errors.New("micro: 这是一个 oneway 调用，你不应该处理任何结果")
	}
	return ReadMsg(conn)
}
