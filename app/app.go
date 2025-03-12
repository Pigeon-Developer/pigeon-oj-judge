package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Pigeon-Developer/pigeon-oj-judge/actuator"
	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
)

var (
	configPath string
)

func Boot() {
	os.MkdirAll("/etc/pigeon-oj-judge", 0755)

	flag.StringVar(&configPath, "config", "./config.json", "config file path")

	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println("读取配置文件失败", err)
		return
	}
	byteValue, _ := io.ReadAll(configFile)

	var appConfig AppConfig

	json.Unmarshal(byteValue, &appConfig)

	if appConfig.Docker.Wait > 0 {
		fmt.Println("等待其他容器启动中...")
		time.Sleep(time.Duration(appConfig.Docker.Wait) * time.Second)
	}

	actuator.Prepare(appConfig.BuiltinRuntime.EnableList, appConfig.BuiltinRuntime.Version)

	solution.NewSolutionPool(appConfig.SolutionSource)

	RunLoop(appConfig.Judge.MaxConcurrent, appConfig.Judge.EmptyWait)
}
