package lua

import (
	"FBI-Analyzer/rule"

	"github.com/yuin/gopher-lua"
)

func GetReqVar(L *lua.LState) int {

	access := L.GetGlobal("access").(*lua.LUserData).Value.(*rule.AccessLog)
	_ = L.CheckAny(1)
	switch L.CheckString(2) {
	case "host":
		L.Push(lua.LString(access.Host))
	case "status":
		L.Push(lua.LNumber(access.Status))
	case "XFF":
		L.Push(lua.LString(access.XFF))
	case "rule":
		L.Push(lua.LString(access.Rule))
	case "size":
		L.Push(lua.LNumber(access.Size))
	case "method":
		L.Push(lua.LString(access.Method))
	case "uri":
		L.Push(lua.LString(access.URI))
	case "reqs":
		L.Push(lua.LString(access.Reqs))
	case "uaddr":
		L.Push(lua.LString(access.Uaddr))
	case "time":
		L.Push(lua.LString(access.Time))
	case "port":
		L.Push(lua.LString(access.Port))
	case "app":
		L.Push(lua.LString(access.APP))
	case "cdn":
		L.Push(lua.LString(access.CDN))
	case "addr":
		L.Push(lua.LString(access.Addr))
	case "urt":
		L.Push(lua.LNumber(access.URT))
	case "pass":
		L.Push(lua.LString(access.Pass))
	case "query":
		L.Push(lua.LString(access.Query))
	case "remote":
		L.Push(lua.LString(access.Remote))
	case "ref":
		L.Push(lua.LString(access.REF))
	case "ua":
		L.Push(lua.LString(access.UA))
	case "risk":
		L.Push(lua.LString(access.Risk))
	case "uname":
		L.Push(lua.LString(access.Uname))
	case "conn":
		L.Push(lua.LString(access.Conn))
	case "ltime":
		L.Push(lua.LString(access.Local))
	default:
		L.Push(lua.LNil)
	}

	return 1
}
