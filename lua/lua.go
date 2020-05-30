package lua

import (
	"FBI-Analyzer/logger"
	"FBI-Analyzer/rule"
	"time"

	"github.com/json-iterator/go"
	"github.com/yuin/gopher-lua"
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

// 测试消费者
func TestConsumer() {
	var sample = []byte(`{"host":"www.k3f.xyz","warn":"nil","saddr":"192.168.123.1","status":200,"xff":"0.0.0.0","rule":"nil","size":67,"method":"POST","uri":"\/interface\/GetData.aspx","reqs":"1","uaddr":"192.168.122.222:80","time":"1590469755.764","port":"33556","app":"X业务","cdn":"nil","addr":"101.102.103.104","urt":0.019,"pass":"nil","query":"key=0.aaa","remote":"101.102.103.104","ref":"http:\/\/www.k3f.xyz\/testuri\/0","ua":"Mozilla\/5.0 (Windows NT 10.0; WOW64) AppleWebKit\/537.36 (KHTML, like Gecko) Chrome\/73.0.3683.86 Safari\/537.36","risk":"ok","uname":"www.k3f.xyz","conn":"74691190","local":"26\/May\/2020:13:09:15 +0800"}`)
	var access rule.AccessLog
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	json.Unmarshal(sample, &access)
	for {
		access.Status += 1
		Kchan <- access
		time.Sleep(5 * time.Second)
	}
}
