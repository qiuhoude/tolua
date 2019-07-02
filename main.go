package main

import (
	"flag"
	"tolua/engine"
)

var outPut = flag.String("o", "luaOutput/", "output dir")

func main() {
	flag.Parse()
	e := engine.ConcurrentEngine{
		Output: *outPut,
	}
	e.Run()
}
