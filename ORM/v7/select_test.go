package v1

import (
	"context"
	"fmt"
	"geektime/ORM/internal/errs"
	"geektime/ORM/v7/internal/model"
	"geektime/ORM/v7/internal/valuer"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelector_Select(t *testing.T) {
	//构建db
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, _ := OpenDB(sqlDB)

	testCases := []struct {
		name      string
		cols      []Selectable
		builder   *Selector[TestModel]
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "normal",
			cols:    []Selectable{C("Id"), C("FirstName")},
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT `id`,`first_name` FROM `test_model`;",
				Args: []any{},
			},
		},
		{
			name:    "invalid column",
			cols:    []Selectable{C("invalid")},
			builder: NewSelector[TestModel](db),
			wantErr: errs.NewErrUnknownField("invalid"),
		},
		{
			name:    "AVG age",
			cols:    []Selectable{Avg("Age")},
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT AVG(`age`) FROM `test_model`;",
				Args: []any{},
			},
		},
		{
			name:    "DISTINCT",
			cols:    []Selectable{Raw("COUNT(DISTINCT `first_name`)", 18)},
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
				Args: []any{18},
			},
		},
		{
			name:    "as name",
			cols:    []Selectable{C("FirstName").As("name")},
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT `first_name` AS `name` FROM `test_model`;",
				Args: []any{},
			},
		},
		{
			name:    "AVG age as avg_age",
			cols:    []Selectable{Avg("Age").As("avg_age")},
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT AVG(`age`) AS `avg_age` FROM `test_model`;",
				Args: []any{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Select(tc.cols...).Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantQuery, q)
		})
	}
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
		{
			name:    "raw expression used as predicate",
			builder: NewSelector[TestModel](db).Where(Raw("`id`<?", 18).AsPredicate()),
			wantQuery: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `test_model` WHERE (`id`<?);",
				Args: []any{18},
			},
		},
		{
			name:    "raw expression used in predicate",
			builder: NewSelector[TestModel](db).Where(C("Id").EQ(Raw("`age`+?", 1))),
			wantQuery: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `test_model` WHERE `id` = (`age`+?);",
				Args: []any{1},
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
	m, err := db.R.Register(TestModel{}, model.ModelWithTableName("new_table"))
	require.NoError(t, err)
	assert.Equal(t, "new_table", m.TableName)
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

		wantModel *model.Model
		wantErr   error
	}{
		{
			name:    "test model",
			field:   "Id",
			colName: "new_id",
			wantModel: &model.Model{
				TableName: "test_model",
				Fields: map[string]*model.Field{
					"Id": {
						ColName: "new_id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"Age": {
						ColName: "age",
					},
					//"LastName": {
					//	ColName: "last_name",
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
			_, err := db.R.Register(TestModel{}, model.ModelWithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			//if err != nil {
			//	return
			//}
			//assert.Equal(t, tc.wantModel.Fields, m.Fields)
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
	db.valCreator = valuer.NewUnsafeValue

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

func BenchmarkSelector_Get(b *testing.B) {
	db, err := Open("sqlite3", fmt.Sprintf("file:benchmark_get.db?cache=shared&mode=memory"))
	if err != nil {
		b.Fatal(err)
	}

	_, err = db.db.Exec(TestModel{}.CreateSQL())
	if err != nil {
		b.Fatal(err)
	}

	res, err := db.db.Exec("INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?) ", 1, "young", 18, "sky")
	if err != nil {
		b.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		b.Fatal(err)
	}
	if affected == 0 {
		b.Fatal()
	}

	b.Run("unsafe", func(b *testing.B) {
		db.valCreator = valuer.NewUnsafeValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("reflect", func(b *testing.B) {
		db.valCreator = valuer.NewReflectValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})

}
