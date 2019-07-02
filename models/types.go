package models

type ColumnVal struct {
	Data []interface{} //每一行的数据,都是数据库中的类型
}

// 每张表的数据
type TableData struct {
	TableName    string      // 表名
	ColumnName   []string    // 表头
	ColumnValues []ColumnVal // 数据
}


