package lua

import (
	"FBI-Analyzer/logger"
	"bytes"
	"time"

	"github.com/yuin/gopher-lua"
)

var timeFns = map[string]lua.LGFunction{
	"unix":   tunix,
	"format": tformat,
}

func logging(L *lua.LState) int {

	buf := new(bytes.Buffer)
	n := L.GetTop()

	for i := 2; i < n+1; i++ {
		buf.WriteString(L.CheckString(i))
		buf.WriteString(" ")
	}

	logger.Println(L.CheckInt(1), buf.String())
	return 0
}

func logging2(L *lua.LState) int {

	n := L.GetTop()
	buf := make([]interface{}, n-1)

	for i := 2; i < n+1; i++ {
		buf[i-2] = L.CheckAny(i)
	}

	logger.Println(L.CheckInt(1), buf...)
	return 0
}

func tunix(L *lua.LState) int {

	L.Push(lua.LNumber(time.Now().Unix()))
	return 1
}

func tformat(L *lua.LState) int {

	L.Push(lua.LString(time.Now().Format("2006-01-02 15:04:05")))
	return 1
}
