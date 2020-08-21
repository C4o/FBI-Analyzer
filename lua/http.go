package lua

import (
	"FBI-Analyzer/rule"

	"github.com/yuin/gopher-lua"
)

func GetReqVar(L *lua.LState) int {

	access := L.GetGlobal("access").(*lua.LUserData).Value.(*rule.AccessLog)
	_ = L.CheckAny(1)
	switch L.CheckString(2) {
	case "remote_addr":
		L.Push(lua.LString(access.RemoteAddr))
	case "remote_user":
		L.Push(lua.LString(access.RemoteUser))
	case "time_local":
		L.Push(lua.LString(access.TimeLocal))
	case "request":
		L.Push(lua.LString(access.Request))
	case "status":
		L.Push(lua.LNumber(access.Status))
	case "size":
		L.Push(lua.LNumber(access.BodyByteSent))
	case "referer":
		L.Push(lua.LString(access.HttpReferer))
	case "request_time":
		L.Push(lua.LString(access.RequestTime))
	case "urt":
		L.Push(lua.LString(access.URT))
	case "ua":
		L.Push(lua.LString(access.UA))
	case "xff":
		L.Push(lua.LString(access.XFF))
	case "body":
		L.Push(lua.LString(access.RequestBody))
	case "accesstoken":
		L.Push(lua.LString(access.AccessToken))
	case "uidgot":
		L.Push(lua.LString(access.UidGot))
	case "host":
		L.Push(lua.LString(access.Host))
	case "cookie":
		L.Push(lua.LString(access.Cookie))
	default:
		L.Push(lua.LNil)
	}
	return 1
}
