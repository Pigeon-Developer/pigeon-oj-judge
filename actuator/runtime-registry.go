package actuator

import "github.com/Pigeon-Developer/pigeon-oj-judge/types"

type ImageConfig struct {
	BuildCmd string `json:"build_cmd"`
	RunCmd   string `json:"run_cmd"`
	Image    string `json:"image"`
}

type BuiltinItem struct {
	Language string // 对应镜像后缀
	Lang     int    // hustoj 内部的编码
}

var (
	RuntimeRegistry = make(map[int]ImageConfig)
	SimpleLangList  = make([]int, 0, 32)
)

func initLanguageResource(CurrentTag string) {
	SimpleLangList = append(SimpleLangList,
		// 1
		types.Language_c,
		types.Language_cpp,
		types.Language_pascal,
		types.Language_java,
		types.Language_ruby,
		// 2
		types.Language_bash,
		types.Language_python,
		types.Language_php,
		types.Language_perl,
		types.Language_csharp,
		// 3
		types.Language_objectivec,
		types.Language_freebasic,
		types.Language_scheme,
		// clang
		// clangpp

		// 4
		types.Language_lua,
		types.Language_javascript,
		types.Language_golang,
		// sql
		types.Language_fortran,

		// 5
		types.Language_matlab,
		types.Language_cobol,
		types.Language_r,
		// types.Language_scratch3,
		types.Language_cangjie)

	for _, v := range SimpleLangList {
		RuntimeRegistry[v] = ImageConfig{
			BuildCmd: "/app/build.sh",
			RunCmd:   "/app/run.sh",
			Image:    "pigeonojdev/runtime-" + types.LanguageMap[v] + ":" + CurrentTag,
		}
	}

	// 下面两个语言共享一个镜像
	RuntimeRegistry[types.Language_clang] = ImageConfig{
		BuildCmd: "/app/build-c.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-clang:" + CurrentTag,
	}
	RuntimeRegistry[types.Language_clangpp] = ImageConfig{
		BuildCmd: "/app/build-cpp.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-clang:" + CurrentTag,
	}
}

func PullBuiltinRuntime(builtinRuntime []string) {
	for _, language := range builtinRuntime {
		langDefine, ok := types.LangMap[language]
		if !ok {
			continue
		}
		config, ok := RuntimeRegistry[langDefine]
		if !ok {
			continue
		}

		ImagePull(config.Image)
	}
}

// 初始化必要的资源
func Prepare(enableList []string, version string) {
	initLanguageResource(version)
	initDockerClient()

	PullBuiltinRuntime(enableList)
}
