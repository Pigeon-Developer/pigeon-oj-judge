package solution

// 还没想好数据要有啥
type AnyConfig struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type SourceConfig struct {
	// 定义怎么获取题目数据 这个版本不做实现 默认读取 hustoj 的本地目录
	ProblemProvider AnyConfig `json:"problem_provider"`
	// 定义怎么获取用户的提交
	ReaderAndWriterConfig AnyConfig `json:"rw_config"`
}

type SolutionInstance struct {
	ID     int
	Source SolutionSource
}

var InstancePool map[int]SolutionInstance
var InstanceNextID = 1

func createSolutionSource(config AnyConfig) SolutionSource {
	ret, err := NewSolutionSourceDB(config.Data["dbType"].(string), config.Data["dsn"].(string))

	if err != nil {
		panic(err)
	}

	return ret
}

// 在内部创建并维护 reader/wiriter 列表
func createSolutionInstance(config SourceConfig) {
	instance := SolutionInstance{
		ID:     InstanceNextID,
		Source: createSolutionSource(config.ReaderAndWriterConfig),
	}

	InstancePool[InstanceNextID] = instance

	InstanceNextID++
}

func GetSolutionInstance(id int) *SolutionInstance {
	ins := InstancePool[id]

	return &ins
}

func NewSolutionPool(config SourceConfig) {
	InstancePool = make(map[int]SolutionInstance)
	createSolutionInstance(config)
}
