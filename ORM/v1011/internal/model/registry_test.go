package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestRegistry_get(t *testing.T) {

	testCases := []struct {
		name      string
		val       any
		wantModel *Model
		wantErr   error
	}{
		//{
		//	name:    "test Model",
		//	val:     TestModel{},
		//	wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		//},
		{
			// 指针
			name: "pointer",
			val:  &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				Fields: map[string]*Field{
					//"Id": {
					//	ColName: "id",
					//},
					//"FirstName": {
					//	ColName: "first_name",
					//},
					"Age": {
						ColName: "age",
					},
					"Name": {
						ColName: "name",
					},
				},
				Columns: map[string]*Field{
					//"id":         {FieldName: "Id", Typ: reflect.TypeOf(12)},
					//"first_name": {FieldName: "FirstName", Typ: reflect.TypeOf("")},
					"age":  {FieldName: "Age", Typ: reflect.TypeOf(1)},
					"name": {FieldName: "Name", Typ: reflect.TypeOf("")},

					//"last_name":  {FieldName: "LastName", Typ: reflect.TypeOf("")},
				},
			},
		},
		{
			// 多级指针
			name: "multiple pointer",
			// 因为 Go 编译器的原因，所以我们写成这样
			val: func() any {
				val := TestModel{}
				return &val
			}(),
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			name:    "map",
			val:     map[string]string{},
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			name:    "slice",
			val:     []int{},
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			name:    "basic type",
			val:     0,
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},

		//test tag functions
		{
			name: "normal tag",
			val: func() any {
				type MyTag struct {
					ColName string `orm:"column=new_col_name""`
				}

				return &MyTag{}
			}(),
			wantModel: &Model{
				TableName: "my_tag",
				Fields: map[string]*Field{
					"ColName": &Field{ColName: "new_col_name"},
				},
				Columns: map[string]*Field{
					"new_col_name": {FieldName: "ColName", Typ: reflect.TypeOf("")},
				},
			},
		},
		{
			name: "without tag",
			val: func() any {
				type MyTag struct {
					ColName string
				}

				return &MyTag{}
			}(),
			wantModel: &Model{
				TableName: "my_tag",
				Fields: map[string]*Field{
					"ColName": &Field{ColName: "col_name"},
				},
				Columns: map[string]*Field{
					"col_name": {FieldName: "ColName", Typ: reflect.TypeOf("")},
				},
			},
		},
		{
			name: "normal table name",
			val:  &CustomTableName{},
			wantModel: &Model{
				TableName: "new_custom_name",
				Fields: map[string]*Field{
					"Name": &Field{ColName: "name"},
				},
				Columns: map[string]*Field{
					"name": {FieldName: "Name", Typ: reflect.TypeOf("")},
				},
			},
		},
		{
			name: "normal table name without ptr",
			val:  CustomTableName{},
			wantModel: &Model{
				TableName: "new_custom_name",
				Fields: map[string]*Field{
					"Name": &Field{ColName: "name"},
				},
				Columns: map[string]*Field{
					"name": {FieldName: "Name", Typ: reflect.TypeOf("")},
				},
			},
		},
	}

	//构建db
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	r := &Registry{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(reflect.TypeOf(tc.val))
			_, err := r.Get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			//if err != nil {
			//	return
			//}
			//assert.Equal(t, tc.wantModel, m)
			//col, _ := m.Columns["last_name"]
			//bytes, _ := json.Marshal(col)
			//fmt.Println(string(bytes))
		})
	}
}

func Test_underscoreName(t *testing.T) {
	testCases := []struct {
		name    string
		srcStr  string
		wantStr string
	}{
		// 我们这些用例就是为了确保
		// 在忘记 underscoreName 的行为特性之后
		// 可以从这里找回来
		// 比如说过了一段时间之后
		// 忘记了 ID 不能转化为 id
		// 那么这个测试能帮我们确定 ID 只能转化为 i_d
		{
			name:    "upper cases",
			srcStr:  "ID",
			wantStr: "i_d",
		},
		{
			name:    "use number",
			srcStr:  "Table1Name",
			wantStr: "table1_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := underscoreName(tc.srcStr)
			assert.Equal(t, tc.wantStr, res)
		})
	}
}

type CustomTableName struct {
	Name string
}

func (c CustomTableName) TableName() string {
	return "new_custom_name"
}
