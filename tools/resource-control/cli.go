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

func testTs() {
	os.RemoveAll("/tmp/pj-run-code/source-code")
	os.RemoveAll("/tmp/pj-run-code/artifacts")
	os.MkdirAll("/tmp/pj-run-code/source-code", os.ModePerm)
	os.MkdirAll("/tmp/pj-run-code/artifacts", os.ModePerm)

	writeFile("/tmp/pj-run-code/test.in", "")
	writeFile(path.Join("/tmp/pj-run-code/source-code", "user_code"), `
console.log("Hello, World!");
`)

	compileResult := actuator.RunInDocker("pigeonojdev/runtime-javascript:0.0.0-alpha.3", []string{"bash", "-l", "-c", "cd /mount/source-code && cp user_code user_code.js && node ./user_code.js"}, []mount.Mount{
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

func testC() {
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

func testOC() {
	os.RemoveAll("/tmp/pj-run-code/source-code")
	os.RemoveAll("/tmp/pj-run-code/artifacts")
	os.MkdirAll("/tmp/pj-run-code/source-code", os.ModePerm)
	os.MkdirAll("/tmp/pj-run-code/artifacts", os.ModePerm)

	writeFile("/tmp/pj-run-code/test.in", "")
	writeFile(path.Join("/tmp/pj-run-code/source-code", "user_code"), `
#import <Foundation/Foundation.h>

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        char input[100];
        while (fgets(input, 100, stdin) != NULL) {
            NSString *inputString = [NSString stringWithUTF8String:input];
            NSArray *numbers = [inputString componentsSeparatedByString:@" "];

            int a = [numbers[0] intValue];
            int b = [numbers[1] intValue];

            int sum = a + b;
            NSLog(@"%d", sum);
        }
    }
    return 0;
}
`)

	// image := "pigeonojdev/runtime-c:0.0.0-alpha.2"
	image := "rt-oc"
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

func testGolang() {
	os.RemoveAll("/tmp/pj-run-code/source-code")
	os.RemoveAll("/tmp/pj-run-code/artifacts")
	os.MkdirAll("/tmp/pj-run-code/source-code", os.ModePerm)
	os.MkdirAll("/tmp/pj-run-code/artifacts", os.ModePerm)

	writeFile("/tmp/pj-run-code/test.in", "")
	writeFile(path.Join("/tmp/pj-run-code/source-code", "user_code"), `
package main

import "fmt"

func main() {
	fmt.Println("Hello World!")
}
`)

	image := "pigeonojdev/runtime-golang:0.0.0-alpha.7"
	compileResult := actuator.RunInDocker(image, []string{"bash", "-c", "/app/build.sh"}, []mount.Mount{
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

func main() {
	testGolang()
}
