package _reflect

import (
	"errors"
	"reflect"
)

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	//检验类型
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Struct && typ.Kind() != reflect.Pointer {
		return nil, errors.New("非法类型")
	}
	//获取entity 相关
	numMtd := typ.NumMethod()
	result := make(map[string]FuncInfo, numMtd)

	//对entity方法进行迭代
	for i := 0; i < numMtd; i++ {
		//循环中，记录每个方法的输入输出类型，零值输入的返回值
		f := typ.Method(i)

		//记录输入
		numIn := f.Type.NumIn()
		inputTypes := make([]reflect.Type, 0, numIn)

		//零值输入的slice
		iptVals := make([]reflect.Value, 0, numIn)
		iptVals = append(iptVals, reflect.ValueOf(entity)) //这里错了，写成typ

		//第一个是接收器自身
		inputTypes = append(inputTypes, f.Type.In(0))
		for j := 1; j < numIn; j++ {
			inT := f.Type.In(j)                          //这里错了，写成typ.In(j)
			iptVals = append(iptVals, reflect.Zero(inT)) //这里 reflect.ValueOf(reflect.Zero(inT))
			inputTypes = append(inputTypes, inT)
		}

		//调用函数
		values := f.Func.Call(iptVals)
		//记录输出

		outNum := f.Type.NumOut()
		outputTypes := make([]reflect.Type, 0, outNum)
		outputValues := make([]any, 0, outNum)
		for k := 0; k < outNum; k++ {
			outT := f.Type.Out(k) //这里错了，写成typ.Out(j)
			outputTypes = append(outputTypes, outT)
			outputValues = append(outputValues, values[k].Interface())
		}

		//append to result
		result[f.Name] = FuncInfo{
			Name:         f.Name,
			InputTypes:   inputTypes,
			OutputTypes:  outputTypes,
			OutputValues: outputValues,
		}
	}
	return result, nil
}

type FuncInfo struct {
	Name         string
	InputTypes   []reflect.Type
	OutputTypes  []reflect.Type
	OutputValues []any
}
