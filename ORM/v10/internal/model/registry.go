package model

import (
	"geektime/ORM/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type Registry struct {
	//models map[reflect.Type]*Model
	//mutex  sync.RWMutex
	models sync.Map
}

var _ Register = &Registry{}

type Register interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOpt) (*Model, error)
}

func (r *Registry) Register(val any, opts ...ModelOpt) (*Model, error) {
	m, err := r.ParseModel(val)
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

func (r *Registry) ParseModel(entity any) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	numField := typ.NumField()
	fields := make(map[string]*Field, numField)
	columns := make(map[string]*Field, numField)
	colSlice := make([]*Field, 0, numField)

	//get all columns
	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		//get column name
		m, err := r.parseTag(f.Tag)
		if err != nil {
			return nil, err
		}
		colName := m[TagKeyColumn]
		if colName == "" {
			colName = underscoreName(f.Name)
		}

		fields[f.Name] = &Field{ColName: colName, Offset: f.Offset} //, Typ: f.Type}
		columns[colName] = &Field{FieldName: f.Name, Typ: f.Type, Offset: f.Offset}
		colSlice = append(colSlice, &Field{ColName: colName, FieldName: f.Name, Offset: f.Offset})
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
		TableName: tblName,
		Fields:    fields,
		Columns:   columns,
		ColSlice:  colSlice,
	}, nil
}

func (r *Registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
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

// Get return *Model of a struct
func (r *Registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}

	var err error
	m, err = r.ParseModel(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), err
}

//func (R *Registry) get1(val any) (*Model, error) {
//	//use double check
//	Typ := reflect.TypeOf(val)
//	R.mutex.RLock()
//	m, ok := R.models[Typ]
//	R.mutex.RUnlock()
//	if ok {
//		return m, nil
//	}
//
//	R.mutex.Lock()
//	defer R.mutex.Unlock()
//	m, ok = R.models[Typ]
//	//再读一遍，如果ok则自己返回
//	if ok {
//		return m, nil
//	}
//	if !ok {
//		var err error
//		m, err = R.ParseModel(val)
//		if err != nil {
//			return nil, err
//		}
//	}
//	R.models[Typ] = m
//	return m, nil
//}
