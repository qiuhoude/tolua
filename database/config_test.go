package database

import (
	"encoding/json"
	"testing"
)

func TestConfig(t *testing.T) {
	c := dbConf{
		User:    "root",
		Pwd:     "root",
		Host:    "127.0.0.1",
		Port:    "3306",
		DBName:  "honor_ini",
		DBParam: "",
	}
	conf, _ := json.MarshalIndent(c, "", "	")
	t.Logf("\n%s", conf)
}
