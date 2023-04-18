package sql_demo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonColumn[T any] struct {
	Val T
	//indicates null
	Valid bool
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	//if null return nil
	if !j.Valid {
		return nil, nil
	}

	//marshal
	return json.Marshal(j.Val)
}

func (j *JsonColumn[T]) Scan(src any) error {
	//according to type
	var bs []byte
	switch data := src.(type) {
	case string:
		bs = []byte(data)
	case []byte:
		bs = data
	case nil:
		return nil
	default:
		return errors.New("orm: 不支持类型")
	}

	err := json.Unmarshal(bs, &j.Val)
	if err == nil {
		j.Valid = true
	}
	return err
}
