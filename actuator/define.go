package actuator

type UserCodeRunResult struct {
	// 用户进程的退出码
	ExitCode int `json:"exit_code"`
	// 内存使用
	MemoryUsage int `json:"memory_usage"`
	// 耗时
	TimeCost int `json:"time_cost"`
	// 是否与正确答案存在不同, -1 未知, 0 一致 1 不一致
	// 其余值为将来判断数据一致但格式不一致预留
	Match int `json:"match"`

	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

const (
	SIGXCPU = 152
)
