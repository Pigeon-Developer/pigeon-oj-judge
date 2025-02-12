package app

import "github.com/Pigeon-Developer/pigeon-oj-judge/solution"

type AppConfig struct {
	// 提交的数据来源
	SolutionSource solution.SourceConfig `json:"solution_source"`
}
