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
	Language_c          = 0
	Language_cpp        = 1
	Language_pascal     = 2
	Language_java       = 3
	Language_ruby       = 4
	Language_bash       = 5
	Language_python     = 6
	Language_php        = 7
	Language_perl       = 8
	Language_csharp     = 9
	Language_objectivec = 10
	Language_freebasic  = 11
	Language_scheme     = 12
	Language_clang      = 13
	Language_clangpp    = 14
	Language_lua        = 15
	Language_javascript = 16
	Language_golang     = 17
	Language_sql        = 18
	Language_fortran    = 19
	Language_matlab     = 20
	Language_cobol      = 21
	Language_r          = 22
	Language_scratch3   = 23
	Language_cangjie    = 24
)

const (
	SIGXCPU = 152
)
