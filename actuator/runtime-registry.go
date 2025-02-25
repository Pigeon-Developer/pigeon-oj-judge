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

var (
	RuntimeRegistry = make(map[int]ImageConfig)
)

func init() {
	CurrentTag := "0.0.0-alpha.0"
	RuntimeRegistry[Language_python] = ImageConfig{
		BuildCmd: "/app/build.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-python:" + CurrentTag,
	}
	RuntimeRegistry[Language_c] = ImageConfig{
		BuildCmd: "/app/build-c.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-cpp:" + CurrentTag,
	}
	RuntimeRegistry[Language_cpp] = ImageConfig{
		BuildCmd: "/app/build-cpp.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-cpp:" + CurrentTag,
	}
	RuntimeRegistry[Language_java] = ImageConfig{
		BuildCmd: "/app/build.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeonojdev/runtime-java:" + CurrentTag,
	}
}

func pull(image string) {
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func PullBuiltinRuntime() {
	for _, v := range RuntimeRegistry {
		pull(v.Image)
	}
}
