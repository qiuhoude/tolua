package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
)

type dbConf struct {
	User    string `json:"user"`
	Pwd     string `json:"pwd"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	DBName  string `json:"dbName"`
	DBParam string `json:"dbParam"`
}

func (c dbConf) toDbUrl() string {
	var dbUrl string

	if c.DBName == "" {
		dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", c.User, c.Pwd, c.Host, c.Port, c.DBName)
	} else {
		dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&%s", c.User, c.Pwd, c.Host, c.Port, c.DBName, c.DBParam)
	}
	return dbUrl
}

var DB *sql.DB

func getDBConf() *dbConf {
	var c dbConf
	bytes, err := ioutil.ReadFile("config/DBConf.json")
	if err != nil {
		log.Fatal("读取配置出错 -> ", err)
	}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		log.Fatal("解析配置出错 -> ", err)
	}
	return &c
}

func InitDB() {
	dbConf := getDBConf()
	dbUrl := dbConf.toDbUrl()
	var err error
	DB, err = sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatal("打开数据库失败 ->", err)
	}
}

type DBDataType uint

const (
	INT DBDataType = iota
	FLOAT
	STRING  // 其他类型都是字符串
)
