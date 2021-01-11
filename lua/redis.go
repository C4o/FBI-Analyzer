package lua

import (
	"FBI-Analyzer/db"
	"time"

	"github.com/go-redis/redis"
	"github.com/yuin/gopher-lua"
)

// golang实现的所有可在lua中使用的redis操作方法
var (
	rdsFns = map[string]lua.LGFunction{
		"incr":   incr,
		"hmget":  hmget,
		"hmset":  hmset,
		"expire": expire,
		"delete": delete,
	}
	plFns = map[string]lua.LGFunction{
		"new":    pnew,
		"exec":   pexec,
		"close":  pclose,
		"incr":   pincr,
		"hmget":  phmget,
		"hmset":  phmset,
		"expire": pexpire,
		"delete": pdelete,
	}
)

// 通过redis:func(n1, n2)相当于redis.func(self, n1, n2)
// 强制使用redis.func(), 否则报错

func pushErr(L *lua.LState, err error) {

	if err == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(err.Error()))
	}
}

func typeCheck(it interface{}) lua.LValue {

	var ok bool

	if it == nil {
		return lua.LNil
	} else if _, ok = it.(string); ok {
		return lua.LString(it.(string))
	} else if _, ok = it.(int); ok {
		return lua.LNumber(it.(int))
	} else {
		return lua.LNumber(1)
	}

}

func incr(L *lua.LState) int {

	var err error
	var result int64
	result, err = db.RedSess.HIncrBy(L.CheckString(1), L.CheckString(2), 1).Result()
	L.Push(lua.LNumber(result))
	pushErr(L, err)
	return 2
}

func hmset(L *lua.LState) int {

	var err error
	var result string
	fmap := make(map[string]interface{})
	fmap[L.CheckString(2)] = L.CheckInt(3)
	result, err = db.RedSess.HMSet(L.CheckString(1), fmap).Result()
	L.Push(lua.LString(result))
	pushErr(L, err)
	return 2
}

func hmget(L *lua.LState) int {

	var err error
	var results []interface{}
	results, err = db.RedSess.HMGet(L.CheckString(1), L.CheckString(2)).Result()
	L.Push(typeCheck(results[0]))
	pushErr(L, err)
	return 2
}

func expire(L *lua.LState) int {

	var err error
	var result bool
	result, err = db.RedSess.Expire(L.CheckString(1), time.Duration(L.CheckInt(2))*time.Second).Result()
	L.Push(lua.LBool(result))
	pushErr(L, err)
	return 2
}

func delete(L *lua.LState) int {

	var err error
	var result int64
	result, err = db.RedSess.Del(L.CheckString(1)).Result()
	L.Push(lua.LNumber(result))
	pushErr(L, err)
	return 2
}

// redis pipeline golang for lua实现
func pnew(L *lua.LState) int {

	var pler = db.RedSess.Pipeline()
	registerPipelineUserData(L, pler)
	return 0
}

func pexec(L *lua.LState) int {

	var err error
	pipeline := L.GetGlobal("pipeline").(*lua.LUserData).Value.(redis.Pipeliner)
	_, err = pipeline.Exec()
	pushErr(L, err)
	return 1
}

func pclose(L *lua.LState) int {

	pipeline := L.GetGlobal("pipeline").(*lua.LUserData).Value.(redis.Pipeliner)
	pipeline.Close()
	return 0
}

func pincr(L *lua.LState) int {

	var err error
	var result int64
	pipeline := L.GetGlobal("pipeline").(*lua.LUserData).Value.(redis.Pipeliner)
	result, err = pipeline.HIncrBy(L.CheckString(1), L.CheckString(2), 1).Result()
	L.Push(lua.LNumber(result))
	pushErr(L, err)
	return 2
}

func phmset(L *lua.LState) int {

	var err error
	var result string
	pipeline := L.GetGlobal("pipeline").(*lua.LUserData).Value.(redis.Pipeliner)
	fmap := make(map[string]interface{})
	fmap[L.CheckString(2)] = L.CheckInt(3)
	result, err = pipeline.HMSet(L.CheckString(1), fmap).Result()
	L.Push(lua.LString(result))
	pushErr(L, err)
	return 2
}

func phmget(L *lua.LState) int {

	var err error
	var results []interface{}
	pipeline := L.GetGlobal("pipeline").(*lua.LUserData).Value.(redis.Pipeliner)
	results, err = pipeline.HMGet(L.CheckString(1), L.CheckString(2)).Result()
	L.Push(typeCheck(results[0]))
	pushErr(L, err)
	return 2
}

func pexpire(L *lua.LState) int {

	var err error
	var result bool
	pipeline := L.GetGlobal("pipeline").(*lua.LUserData).Value.(redis.Pipeliner)
	result, err = pipeline.Expire(L.CheckString(1), time.Duration(L.CheckInt(2))*time.Second).Result()
	L.Push(lua.LBool(result))
	pushErr(L, err)
	return 2
}

func pdelete(L *lua.LState) int {

	var err error
	var result int64
	pipeline := L.GetGlobal("pipeline").(*lua.LUserData).Value.(redis.Pipeliner)
	result, err = pipeline.Del(L.CheckString(1)).Result()
	L.Push(lua.LNumber(result))
	pushErr(L, err)
	return 2
}

// 判断栈顶第一位是否是Table类型
//func typeCheck(it interface{}) lua.LValue {

//var ok bool

//if it == nil {
//return lua.LNil
//} else if _, ok = it.(string); ok {
//return lua.LString(it.(string))
//} else if _, ok = it.(int); ok {
//return lua.LNumber(it.(int))
//} else {
//return lua.LNumber(1)
//}

//}

//func incr(L *lua.LState) int {

//var err error
//var result int64
//if L.Get(1).Type() == lua.LTTable {
//result, err = db.RedSess.HIncrBy(L.CheckString(2), L.CheckString(3), 1).Result()
//} else {
//result, err = db.RedSess.HIncrBy(L.CheckString(1), L.CheckString(2), 1).Result()
//}
//L.Push(lua.LNumber(result))
//pushErr(L, err)
//return 2
//}

//func hmset(L *lua.LState) int {

//var err error
//var result bool
//if L.Get(1).Type() == lua.LTTable {
//result, err = db.RedSess.HMSet(L.CheckString(2), L.CheckString(3), L.CheckInt(4)).Result()
//} else {
//result, err = db.RedSess.HMSet(L.CheckString(1), L.CheckString(2), L.CheckInt(3)).Result()
//}
//L.Push(lua.LBool(result))
//pushErr(L, err)
//return 2
//}

//func hmget(L *lua.LState) int {

//var err error
//var results []interface{}
//if L.Get(1).Type() == lua.LTTable {
//results, err = db.RedSess.HMGet(L.CheckString(2), L.CheckString(3)).Result()
//} else {
//results, err = db.RedSess.HMGet(L.CheckString(1), L.CheckString(2)).Result()
//}
//L.Push(typeCheck(results[0]))
//pushErr(L, err)
//return 2
//}

//func expire(L *lua.LState) int {

//var err error
//var result bool
//if L.Get(1).Type() == lua.LTTable {
//result, err = db.RedSess.Expire(L.CheckString(2), time.Duration(L.CheckInt(3))*time.Second).Result()
//} else {
//result, err = db.RedSess.Expire(L.CheckString(1), time.Duration(L.CheckInt(2))*time.Second).Result()
//}
//L.Push(lua.LBool(result))
//pushErr(L, err)
//return 2
//}

//func delete(L *lua.LState) int {

//var err error
//var result int64
//if L.Get(1).Type() == lua.LTTable {
//result, err = db.RedSess.Del(L.CheckString(2)).Result()
//} else {
//result, err = db.RedSess.Del(L.CheckString(1)).Result()
//}
//L.Push(lua.LNumber(result))
//pushErr(L, err)
//return 2
//}
