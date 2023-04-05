package session

import "geektime/web/v789"

type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

//GetSession 获取session ，先从ctx中拿
//并缓存在context中
func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	//这个sess是局部变量
	if sess, ok := ctx.UserValues[m.SessCtxKey]; ok {
		return sess.(Session), nil
	}
	//从请求里拿
	//id --> session
	id, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}

	sess, err := m.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	//如果map为空，则new一个
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	ctx.UserValues[m.SessCtxKey] = sess
	return sess, nil
}

// InitSession 初始化一个 session，并且注入到 http response 里面
func (m *Manager) InitSession(ctx *web.Context, id string) (Session, error) {
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	//注入response中
	err = m.Inject(id, ctx.Resp)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// RefreshSession 刷新 Session
func (m *Manager) RefreshSession(ctx *web.Context) (Session, error) {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	//refresh time
	err = m.Refresh(ctx.Req.Context(), sess.ID())
	if err != nil {
		return nil, err
	}

	//注入http 响应中
	err = m.Inject(sess.ID(), ctx.Resp)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// RemoveSession 删除 Session
func (m *Manager) RemoveSession(ctx *web.Context) error {
	//获取session
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	//从两个层面移除session
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}

	err = m.Propagator.Remove(ctx.Resp)
	return err
}
