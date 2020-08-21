package lua

import (
	"FBI-Analyzer/logger"
	"FBI-Analyzer/rule"
	"time"

	"github.com/yuin/gopher-lua"
	"github.com/tr3ee/ngx-go"
)

var (
	Kchan = make(chan rule.AccessLog, 409600)
)

const (
	ReqType = "req"
)

// 启动一个线程启动lua虚拟机进行消费
func LuaThread(i int) {

	var luaFuncs = make(map[string]*lua.LFunction)
	logger.Print(logger.INFO, "lua thread no.%d started.", i)
	L := lua.NewState()
	defer L.Close()
	registerHttpType(L)
	L.PreloadModule("redis", luaRedis)
	L.PreloadModule("re", luaRe)
	L.PreloadModule("time", luaTime)
	// 监控规则更新
	go rule.SyncRule(L, luaFuncs)
	// 启动lua协程，执行所有缓存的方法
	for {
		select {
		case access := <-Kchan:
			luaCoroutines(L, access, luaFuncs)
		}
	}
}

// 新建lua协程
func luaCoroutines(L *lua.LState, access rule.AccessLog, luaFuncs map[string]*lua.LFunction) {

	var fn *lua.LFunction
	var err error
	var st lua.ResumeState
	//var values []lua.LValue
	co, _ := L.NewThread()
	registerAccessUserData(co, access)
	for _, fn = range luaFuncs {
		//L.Resume(co, fn)
		st, err, _ = L.Resume(co, fn)
		if st == lua.ResumeError {
			logger.Print(logger.ERROR, "coroutines failed : %v.", err.Error())
		}
		//if st == lua.ResumeOK {
		//logger.Print(logger.INFO, "coroutines ok")
		//}
	}
	co.Close()
}

// 测试消费者，日志ngxRaw格式
func TestConsumerRaw() {
	var sampleRaw = []byte(`1.1.1.1 - - [20/Aug/2020:16:44:22 +0800] "GET /api/queryKey?sdkAppId=x&platform=ios&method=HMACSHA256&gateway=x HTTP/1.1" 403 62 "-" "0.000" "-" "PostmanRuntime/7.6.0" "-" "post-body-data" "-" uid:"-x_uid=id_value" "test-1.host.x.net" "-"`)
	var sampleFmt = `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$request_time" "$upstream_response_time" "$http_user_agent" "$http_x_forwarded_for" "$request_body" "$http_accesstoken" uid:"$uid_got$uid_set" "$host" "$http_cookie"`
	var access rule.AccessLog
	var ngxc, _ = ngx.Compile(sampleFmt)

	ngxc.Unmarshal(sampleRaw, &access)
	for {
		access.Status += 1
		Kchan <- access
		time.Sleep(5 * time.Second)
	}
}