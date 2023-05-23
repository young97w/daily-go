package cache

//func TestRedisCache_Get(t *testing.T) {
//	testCases := []struct {
//		name    string
//		mock    func(ctrl *gomock.Controller) redis.Cmdable
//		key     string
//		wantErr error
//		wantVal string
//	}{
//		{
//			name: "get value",
//			mock: func(ctrl *gomock.Controller) redis.Cmdable {
//				cmd := mocks.NewMockCmdable(ctrl)
//
//				str := redis.NewStringCmd(context.Background())
//				str.SetVal("value1")
//
//				cmd.EXPECT().Get(context.Background(), "key1").Return(str)
//				return cmd
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			c := NewRedisCache(tc.mock(ctrl))
//
//			val, err := c.Get(context.Background(), tc.key)
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//			assert.Equal(t, tc.wantVal, val)
//		})
//	}
//}
