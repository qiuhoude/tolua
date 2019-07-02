package engine

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"tolua/database"
	"tolua/persist"
	"tolua/work"
)

type ConcurrentEngine struct {
	Output string
}

func (c *ConcurrentEngine) Run() {
	t1 := time.Now()
	database.InitDB() //初始化数据库
	if c.Output == "" {
		panic("output 路径不能为空")
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	tableNames := searchTableNames()
	var wg sync.WaitGroup
	p := work.New(len(tableNames) / 4)
	for i, tableName := range tableNames {
		wg.Add(1)
		f := createTaskFunc(c, tableName, i)
		go func() {
			p.Run(f)
			wg.Done()
		}()
	}
	wg.Wait()
	p.Shutdown()
	costTime := time.Since(t1)
	fmt.Println("功花掉了生命 ", costTime)
	// 单携程运行
	//r := runner.New(10 * time.Second)
	//for _, tname := range tableNames {
	//	task := createTask(c, tname)
	//	r.Add(task)
	//}
	//r.Start()
}

type taskFunc func()

func (f taskFunc) Task() {
	f()
}

func createTaskFunc(c *ConcurrentEngine, tableName string, index int) taskFunc {
	return func() {
		d := queryDataByTable(tableName)
		err := persist.SaveDataToFile(c.Output, d)
		if err == nil {
			fmt.Printf("------------【%d】保存 %s 表数据完成-----------\n", index, tableName)
		} else {
			fmt.Printf("------------【%d】保存 %s 表数据失败-----------\n", index, tableName)
		}
	}
}

// 创建多个
func createTask(c *ConcurrentEngine, tableName string) func(int) {
	return func(id int) {
		d := queryDataByTable(tableName)
		//persist.MarshalLua(os.Stdout,d)
		err := persist.SaveDataToFile(c.Output, d)
		if err == nil {
			fmt.Printf("------------【%d】保存 %s 表数据完成-----------\n", id, tableName)
		} else {
			fmt.Printf("------------【%d】保存 %s 表数据失败-----------\n", id, tableName)
		}
	}
}
