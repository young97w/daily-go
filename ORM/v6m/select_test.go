package v1

import (
	"context"
	"geektime/ORM/internal/errs"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	//LastName  string //*sql.NullString
}

func TestSelector_Build(t *testing.T) {
	//构建db
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, _ := OpenDB(sqlDB)
	testCase := []struct {
		name      string
		builder   *Selector[TestModel]
		wantQuery *Query
		wantErr   error
	}{
		{
			name:      "no from",
			builder:   NewSelector[TestModel](db),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model`;", Args: make([]any, 0)},
		},
		{
			name:      "with from",
			builder:   NewSelector[TestModel](db).From("`test_table`"),
			wantQuery: &Query{SQL: "SELECT * FROM `test_table`;", Args: make([]any, 0)},
		},
		{
			name:      "empty form",
			builder:   NewSelector[TestModel](db).From(""),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model`;", Args: make([]any, 0, 4)},
		},
		{
			// 单一简单条件
			name:    "single and simple predicate",
			builder: NewSelector[TestModel](db).From("`test_model_t`").Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name:    "multiple predicates",
			builder: NewSelector[TestModel](db).Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name:    "and",
			builder: NewSelector[TestModel](db).Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name:    "or",
			builder: NewSelector[TestModel](db).Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name:    "not",
			builder: NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` > ?);",
				Args: []any{18},
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestModelWithTableName(t *testing.T) {
	//构建db
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, _ := OpenDB(sqlDB)
	m, err := db.r.Register(TestModel{}, ModelWithTableName("new_table"))
	require.NoError(t, err)
	assert.Equal(t, "new_table", m.tableName)
}

func TestModelWithColumnName(t *testing.T) {
	//构建db
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, _ := OpenDB(sqlDB)
	testCase := []struct {
		name    string
		field   string
		colName string

		wantModel *Model
		wantErr   error
	}{
		{
			name:    "test model",
			field:   "Id",
			colName: "new_id",
			wantModel: &Model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "new_id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					//"LastName": {
					//	colName: "last_name",
					//},
				},
			},
		},
		{
			name:    "empty column name with error",
			field:   "Id",
			colName: "",
			wantErr: errs.NewErrUnknownField("Id"),
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			m, err := db.r.Register(TestModel{}, ModelWithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel.fields, m.fields)
		})
	}
}

type mockModel struct {
	Id        int
	FirstName string
}

func TestSelector_Get(t *testing.T) {
	//构建db
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	//new DB
	db, err := OpenDB(sqlDB)

	defer sqlDB.Close()
	require.NoError(t, err)

	//构建表数据
	mockRows := sqlmock.NewRows([]string{"id", "first_name"})
	mockRows.AddRow(1, "young")

	//构建sql表达式返回结果
	mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)

	//rows, _ := db.db.QueryContext(context.Background(), "SELECT *")
	//for rows.Next() {
	//	mm := mockModel{}
	//	rows.Scan(&mm.Id, &mm.FirstName)
	//	bytes, _ := json.Marshal(mm)
	//	fmt.Println(string(bytes))
	//}

	testCase := []struct {
		name    string
		builder *Selector[mockModel]

		wantErr error
		wantRes *mockModel
	}{
		{
			name:    "normal model",
			builder: NewSelector[mockModel](db),
			wantRes: &mockModel{
				Id:        1,
				FirstName: "young",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.builder.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}
}
