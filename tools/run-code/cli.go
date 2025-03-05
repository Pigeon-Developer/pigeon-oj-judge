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
)

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

	writeFile("/tmp/pj-run-code/test.in", "")

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
		image := imageName + ":0.0.0-alpha.2"
		fmt.Printf("use %s test %s\n", image, v.Name())
		compileResult := actuator.RunInDocker(image, []string{"bash", "-l", "-c", "/app/build.sh"}, []mount.Mount{
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

		fmt.Printf("%s - %+v\n", v.Name(), compileResult)
	}
}

// 测试资源消耗的收集
func testResourceCollect() {

}

func main() {
	testRumtimeErrorCode()
	testResourceCollect()
}

// docker pull pigeonojdev/runtime-cpp:0.0.0-alpha.2
// docker pull pigeonojdev/runtime-c:0.0.0-alpha.2
// docker pull pigeonojdev/runtime-python:0.0.0-alpha.2
