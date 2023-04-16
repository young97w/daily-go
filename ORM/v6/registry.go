package v1

import (
	"geektime/ORM/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type registry struct {
	//models map[reflect.Type]*Model
	//mutex  sync.RWMutex
	models sync.Map
}

var _ Register = &registry{}

type Register interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOpt) (*Model, error)
}

func (r *registry) Register(val any, opts ...ModelOpt) (*Model, error) {
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err = opt(m)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (r *registry) parseModel(entity any) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	numField := typ.NumField()
	fields := make(map[string]*field, numField)

	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		//get column name
		m, err := r.parseTag(f.Tag)
		if err != nil {
			return nil, err
		}
		colName := m[tagKeyColumn]
		if colName == "" {
			colName = underscoreName(f.Name)
		}

		fields[f.Name] = &field{colName: colName}
	}

	//get table name
	var tblName string
	if tn, ok := entity.(TableName); ok {
		tblName = tn.TableName()
	}
	if tblName == "" {
		tblName = underscoreName(typ.Name())
	}

	return &Model{
		tableName: tblName,
		fields:    fields,
	}, nil
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	//get by "orm"
	res := make(map[string]string, 1)

	val := tag.Get("orm")
	if val == "" {
		return res, nil
	}

	pairs := strings.Split(val, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return res, errs.NewErrInvalidTagContent(pair)
		}

		res[kv[0]] = kv[1]
	}

	return res, nil
}

func underscoreName(s string) string {
	var buf []byte
	for i, v := range s {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}

	var err error
	m, err = r.parseModel(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), err
}

//func (r *registry) get1(val any) (*Model, error) {
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
