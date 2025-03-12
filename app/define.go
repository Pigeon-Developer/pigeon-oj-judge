package app

import "github.com/Pigeon-Developer/pigeon-oj-judge/solution"

type JudgeConfig struct {
	// 最大判题并发
	MaxConcurrent int `json:"max_concurrent"`
	// 无任务时，等待时间
	EmptyWait int `json:"empty_wait"`
}

type BuiltinRuntimeConfig struct {
	EnableList []string `json:"enable_list"`
	Version    string   `json:"version"`
}

type DockerConfig struct {
	Wait int `json:"wait"`
}

type AppConfig struct {
	// 提交的数据来源
	SolutionSource solution.SourceConfig `json:"solution_source"`
	// 启用哪些内置语言镜像
	BuiltinRuntime BuiltinRuntimeConfig `json:"builtin_runtime"`
	Judge          JudgeConfig          `json:"judge"`
	Docker         DockerConfig         `json:"docker"`
}
