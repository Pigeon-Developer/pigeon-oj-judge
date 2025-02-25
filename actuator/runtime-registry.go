package actuator

type ImageConfig struct {
	BuildCmd string `json:"build_cmd"`
	RunCmd   string `json:"run_cmd"`
	Image    string `json:"image"`
}

var (
	RuntimeRegistry = make(map[int]ImageConfig)
)

func init() {
	RuntimeRegistry[Language_python] = ImageConfig{
		BuildCmd: "/app/build.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeon-oj/runtime-python:DATE-TAG",
	}
	RuntimeRegistry[Language_c] = ImageConfig{
		BuildCmd: "/app/build-c.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeon-oj/runtime-cpp:DATE-TAG",
	}
	RuntimeRegistry[Language_cpp] = ImageConfig{
		BuildCmd: "/app/build-cpp.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeon-oj/runtime-python:DATE-TAG",
	}
	RuntimeRegistry[Language_java] = ImageConfig{
		BuildCmd: "/app/build.sh",
		RunCmd:   "/app/run.sh",
		Image:    "pigeon-oj/runtime-java:DATE-TAG",
	}
}
