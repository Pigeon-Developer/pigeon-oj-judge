package app

import "github.com/Pigeon-Developer/pigeon-oj-judge/solution"

type JudgeConfig struct {
	// 最大判题并发
	MaxConcurrent int `json:"max_concurrent"`
}

type AppConfig struct {
	// 提交的数据来源
	SolutionSource solution.SourceConfig `json:"solution_source"`
	// 启用哪些内置语言镜像
	BuiltinRuntime map[string]bool `json:"builtin_runtime"`
	Judge          JudgeConfig     `json:"judge"`
}
