package actuator

import (
	"os"
	"os/exec"
)

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
	LanguageMap     = make(map[int]string)
	LangMap         = make(map[string]int)
)

func init() {
	// 1
	LanguageMap[Language_c] = "c"
	LanguageMap[Language_cpp] = "cpp"
	LanguageMap[Language_pascal] = "pascal"
	LanguageMap[Language_java] = "java"
	LanguageMap[Language_ruby] = "ruby"
	// 2
	LanguageMap[Language_bash] = "bash"
	LanguageMap[Language_python] = "python"
	LanguageMap[Language_php] = "php"
	LanguageMap[Language_perl] = "perl"
	LanguageMap[Language_csharp] = "csharp"
	// 3
	LanguageMap[Language_objectivec] = "objectivec"
	LanguageMap[Language_freebasic] = "freebasic"
	LanguageMap[Language_scheme] = "scheme"
	LanguageMap[Language_clang] = "clang"
	LanguageMap[Language_clangpp] = "clangpp"
	// 4
	LanguageMap[Language_lua] = "lua"
	LanguageMap[Language_javascript] = "javascript"
	LanguageMap[Language_golang] = "golang"
	LanguageMap[Language_sql] = "sql"
	LanguageMap[Language_fortran] = "fortran"
	// 5
	LanguageMap[Language_matlab] = "matlab"
	LanguageMap[Language_cobol] = "cobol"
	LanguageMap[Language_r] = "r"
	LanguageMap[Language_scratch3] = "scratch3"
	LanguageMap[Language_cangjie] = "cangjie"

	// 1
	LangMap["c"] = Language_c
	LangMap["cpp"] = Language_cpp
	LangMap["pascal"] = Language_pascal
	LangMap["java"] = Language_java
	LangMap["ruby"] = Language_ruby
	// 2
	LangMap["bash"] = Language_bash
	LangMap["python"] = Language_python
	LangMap["php"] = Language_php
	LangMap["perl"] = Language_perl
	LangMap["csharp"] = Language_csharp
	// 3
	LangMap["objectivec"] = Language_objectivec
	LangMap["freebasic"] = Language_freebasic
	LangMap["scheme"] = Language_scheme
	LangMap["clang"] = Language_clang
	LangMap["clangpp"] = Language_clangpp
	// 4
	LangMap["lua"] = Language_lua
	LangMap["javascript"] = Language_javascript
	LangMap["golang"] = Language_golang
	LangMap["sql"] = Language_sql
	LangMap["fortran"] = Language_fortran
	// 5
	LangMap["matlab"] = Language_matlab
	LangMap["cobol"] = Language_cobol
	LangMap["r"] = Language_r
	LangMap["scratch3"] = Language_scratch3
	LangMap["cangjie"] = Language_cangjie

	list := []int{
		// 1
		Language_c,
		Language_cpp,
		Language_pascal,
		Language_java,
		Language_ruby,
		// 2
		Language_bash,
		Language_python,
		Language_php,
		Language_perl,
		Language_csharp,
		// 3
		// objectivec
		Language_freebasic,
		Language_scheme,
		// clang
		// clangpp

		// 4
		Language_lua,
		Language_javascript,
		Language_golang,
		// sql
		Language_fortran,

		// 5
		Language_matlab,
		Language_cobol,
		Language_r,
		Language_scratch3,
		Language_cangjie,
	}

	CurrentTag := "0.0.0-alpha.0"

	for _, v := range list {
		RuntimeRegistry[v] = ImageConfig{
			BuildCmd: "/app/build.sh",
			RunCmd:   "/app/run.sh",
			Image:    "pigeonojdev/runtime-" + LanguageMap[v] + ":" + CurrentTag,
		}
	}

	// 下面三个语言共享一个镜像
	RuntimeRegistry[Language_clang] = ImageConfig{
		BuildCmd: "/app/build-c.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-clang:" + CurrentTag,
	}
	RuntimeRegistry[Language_clangpp] = ImageConfig{
		BuildCmd: "/app/build-cpp.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-clang:" + CurrentTag,
	}
	RuntimeRegistry[Language_objectivec] = ImageConfig{
		BuildCmd: "/app/build-objectivec.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-clang:" + CurrentTag,
	}
}

func pull(image string) {
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func PullBuiltinRuntime(builtinRuntime map[string]bool) {
	for language, enable := range builtinRuntime {
		if !enable {
			continue
		}
		langDefine, ok := LangMap[language]
		if !ok {
			continue
		}
		config, ok := RuntimeRegistry[langDefine]
		if !ok {
			continue
		}

		pull(config.Image)
	}
}
