package rule

import (
	"FBI-Analyzer/logger"
	"io/ioutil"
	"os"
	"time"

	"github.com/yuin/gopher-lua"
)

// 同步规则配置
func SyncRule(L *lua.LState, luaFuncs map[string]*lua.LFunction) {

	var scriptsPath = "scripts/"
	var ruleFileInfo = make(map[string]int64)
	var fileList []os.FileInfo
	var file os.FileInfo
	var filename string
	var fileModTime int64
	var err error
	var ok bool
	s3 := time.NewTicker(3 * time.Second)
	defer s3.Stop()

	for {
		select {
		case <-s3.C:
			fileList, err = ioutil.ReadDir(scriptsPath)
			if err != nil {
				logger.Print(logger.ERROR, "get scripts info failed. %v", err)
			}
			for _, file = range fileList {
				filename = file.Name()
				fileModTime = file.ModTime().Unix()
				if _, ok = ruleFileInfo[filename]; ok {
					if ruleFileInfo[filename] != fileModTime {
						ruleFileInfo[filename] = fileModTime
						compileLua(L, scriptsPath+filename, luaFuncs)
						logger.Print(logger.INFO, "compile modified lua script %s.", filename)
					}
				} else {
					ruleFileInfo[filename] = fileModTime
					compileLua(L, scriptsPath+filename, luaFuncs)
					logger.Print(logger.INFO, "compile new lua script %s.", filename)
				}
			}
		}
	}
}
