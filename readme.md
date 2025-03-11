[WIP] 一个简易的 docker 判题机实现

- solution 是对用户提交操作的抽象，提供读写用户提交以及题目信息的操作
- actuator 编译并执行用户代码，获取用户的输出，对比用户的输出是否与标准输出一致，返回对比结果
- app 中为控制功能，初始化数据库连接，多久轮询一次是否存在提交，对比结果如何回写等都在这里面实现
- runtime 中为每个语言的编译/运行 docker 定义

todo list

- [ ] 适配运行在容器中的情况，允许 docker compose 一键启动
- [ ] runtime 适配
  - [ ] 支持 sql
  - [ ] 支持 scratch3
  - [ ] 支持自定义 runtime 配置
  - [ ] 允许对指定 cid/pid 配置 runtime
- [ ] 增加执行异常的信息
  - [ ] 构建/执行时 CPU/内存/IO 超过限制
- [ ] 资源限制
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

| image      | 体积                                                                                                                    |
| ---------- | ----------------------------------------------------------------------------------------------------------------------- |
| c          | ![runtime-c-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-c/0.0.0-alpha.8)                   |
| cpp        | ![runtime-cpp-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-cpp/0.0.0-alpha.8)               |
| pascal     | ![runtime-pascal-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-pascal/0.0.0-alpha.8)         |
| java       | ![runtime-java-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-java/0.0.0-alpha.8)             |
| ruby       | ![runtime-ruby-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-ruby/0.0.0-alpha.8)             |
| bash       | ![runtime-bash-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-bash/0.0.0-alpha.8)             |
| python     | ![runtime-python-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-python/0.0.0-alpha.8)         |
| php        | ![runtime-php-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-php/0.0.0-alpha.8)               |
| perl       | ![runtime-perl-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-perl/0.0.0-alpha.8)             |
| csharp     | ![runtime-csharp-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-csharp/0.0.0-alpha.8)         |
| objectivec | ![runtime-objectivec-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-objectivec/0.0.0-alpha.8) |
| freebasic  | ![runtime-freebasic-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-freebasic/0.0.0-alpha.8)   |
| scheme     | ![runtime-scheme-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-scheme/0.0.0-alpha.8)         |
| lua        | ![runtime-lua-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-lua/0.0.0-alpha.8)               |
| javascript | ![runtime-javascript-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-javascript/0.0.0-alpha.8) |
| golang     | ![runtime-golang-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-golang/0.0.0-alpha.8)         |
| fortran    | ![runtime-fortran-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-fortran/0.0.0-alpha.8)       |
| matlab     | ![runtime-matlab-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-matlab/0.0.0-alpha.8)         |
| cobol      | ![runtime-cobol-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-cobol/0.0.0-alpha.8)           |
| r          | ![runtime-r-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-r/0.0.0-alpha.8)                   |
| scratch3   | ![runtime-scratch3-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-scratch3/0.0.0-alpha.8)     |
| cangjie    | ![runtime-cangjie-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-cangjie/0.0.0-alpha.8)       |
| clang      | ![runtime-clang-image-size](https://img.shields.io/docker/image-size/pigeonojdev/runtime-clang/0.0.0-alpha.8)           |

## 判题开销

以下为 a+b 的判题开销

```
+------------+-----------------+---------------+
|            | COMPILE COST    | RUN COST      |
+------------+-----------------+---------------+
| bash       | 4.7 MB / 28ms   | 4.6 MB / 30ms |
| c          | 7.9 MB / 64ms   | 6.9 MB / 27ms |
| cangjie    | 187 MB / 1126ms | 6.6 MB / 39ms |
| cobol      | 26 MB / 143ms   | 7.2 MB / 30ms |
| cpp        | 53 MB / 316ms   | 7.0 MB / 29ms |
| csharp     | 240 MB / 3602ms | 5.8 MB / 60ms |
| fortran    | 15 MB / 77ms    | 7.1 MB / 29ms |
| freebasic  | 26 MB / 97ms    | 6.8 MB / 27ms |
| golang     | 341 MB / 8223ms | 8.8 MB / 28ms |
| java       | 64 MB / 939ms   | 27 MB / 114ms |
| javascript | 8.2 MB / 30ms   | 27 MB / 130ms |
| lua        | 7.0 MB / 29ms   | 4.5 MB / 30ms |
| matlab     | 5.2 MB / 29ms   | 30 MB / 164ms |
| objectivec | 59 MB / 313ms   | 10 MB / 53ms  |
| pascal     | 18 MB / 68ms    | 4.9 MB / 27ms |
| perl       | 6.9 MB / 32ms   | 4.9 MB / 31ms |
| php        | 15 MB / 49ms    | 8.8 MB / 46ms |
| python     | 22 MB / 186ms   | 6.9 MB / 46ms |
| r          | 4.8 MB / 33ms   | 75 MB / 180ms |
| ruby       | 7.3 MB / 38ms   | 9.2 MB / 85ms |
| scheme     | 9.0 MB / 28ms   | 19 MB / 133ms |
| clang-c    | 25 MB / 93ms    | 5.1 MB / 28ms |
| clang-cpp  | 41 MB / 373ms   | 7.0 MB / 29ms |
+------------+-----------------+---------------+
```
