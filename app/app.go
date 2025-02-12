package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Pigeon-Developer/pigeon-oj-judge/solution"
)

var (
	configPath string
)

func Boot() {
	flag.StringVar(&configPath, "config", "./config.json", "config file path")

	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println("读取配置文件失败", err)
	}
	byteValue, _ := io.ReadAll(configFile)

	var appConfig AppConfig

	json.Unmarshal(byteValue, &appConfig)

	solution.NewSolutionPool(appConfig.SolutionSource)

	RunLoop()
}
