package main

import (
	"fmt"
	"os"
	"path"

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

func main() {
	os.RemoveAll("/tmp/pj-run-code/source-code")
	os.RemoveAll("/tmp/pj-run-code/artifacts")
	os.MkdirAll("/tmp/pj-run-code/source-code", os.ModePerm)
	os.MkdirAll("/tmp/pj-run-code/artifacts", os.ModePerm)

	writeFile("/tmp/pj-run-code/test.in", "")
	writeFile(path.Join("/tmp/pj-run-code/source-code", "user_code"), `
#include <stdio.h>
int main()
{
   printf("Hello, World!");
   return 0;
}
`)

	// image := "pigeonojdev/runtime-c:0.0.0-alpha.2"
	image := "rt-c"
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

	fmt.Printf("compileResult: %+v\n", compileResult)
}
