package app

import "github.com/Pigeon-Developer/pigeon-oj-judge/solution"

type AppConfig struct {
	// 提交的数据来源
	SolutionSource solution.SourceConfig `json:"solution_source"`
	// 启用哪些内置语言镜像
	BuiltinRuntime map[string]bool `json:"builtin_runtime"`
}
