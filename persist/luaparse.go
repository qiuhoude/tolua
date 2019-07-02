package persist

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"tolua/models"
)

var luaFuncMap = template.FuncMap{
	"printSlice": printSlice,
	"printTitle": printTitle,
	"printComma": printComma,
}

func MarshalLua(wr io.Writer, data *models.TableData) error {
	luaTemplate := template.Must(template.New("").Funcs(luaFuncMap).ParseFiles("config/luatemp.txt"))
	return luaTemplate.ExecuteTemplate(wr, "lua", data)
}

func SaveDataToFile(outDir string, data *models.TableData) error {
	checkDir(outDir)
	f := filepath.Join(outDir, data.TableName+".lua")
	file, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	return MarshalLua(file, data)
}

func checkDir(dir string) {
	if exists, err := PathExists(dir); !exists || err != nil {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func printComma(len, index int) string {
	if index < len-1 {
		return ","
	} else {
		return ""
	}
}

// 打印表头
func printTitle(titles []string) string {
	sb := strings.Builder{}
	length := len(titles)
	for i, t := range titles {
		sb.WriteRune('"')
		sb.WriteString(t)
		sb.WriteRune('"')
		if i < length-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func printSlice(data []interface{}) string {
	var bf bytes.Buffer
	length := len(data)
	for i, v := range data {
		switch t := v.(type) {
		case int, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
			bf.WriteString(fmt.Sprintf(`%v`, t))
		case string:
			if t == "" {
				bf.WriteString("nil")
			} else {
				bf.WriteString(fmt.Sprintf(`"%s"`, t))
			}
		case []byte:
			bf.WriteString(fmt.Sprintf(`"%s"`, string(t)))
		default:
			bf.WriteString(fmt.Sprintf(`"%v"`, t))
		}
		if i < length-1 {
			bf.WriteRune(',')
		}
	}
	return bf.String()
}
