package main

import (
	"fmt"
	"strings"

	"github.com/Pigeon-Developer/pigeon-oj-judge/actuator"
)

const Prefix = `name: release-runtime

on:
  workflow_dispatch:
    inputs:
      version:
        description: code tag
        required: true

env:
  cr_coding: g-yuie8424-docker.pkg.coding.net
  cr_acr: registry.cn-hangzhou.aliyuncs.com
  cr_tcr: ccr.ccs.tencentyun.com

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Login to Coding Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.cr_coding }}
          username: ${{ secrets.CODING_USERNAME }}
          password: ${{ secrets.CODING_TOKEN }}
      - name: Login to acr
        uses: docker/login-action@v3
        with:
          registry: ${{ env.cr_acr }}
          username: ${{ secrets.ACR_USERNAME }}
          password: ${{ secrets.CR_TOKEN }}
      - name: Login to tcr
        uses: docker/login-action@v3
        with:
          registry: ${{ env.cr_tcr }}
          username: ${{ secrets.TCR_USERNAME }}
          password: ${{ secrets.CR_TOKEN }}`

func renderTemplate(language string) string {
	// 	_template := `
	//       - name: Build and push - %s
	//         uses: docker/build-push-action@v6
	//         with:
	//           file: ./runtime/%s/dockerfile
	//           context: ./runtime/%s
	//           push: true
	//           tags: |
	//             pigeonojdev/runtime-%s:${{ github.event.inputs.version }}
	//             ${{ env.cr_coding }}/pigeon-oj/release/runtime-%s:${{ github.event.inputs.version }}
	//             ${{ env.cr_acr }}/pigeon-oj/runtime-%s:${{ github.event.inputs.version }}
	//             ${{ env.cr_tcr }}/pigeon-oj/runtime-%s:${{ github.event.inputs.version }}
	// `
	_template := `
      - name: Build and push - %s
        uses: docker/build-push-action@v6
        with:
          file: ./runtime/%s/dockerfile
          context: ./runtime/%s
          push: true
          tags: |
            pigeonojdev/runtime-%s:${{ github.event.inputs.version }}
`
	return strings.Trim(fmt.Sprintf(_template, language, language, language, language), "\n")
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
