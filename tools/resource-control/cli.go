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
	code := `
console.log("Hello, World!");
`
	runAndDump("pigeonojdev/runtime-javascript:0.0.0-alpha.8", code)
}

func testC() {
	code := `
	#include <stdio.h>
	int main()
	{
	   printf("Hello, World!");
	   return 0;
	}
	`

	runAndDump("pigeonojdev/runtime-c:0.0.0-alpha.8", code)
}

func testOC() {
	code := `
#import <Foundation/Foundation.h>

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSFileHandle *input = [NSFileHandle fileHandleWithStandardInput];
        NSData *inputData = [input availableData];
        NSString *inputString = [[NSString alloc] initWithData:inputData encoding:NSUTF8StringEncoding];


        inputString = [inputString stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceAndNewlineCharacterSet]];

  
        NSArray *components = [inputString componentsSeparatedByString:@" "];
        if ([components count] != 2) {
            return 1;
        }
        NSNumberFormatter *formatter = [[NSNumberFormatter alloc] init];
        NSNumber *number1 = [formatter numberFromString:[components objectAtIndex:0]];
        NSNumber *number2 = [formatter numberFromString:[components objectAtIndex:1]];

        if (number1 == nil || number2 == nil) {
  
            return 1;
        }

        int sum = [number1 intValue] + [number2 intValue];

        printf("%d", sum);
    }
    return 0;
}
`

	image := "pigeonojdev/runtime-objectivec:0.0.0-alpha.8"
	runAndDump(image, code)

}

func testGolang() {
	code := `
package main

import "fmt"

func main() {
	fmt.Println("Hello World!")
}
	`

	image := "pigeonojdev/runtime-golang:0.0.0-alpha.8"
	runAndDump(image, code)
}

func testFreeBasic() {
	runAndDump("pigeonojdev/runtime-freebasic:0.0.0-alpha.8", `
DIM a AS INTEGER
DIM b AS INTEGER
Dim sum As Integer

INPUT "", a, b
sum = a + b
Print "" & sum
`)
}

func testBash() {
	code := `#!/bin/bash

while read -r line; do
    a=$(echo $line | cut -d ' ' -f 1)
    b=$(echo $line | cut -d ' ' -f 2)
    sum=$((a + b))
    echo $sum
done

`
	runAndDump("pigeonojdev/runtime-bash:0.0.0-alpha.8", code)
}

func testScheme() {
	code := `
(use-modules (ice-9 rdelim))  ; Guile Scheme

(define (read-numbers)
  (let* ((line (read-line))
         (numbers (string-split line #\space)))
    (map string->number numbers)))

(define (string-split str ch)
  (let loop ((str str) (result '()))
    (let ((pos (string-index str ch)))
      (if pos
          (loop (substring str (+ pos 1) (string-length str))
                (cons (substring str 0 pos) result))
          (reverse (cons str result))))))

(let ((numbers (read-numbers)))
  (if (>= (length numbers) 2)
      (display (+ (car numbers) (cadr numbers)))
      (display "Error: Need two numbers"))
  (newline))

`
	runAndDump("pigeonojdev/runtime-scheme:0.0.0-alpha.8", code)
}

func runAndDump(image, code string) {
	os.RemoveAll("/tmp/pj-run-code/source-code")
	os.RemoveAll("/tmp/pj-run-code/source-code")
	os.RemoveAll("/tmp/pj-run-code/artifacts")
	os.MkdirAll("/tmp/pj-run-code/source-code", os.ModePerm)
	os.MkdirAll("/tmp/pj-run-code/artifacts", os.ModePerm)

	writeFile("/tmp/pj-run-code/test.in", "1111 2222\r\n")
	writeFile("/tmp/pj-run-code/test.out", "")
	writeFile(path.Join("/tmp/pj-run-code/source-code", "user_code"), code)

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

	runResult := actuator.RunInDocker(image, []string{"bash", "-c", "/app/run.sh"}, []mount.Mount{
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

	fmt.Printf("runResult: %+v\n", runResult)

	outfileBytes, err := os.ReadFile("/tmp/pj-run-code/test.out") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	outfileStr := string(outfileBytes)

	fmt.Printf("outfile byte: [%+v] str: [%s]\n", outfileBytes, outfileStr)
}

func testR() {
	code := `
input <- readLines("stdin", n=1)
nums <- as.numeric(strsplit(input, " ")[[1]])
cat(nums[1] + nums[2])
`

	runAndDump("pigeonojdev/runtime-r:0.0.0-alpha.8", code)
}

func main() {
	testR()
}
