package lua

import (
	"FBI-Analyzer/logger"
	"FBI-Analyzer/rule"

	"github.com/go-redis/redis"
	"github.com/yuin/gopher-lua"
)

// 注册redis相关方法到lua虚拟机
func luaRedis(L *lua.LState) int {

	// 注册方法
	mod := L.SetFuncs(L.NewTable(), rdsFns)
	mod.RawSet(lua.LString("pipeline"), L.SetFuncs(L.NewTable(), plFns))
	L.Push(mod)
	return 1
}

// 注册re相关方法到lua虚拟机
func luaRe(L *lua.LState) int {

	// 注册方法
	mod := L.SetFuncs(L.NewTable(), reFns)
	L.Push(mod)
	return 1
}

// 注册time相关方法到lua虚拟机
func luaTime(L *lua.LState) int {

	// 注册方法
	mod := L.SetFuncs(L.NewTable(), timeFns)
	mod.RawSet(lua.LString("zero"), lua.LNumber(1590829200))
	L.Push(mod)
	return 1
}

// 注册Request到lua虚拟机
func registerHttpType(L *lua.LState) {

	var mt = L.NewTypeMetatable(ReqType)
	var reqTb = &lua.LTable{}
	var fbiTb = &lua.LTable{}
	// 设置req的表和元表
	L.SetFuncs(mt, map[string]lua.LGFunction{"__index": GetReqVar})
	reqTb.Metatable = mt
	// 设置fbi的表
	fbiTb.RawSet(lua.LString("var"), reqTb)
	// 设置logging方法和日志等级
	fbiTb.RawSet(lua.LString("log"), L.NewFunction(logging))
	fbiTb.RawSet(lua.LString("log2"), L.NewFunction(logging2))
	fbiTb.RawSet(lua.LString("ERROR"), lua.LNumber(logger.ERROR))
	fbiTb.RawSet(lua.LString("DEBUG"), lua.LNumber(logger.DEBUG))
	fbiTb.RawSet(lua.LString("INFO"), lua.LNumber(logger.INFO))
	L.SetGlobal("fbi", fbiTb)
}

// 注册userdata区域并填充日志数据
func registerAccessUserData(L *lua.LState, access rule.AccessLog) {

	var ud = L.NewUserData()
	ud.Value = &access
	L.SetGlobal("access", ud)
}

// 在新建pipeline的时候填充一块userdata
func registerPipelineUserData(L *lua.LState, pipeline redis.Pipeliner) {

	var ud = L.NewUserData()
	ud.Value = pipeline
	L.SetGlobal("pipeline", ud)
}
