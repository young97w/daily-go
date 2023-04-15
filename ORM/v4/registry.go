package v1

import (
	"reflect"
	"sync"
)

type registry struct {
	models map[reflect.Type]*model
	mutex  sync.RWMutex
}

func (r *registry) get(val any) (*model, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	typ := reflect.TypeOf(val)

	m, ok := r.models[typ]
	if !ok {
		var err error
		m, err = r.parseModel(val)
		if err != nil {
			return nil, err
		}
	}

	r.models[typ] = m
	return m, nil
}
