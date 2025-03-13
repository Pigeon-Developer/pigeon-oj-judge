package solution

// 数据库字段
type SolutionRecord struct {
	SolutionId int `json:"solution_id" db:"solution_id"`
	ProblemId  int `json:"problem_id" db:"problem_id"`
	ContestId  int `json:"contest_id" db:"contest_id"`
	Language   int `json:"language" db:"language"`
	Result     int `json:"result" db:"result"`
}

type ProblemRecord struct {
	ProblemId   int     `json:"problem_id" db:"problem_id"`
	TimeLimit   float64 `json:"time_limit" db:"time_limit"`
	MemoryLimit int     `json:"memory_limit" db:"memory_limit"`
}

type SourceCodeRecord struct {
	SolutionId int    `json:"solution_id" db:"solution_id"`
	Source     string `json:"source" db:"source"`
}

const (
	// 等待
	Result_PD = 0
	// 等待重判
	Result_PR = 1
	// 编译中
	Result_CI = 2
	// 运行并评判
	Result_RJ = 3
	// 正确
	Result_AC = 4

	// 格式错误
	Result_PE = 5
	// 答案错误
	Result_WA = 6
	// 时间超限
	Result_TLE = 7
	// 内存超限
	Result_MLE = 8
	// 输出超限
	Result_OLE = 9

	// 运行时错误
	Result_RE = 10
	// 编译错误
	Result_CE = 11
	// 编译成功
	Result_CO = 12
	// 测试运行
	Result_TR = 13
	// 待裁判确认
	Result_MC = 14

	// 远程提交中
	Result_REMOTE_SUBMITTING = 15
	// 远程等待
	Result_REMOTE_RP = 16
	// 远程判题中
	Result_REMOTE_RJ = 17
)

type Solution struct {
	SolutionId  int
	ProblemId   int
	ContestId   int
	Language    int
	TimeLimit   float64
	MemoryLimit int
	Code        string
}

type SolutionResult struct {
	// 最终的结果
	Result int
}

// 获取 solution 的接口
type SolutionSource interface {
	// 释放对应资源
	Close()
	// 获取一个待判题的提交，如果没有数据，则立即返回
	// 如果获取到了数据，在函数返回前，就会修改提交状态到判题中
	GetOne(languageList []int) (*Solution, error)
	// 更新用户的 solution 状态
	Update(solutionId int, result SolutionResult) error
}

type JudgeJob struct {
	// Solution instance ID，用于标识来源
	SourceID int
	Data     *Solution
}
