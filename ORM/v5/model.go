package v1

// 我们支持的全部标签上的 key 都放在这里
// 方便用户查找，和我们后期维护
const (
	tagKeyColumn = "column"
)

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}

type Model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string
}
