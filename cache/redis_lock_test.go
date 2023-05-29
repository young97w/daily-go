package cache

import (
	"context"
	"fmt"
	"geektime/cache/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestClient_TryLock(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) redis.Cmdable
		key      string
		wantLock Lock
		wantErr  error
	}{
		{
			name: "get lock",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(true, nil)

				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Minute).Return(res)

				return cmd
			},
			key: "key1",
			wantLock: Lock{
				key:        "key1",
				expiration: time.Minute,
			},
		},
		{
			name: "time out",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(false, context.DeadlineExceeded)

				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Minute).Return(res)

				return cmd
			},
			key:     "key1",
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "fail to get lock",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(false, nil)

				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Minute).Return(res)

				return cmd
			},
			key:     "key1",
			wantErr: ErrFailedToPreemptLock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewClient(tc.mock(ctrl))

			l, err := c.TryLock(context.Background(), tc.key, time.Minute)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLock.key, l.key)
			assert.Equal(t, tc.wantLock.expiration, l.expiration)
			// 赋予值了
			assert.NotEmpty(t, l.value)
		})
	}
}

func TestLock_Unlock(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		key     string
		value   string
		wantErr error
	}{
		{
			name: "unlock",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(1))

				cmd.EXPECT().Eval(context.Background(), luaUnlock, []string{"key1"}, []any{"value1"}).Return(res)

				return cmd
			},
			key:   "key1",
			value: "value1",
		},
		{
			name: "not hold",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))

				cmd.EXPECT().Eval(context.Background(), luaUnlock, []string{"key1"}, []any{"value1"}).Return(res)

				return cmd
			},
			key:     "key1",
			value:   "value1",
			wantErr: ErrLockNotHold,
		},
		{
			name: "unlock err",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(context.DeadlineExceeded)

				cmd.EXPECT().Eval(context.Background(), luaUnlock, []string{"key1"}, []any{"value1"}).Return(res)

				return cmd
			},
			key:     "key1",
			value:   "value1",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			lock := &Lock{
				client: tc.mock(ctrl),
				key:    tc.key,
				value:  tc.value,
			}
			err := lock.Unlock(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestLock_Refresh(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) redis.Cmdable
		key        string
		value      string
		expiration time.Duration
		wantErr    error
	}{
		{
			name: "refreshed",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(1))

				cmd.EXPECT().Eval(context.Background(), luaRefresh, []string{"key1"}, []any{"value1", float64(60)}).Return(res)

				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Minute,
		},
		{
			name: "not hold",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))

				cmd.EXPECT().Eval(context.Background(), luaRefresh, []string{"key1"}, []any{"value1", float64(60)}).Return(res)

				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Minute,
			wantErr:    ErrLockNotHold,
		},
		{
			name: "unlock err",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(context.DeadlineExceeded)

				cmd.EXPECT().Eval(context.Background(), luaRefresh, []string{"key1"}, []any{"value1", float64(60)}).Return(res)

				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Minute,
			wantErr:    context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			lock := &Lock{
				client:     tc.mock(ctrl),
				key:        tc.key,
				value:      tc.value,
				expiration: tc.expiration,
			}
			err := lock.Refresh(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func ExampleLock_Refresh() {
	//跟业务代码一起配合，使用refresh
	//1.注意refresh时候出现错误，如有，则中断业务代码
	//2.反之，业务代码出错时候，则终止refresh
	//3.执行完毕后，关闭refresh的goroutine
	//当刷新出错时，发送err，如果时超时可以不发送，其它错误发送
	errChan := make(chan error)
	//超时的channel，给一个buffer，timeout配合超时重试机制
	timeOutChan := make(chan struct{}, 1)
	//业务结束后通知
	stopChan := make(chan struct{})

	var l *Lock
	//开一个刷新的goroutine
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		//重试次数
		timeoutRetry := 0
		for {
			select {
			case <-ticker.C:
				//refresh
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				err := l.Refresh(ctx)
				cancel()
				if err == context.DeadlineExceeded {
					timeOutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
				}
				timeoutRetry = 0
			case <-timeOutChan:
				//time out , retry
				timeoutRetry++
				if timeoutRetry > 20 {
					errChan <- context.DeadlineExceeded
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				err := l.Refresh(ctx)
				cancel()
				if err == context.DeadlineExceeded {
					timeOutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
				}
			case <-stopChan:
				return
			}
		}
	}()

	//biz1
	for i := 0; i < 100; i++ {
		select {
		case <-errChan:
			l.Unlock(context.Background())
			break
		default:
			//biz
		}
	}

	//biz2
	select {
	case <-errChan:
		l.Unlock(context.Background())
		break
	default:
		//biz
	}

	//finish
	stopChan <- struct{}{}
	err := l.Unlock(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Hi")
	// Output:
	// Hi
}
