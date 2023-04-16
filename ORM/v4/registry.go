package v1

import (
	"reflect"
	"sync"
)

type registry struct {
	//models map[reflect.Type]*model
	//mutex  sync.RWMutex
	models sync.Map
}

func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}

	var err error
	m, err = r.parseModel(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*model), err
}

//func (r *registry) get1(val any) (*model, error) {
//	//use double check
//	typ := reflect.TypeOf(val)
//	r.mutex.RLock()
//	m, ok := r.models[typ]
//	r.mutex.RUnlock()
//	if ok {
//		return m, nil
//	}
//
//	r.mutex.Lock()
//	defer r.mutex.Unlock()
//	m, ok = r.models[typ]
//	//再读一遍，如果ok则自己返回
//	if ok {
//		return m, nil
//	}
//	if !ok {
//		var err error
//		m, err = r.parseModel(val)
//		if err != nil {
//			return nil, err
//		}
//	}
//	r.models[typ] = m
//	return m, nil
//}
