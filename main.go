package main

import (
	"FBI-Analyzer/conf"
	"FBI-Analyzer/logger"
	"FBI-Analyzer/lua"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

var AgentInfo = `[*] 更新日志-20200528

[+] 2020年5月24日
    完成内置redis方法,incr/hmset/hmget/expire/delete
    完成内置变量fbi.var,可获取各类access日志参数

[+] 2020年5月26日
    完成log方法,fbi.log(fbi.ERROR/INFO/DEBUG, "s1", "s2", "s3", ... , "sn")
    定义方法只能通过类似redis.incr()调用,否则报错

[+] 2020年5月29日
    完成redis的pipeline方法,使用方式redis.pipeline.incr()

[+] 2020年5月29日
    完成正则匹配方法,find和match
`

func main() {

	var err error
	// 配置文件初始化
	c := flag.String("c", "conf.yaml", "configure file path")
	v := flag.Bool("v", false, "agent info")
	d := flag.Bool("d", false, "debug switch to open pprof")

	flag.Parse()

	if *d {
		var fc, fm *os.File
		// 使用pprof
		fc, err = os.Create("cpu.prof")
		if err != nil {
			return
		}
		defer fc.Close() // error handling omitted for example
		if err = pprof.StartCPUProfile(fc); err != nil {
			return
		}

		fm, err = os.Create("mem.prof")
		if err != nil {
			return
		}
		defer fm.Close() // error handling omitted for example
		runtime.GC()     // get up-to-date statistics
		if err = pprof.WriteHeapProfile(fm); err != nil {
			return
		}

		signalChan := make(chan os.Signal)
		signal.Notify(signalChan)
		go func() {
			for {
				select {
				case s := <-signalChan:
					switch s {
					case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
						// 进程退出时，保存offset
						pprof.StopCPUProfile()
						os.Exit(0)
					}
				}
			}
		}()
	}

	// 打印版本信息
	if *v {
		fmt.Println(AgentInfo)
		return
	}
	// 初始化配置文件
	err = conf.New(*c)
	if err != nil {
		fmt.Printf("[ERROR] new config struct error : %v", err)
		return
	}
	// 初始化日志
	err = logger.New(conf.Cfg.Path)
	if err != nil {
		fmt.Printf("[ERROR] new logger error : %v", err)
		return
	}
	// 初始化redis,连接和健康检查
	//red := db.Redis{
	//	RedisAddr: conf.Cfg.RedAddr,
	//	RedisPass: conf.Cfg.RedPass,
	//	RedisDB:   conf.Cfg.DB,
	//}
	//red.Conn()
	// 初始化kafka配置
	//kaf := db.Kafka{
	//	Broker:  conf.Cfg.Broker,
	//	GroupID: conf.Cfg.GroupID,
	//	Topic:   conf.Cfg.Topic,
	//	Offset:  conf.Cfg.Offset,
	//}
	// 启动lua进程
	for i := 0; i < runtime.NumCPU(); i++ {
		go lua.LuaThread(i)
		//go kaf.Consumer(lua.Kchan, i)
	}
	// 本地模拟消费者，不使用kafka
	lua.TestConsumerRaw()
	// redis健康检查卡住主进程，redis异常断开程序终止
	//red.Health()
}
