package main

import (
	"fmt"
	"strings"

	"github.com/Pigeon-Developer/pigeon-oj-judge/actuator"
)

func renderTemplate(language string) string {
	_template := `
      - name: Build and push - %s
        uses: docker/build-push-action@v6
        with:
          file: ./runtime/%s/dockerfile
          context: ./runtime/%s
          push: true
          tags: |
            pigeonojdev/runtime-%s:${{ github.event.inputs.version }}
            ${{ env.cr_coding }}/pigeon-oj/release/runtime-%s:${{ github.event.inputs.version }}
            ${{ env.cr_acr }}/pigeon-oj/runtime-%s:${{ github.event.inputs.version }}
            ${{ env.cr_tcr }}/pigeon-oj/runtime-%s:${{ github.event.inputs.version }}
`
	return strings.Trim(fmt.Sprintf(_template, language, language, language, language, language, language, language), "\n")
}

func main() {
	result := make([]string, 0, 32)
	for _, lang := range actuator.SimpleLangList {
		language := actuator.LanguageMap[lang]
		result = append(result, renderTemplate(language))
	}

	result = append(result, renderTemplate("clang"))

	fmt.Println(strings.Join(result, "\n"))
}
