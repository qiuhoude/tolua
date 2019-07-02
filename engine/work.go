// 查询数据库
package engine

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"tolua/database"
	"tolua/models"
)

// 查询所有表名
func searchTableNames() []string {
	rows, err := database.DB.Query("SHOW TABLES")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//读出查询出的列字段名
	var tableNames []string
	for rows.Next() {
		var tname string
		rows.Scan(&tname)
		tableNames = append(tableNames, tname)
	}
	return tableNames
}

func queryTableType(tableName string) (map[string]database.DBDataType, error) {
	rows, err := database.DB.Query("SELECT COLUMN_NAME,DATA_TYPE FROM `information_schema`.`COLUMNS`  WHERE TABLE_NAME=?", tableName)
	if err != nil {
		return nil, err
	}
	res := make(map[string]database.DBDataType)
	for rows.Next() {
		var columnName, datatype string
		rows.Scan(&columnName, &datatype)
		switch typeStr := strings.ToLower(datatype); typeStr {
		case "int":
			res[columnName] = database.INT
		case "float":
			res[columnName] = database.FLOAT
		default:
			res[columnName] = database.STRING
		}
	}
	return res, nil
}

// 查询单表
func queryDataByTable(tableName string) *models.TableData {
	typeMap, err := queryTableType(tableName)
	if err != nil {
		fmt.Println(tableName, "表查询表结构出错 -> ", err.Error())
		return nil
	}
	rows, err := database.DB.Query("SELECT * FROM " + tableName)
	if err != nil {
		fmt.Println(tableName, "表查询出错 -> ", err.Error())
		return nil
	}
	// table数据
	tableData := models.TableData{
		TableName: tableName,
	}
	//读出查询出的列字段名
	cols, _ := rows.Columns()
	tableData.ColumnName = cols
	colsLen := len(cols)
	values := make([]sql.RawBytes, colsLen)
	// 每个字段的指针地址
	columnPointers := make([]interface{}, colsLen)
	// 赋值指针
	for i := range values {
		columnPointers[i] = &values[i]
	}
	var res []models.ColumnVal
	for rows.Next() {
		if err := rows.Scan(columnPointers...); err != nil {
			log.Println(err)
			continue
		}
		tv := models.ColumnVal{
			Data: make([]interface{}, colsLen),
		}
		for i, col := range values {
			switch typeMap[cols[i]] {
			case database.INT:
				if col == nil {
					tv.Data[i] = 0
				} else {
					f, _ := strconv.Atoi(string(col))
					tv.Data[i] = f
				}
			case database.FLOAT:
				if col == nil {
					tv.Data[i] = 0
				} else {
					f, _ := strconv.ParseFloat(string(col), 32/64)
					tv.Data[i] = f
				}
			case database.STRING:
				tv.Data[i] = string(col)
			}
		}
		res = append(res, tv)
	}
	tableData.ColumnValues = res
	return &tableData
}
