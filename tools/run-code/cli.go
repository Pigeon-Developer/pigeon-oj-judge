package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Pigeon-Developer/pigeon-oj-judge/actuator"
	"github.com/docker/docker/api/types/mount"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
)

const RuntimeImageTag = ":0.0.0-alpha.3"

func writeFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(content)
}

func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 确保写入完成
	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func buildAndRun(beginDesc, image, buildCmd string) {
	if buildCmd == "" {
		buildCmd = "/app/build.sh"
	}
	writeFile("/tmp/pj-run-code/test.out", "")
	fmt.Printf(beginDesc)

	compileResult := actuator.RunInDocker(image, []string{"bash", "-l", "-c", buildCmd}, []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: "/tmp/pj-run-code/source-code",
			Target: "/mount/source-code",
		},
		{
			Type:   mount.TypeBind,
			Source: "/tmp/pj-run-code/artifacts",
			Target: "/mount/artifacts",
		},
	}, 10)

	if compileResult.ExitCode != 0 {
		printResult(compileResult, actuator.RunResult{})
		return
	}

	runResult := actuator.RunInDocker(image, []string{"bash", "-l", "-c", "/app/run.sh"}, []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: "/tmp/pj-run-code/test.in",
			Target: "/app/data.in",
		},
		{
			Type:   mount.TypeBind,
			Source: "/tmp/pj-run-code/test.out",
			Target: "/app/data.out",
		},
		{
			Type:   mount.TypeBind,
			Source: "/tmp/pj-run-code/artifacts",
			Target: "/mount/artifacts",
		},
	}, 10)

	printResult(compileResult, runResult)
}

// 测试会引起 runtime 异常的代码
func testRumtimeErrorCode() {
	os.RemoveAll("/tmp/pj-run-code")
	os.MkdirAll("/tmp/pj-run-code", os.ModePerm)
	cmd := exec.Command("bash", "-c", "cd /tmp/pj-run-code && git clone https://github.com/vijos/malicious-code.git")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return
	}

	writeFile("/tmp/pj-run-code/test.in", "1 2")

	entries, _ := os.ReadDir("/tmp/pj-run-code/malicious-code")
	for _, v := range entries {
		os.RemoveAll("/tmp/pj-run-code/source-code")
		os.RemoveAll("/tmp/pj-run-code/artifacts")
		os.MkdirAll("/tmp/pj-run-code/source-code", os.ModePerm)
		os.MkdirAll("/tmp/pj-run-code/artifacts", os.ModePerm)

		imageName := ""
		if strings.HasSuffix(v.Name(), ".cpp") {
			imageName = "pigeonojdev/runtime-cpp"
		} else if strings.HasSuffix(v.Name(), ".c") {
			imageName = "pigeonojdev/runtime-c"
		} else if strings.HasSuffix(v.Name(), ".py") {
			imageName = "pigeonojdev/runtime-python"
		} else {
			continue
		}
		// if strings.HasSuffix(v.Name(), ".cpp") {
		// 	imageName = "rt-cpp"
		// } else if strings.HasSuffix(v.Name(), ".c") {
		// 	imageName = "rt-c"
		// } else if strings.HasSuffix(v.Name(), ".py") {
		// 	imageName = "rt-py"
		// } else {
		// 	continue
		// }

		copyFile(path.Join("/tmp/pj-run-code/malicious-code", v.Name()), path.Join("/tmp/pj-run-code/source-code", "user_code"))

		// image := imageName
		image := imageName + RuntimeImageTag

		buildAndRun(fmt.Sprintf("use %s test %s\n", image, v.Name()), image, "")
	}
}

// 测试资源消耗的收集
func testResourceCollect() {
	os.RemoveAll("/tmp/pj-run-code")
	os.MkdirAll("/tmp/pj-run-code", os.ModePerm)
	cmd := exec.Command("bash", "-c", "cd /tmp/pj-run-code && git clone https://github.com/Pigeon-Developer/language-test-code.git")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return
	}

	writeFile("/tmp/pj-run-code/test.in", "1 2")

	// 先来测试基础代码
	codeBase := "/tmp/pj-run-code/language-test-code"
	langs, _ := os.ReadDir(codeBase + "/runtime")
	for _, lang := range langs {
		_, hasLang := actuator.LangMap[lang.Name()]
		isClang := strings.HasPrefix(lang.Name(), "clang-")

		if !isClang && !hasLang {
			continue
		}

		files, _ := os.ReadDir(codeBase + "/runtime/" + lang.Name())
		for _, file := range files {
			os.RemoveAll("/tmp/pj-run-code/source-code")
			os.RemoveAll("/tmp/pj-run-code/artifacts")
			os.MkdirAll("/tmp/pj-run-code/source-code", os.ModePerm)
			os.MkdirAll("/tmp/pj-run-code/artifacts", os.ModePerm)
			writeFile("/tmp/pj-run-code/test.out", "")

			imageName := "pigeonojdev/runtime-" + lang.Name()
			if isClang {
				imageName = "pigeonojdev/runtime-clang"
			}

			copyFile(path.Join(codeBase, "runtime", lang.Name(), file.Name()), path.Join("/tmp/pj-run-code/source-code", "user_code"))

			// image := imageName
			image := imageName + RuntimeImageTag

			buildCmd := ""

			if isClang {
				strs := strings.Split(lang.Name(), "-")
				buildCmd = "/app/build-" + strs[1] + ".sh"
			}

			buildAndRun(fmt.Sprintf("use %s test %s/%s\n", image, lang.Name(), file.Name()), image, buildCmd)
		}
	}

	// 再来测试异常代码
	resources, _ := os.ReadDir(codeBase + "/resource")
	for _, file := range resources {
		os.RemoveAll("/tmp/pj-run-code/source-code")
		os.RemoveAll("/tmp/pj-run-code/artifacts")
		os.MkdirAll("/tmp/pj-run-code/source-code", os.ModePerm)
		os.MkdirAll("/tmp/pj-run-code/artifacts", os.ModePerm)
		writeFile("/tmp/pj-run-code/test.out", "")

		imageName := "pigeonojdev/runtime-c"

		copyFile(path.Join(codeBase, "resource", file.Name()), path.Join("/tmp/pj-run-code/source-code", "user_code"))

		// image := imageName
		image := imageName + RuntimeImageTag

		buildAndRun(fmt.Sprintf("use %s test %s/%s\n", image, "resource", file.Name()), image, "")
	}
}

func printResult(compile, run actuator.RunResult) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"", "compile", "run"})
	t.AppendRows([]table.Row{
		{"exit code", compile.ExitCode, run.ExitCode},
		{"timecost/ms", compile.TimeCost, run.TimeCost},
		{"mmory usage", humanize.Bytes(uint64(compile.MemoryUsage)), humanize.Bytes(uint64(run.MemoryUsage))},
		{"stdout", compile.Stdout, run.Stdout},
		{"stderr", compile.Stderr, run.Stderr},
	})
	t.Render()
}

func pull(image string) {
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	pull("pigeonojdev/runtime-c" + RuntimeImageTag)
	pull("pigeonojdev/runtime-cpp" + RuntimeImageTag)
	pull("pigeonojdev/runtime-python" + RuntimeImageTag)

	testRumtimeErrorCode()
	testResourceCollect()
}
