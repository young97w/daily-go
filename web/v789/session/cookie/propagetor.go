package cookie

import (
	"net/http"
)

type PropagatorOption func(propagator *Propagator)

type Propagator struct {
	cookieName string
	cookieOpt  func(c *http.Cookie)
}

func NewPropagator(cookieName string, opt PropagatorOption) *Propagator {
	res := &Propagator{
		cookieName: cookieName,
	}
	if opt != nil {
		opt(res)
	}
	return res
}

func WithCookieOption(opt func(c *http.Cookie)) PropagatorOption {
	return func(propagator *Propagator) {
		propagator.cookieOpt = opt
	}
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	//将ID注入响应中
	cookie := &http.Cookie{
		Name:     p.cookieName,
		Value:    id,
		HttpOnly: false,
	}
	if p.cookieOpt != nil {
		p.cookieOpt(cookie)
	}
	http.SetCookie(writer, cookie)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	id := cookie.Value
	return id, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	//实际就是设置过期时间
	cookie := &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	}
	http.SetCookie(writer, cookie)
	return nil
}
