package actuator

type UserCodeRunResult struct {
	// 用户进程的退出码
	ExitCode int `json:"exit_code"`
	// 内存使用
	MemoryUsage int `json:"memory_usage"`
	// 耗时
	TimeCost int `json:"time_cost"`
	// 是否与正确答案一致
	Match int `json:"match"`

	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

const (
	SIGXCPU = 152
)
