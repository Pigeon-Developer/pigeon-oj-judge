[WIP] 一个简易的 docker 判题机实现

- solution 是对用户提交操作的抽象，提供读写用户提交以及题目信息的操作
- actuator 编译并执行用户代码，获取用户的输出，对比用户的输出是否与标准输出一致，返回对比结果
- app 中为控制功能，初始化数据库连接，多久轮询一次是否存在提交，对比结果如何回写等都在这里面实现
- runtime 中为每个语言的编译/运行 docker 定义

todo list

- [ ] 适配运行在容器中的情况，允许 docker compose 一键启动
- [ ] runtime 适配
  - [ ] sql
  - [ ] 支持自定义 runtime 配置
  - [ ] 允许对指定 cid/pid 配置 runtime
- [ ] 增加执行异常的信息
  - [ ] 构建/执行时 CPU/内存/IO 超过限制
- [ ] 资源限制
  - [ ] 适配 Cgroup Driver: systemd
  - [ ] 构建/执行时的 CPU/内存/IO 限制可以使用配置文件
  - [ ] 允许分语言配置资源限制
  - [ ] 限制用户代码只能跑在 1 个核心上，在不同核心中均匀分配
- [ ] hustoj-php 兼容
  - [ ] 支持 http 判题
  - [ ] 支持 udp 判题
  - [ ] 支持 redis 判题
  - [ ] 支持 spj

## 差异

编译与运行的镜像一般采用下面两种方式构建

- 语言官方维护的 debian 镜像
- 尽可能的使用 debian bookworm 环境，安装 debian 源自带的对应语言的编译器

部分语言的编译与执行同 hustoj 存在一定的差异，其中差异较大的为

| 语言        | hustoj 处理方式           | pigeon-oj-judge 处理方式      |
| ----------- | ------------------------- | ----------------------------- |
| C#          | msc 编译，mono 执行       | dotnet cli 编译，直接执行产物 |
| Objective-C | gcc 编译，带 GNUstep 环境 | clang 编译                    |

## 镜像体积

内置的判题为每种 hustoj 支持的语言单独打包了一个 image

| image      | 体积                                                                                                      |
| ---------- | --------------------------------------------------------------------------------------------------------- |
| c          | ![runtime-c-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-c)                   |
| cpp        | ![runtime-cpp-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-cpp)               |
| pascal     | ![runtime-pascal-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-pascal)         |
| java       | ![runtime-java-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-java)             |
| ruby       | ![runtime-ruby-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-ruby)             |
| bash       | ![runtime-bash-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-bash)             |
| python     | ![runtime-python-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-python)         |
| php        | ![runtime-php-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-php)               |
| perl       | ![runtime-perl-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-perl)             |
| csharp     | ![runtime-csharp-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-csharp)         |
| objectivec | ![runtime-objectivec-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-objectivec) |
| freebasic  | ![runtime-freebasic-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-freebasic)   |
| scheme     | ![runtime-scheme-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-scheme)         |
| lua        | ![runtime-lua-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-lua)               |
| javascript | ![runtime-javascript-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-javascript) |
| golang     | ![runtime-golang-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-golang)         |
| fortran    | ![runtime-fortran-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-fortran)       |
| matlab     | ![runtime-matlab-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-matlab)         |
| cobol      | ![runtime-cobol-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-cobol)           |
| r          | ![runtime-r-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-r)                   |
| scratch3   | ![runtime-scratch3-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-scratch3)     |
| cangjie    | ![runtime-cangjie-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-cangjie)       |
| clang      | ![runtime-clang-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-clang)           |

## 判题开销

以下为 a+b 的判题开销

- [ ] C# `File size limit exceeded(core dumped) dotnet build --property:Configuration=Release -o /mount/artifacts`
- [ ] freebasic `ld: cannot open output file build_result/main.bin: No such file or directory`
- [ ] golang `go: command not found`
- [ ] Objective-C a+b 输出不对
- [ ] R a+b 输出不对
- [ ] scheme `Unbound variable: read-line`
- [ ] scratch3 `scratch-run: command not found`
