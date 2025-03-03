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
	Language_c          = 0b0000000000000000000000001
	Language_cpp        = 0b0000000000000000000000010
	Language_pascal     = 0b0000000000000000000000100
	Language_java       = 0b0000000000000000000001000
	Language_ruby       = 0b0000000000000000000010000
	Language_bash       = 0b0000000000000000000100000
	Language_python     = 0b0000000000000000001000000
	Language_php        = 0b0000000000000000010000000
	Language_perl       = 0b0000000000000000100000000
	Language_csharp     = 0b0000000000000001000000000
	Language_objectivec = 0b0000000000000010000000000
	Language_freebasic  = 0b0000000000000100000000000
	Language_scheme     = 0b0000000000001000000000000
	Language_clang      = 0b0000000000010000000000000
	Language_clangpp    = 0b0000000000100000000000000
	Language_lua        = 0b0000000001000000000000000
	Language_javascript = 0b0000000010000000000000000
	Language_golang     = 0b0000000100000000000000000
	Language_sql        = 0b0000001000000000000000000
	Language_fortran    = 0b0000010000000000000000000
	Language_matlab     = 0b0000100000000000000000000
	Language_cobol      = 0b0001000000000000000000000
	Language_r          = 0b0010000000000000000000000
	Language_scratch3   = 0b0100000000000000000000000
	Language_cangjie    = 0b1000000000000000000000000
)
