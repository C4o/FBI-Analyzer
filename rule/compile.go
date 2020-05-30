package rule

import (
	"FBI-Analyzer/logger"
	"bufio"
	"os"

	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/ast"
	"github.com/yuin/gopher-lua/parse"
)

// 编译lua脚本并缓存方法
func compileLua(L *lua.LState, filepath string, luaFuncs map[string]*lua.LFunction) {

	var err error
	var file *os.File
	var chunk []ast.Stmt
	var proto *lua.FunctionProto

	file, err = os.OpenFile(filepath, os.O_RDONLY, 0444)
	if err != nil {
		logger.Print(logger.ERROR, "script %s not found!", filepath)
		return
	}
	defer file.Close()
	chunk, err = parse.Parse(bufio.NewReader(file), filepath)
	if err != nil {
		logger.Print(logger.ERROR, "parse script %s failed for %v.", filepath, err)
		return
	}
	proto, err = lua.Compile(chunk, filepath)
	if err != nil {
		logger.Print(logger.ERROR, "compile script %s failed for %v.", filepath, err)
		return
	}

	luaFuncs[filepath] = L.NewFunctionFromProto(proto)
	return
}
