package main

import (
	"fmt"
	"strings"

	"github.com/Pigeon-Developer/pigeon-oj-judge/actuator"
)

const Prefix = `
| image       | 体积                      |
| ----------- | ------------------------- |`

func renderTemplate(language string) string {
	_template := `| %s          | ![runtime-%s-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-%s/0.0.0-alpha.6)  |`
	return strings.Trim(fmt.Sprintf(_template, language, language, language), "\n")
}

func main() {
	result := make([]string, 0, 32)
	result = append(result, Prefix)

	for _, lang := range actuator.SimpleLangList {
		language := actuator.LanguageMap[lang]
		result = append(result, renderTemplate(language))
	}

	result = append(result, renderTemplate("clang"))

	fmt.Println(strings.Join(result, "\n"))
}
