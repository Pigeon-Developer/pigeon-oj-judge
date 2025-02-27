[WIP] 一个简易的 docker 判题机实现

- solution 是对用户提交操作的抽象，提供读写用户提交以及题目信息的操作
- actuator 编译并执行用户代码，获取用户的输出，对比用户的输出是否与标准输出一致，返回对比结果
- app 中为控制功能，初始化数据库连接，多久轮询一次是否存在提交，对比结果如何回写等都在这里面实现
- runtime 中为每个语言的编译/运行 docker 定义

todo list

- [ ] 为现在 hustoj 在 debian12 下支持的每个语言构建镜像
  - [ ] sql
- [ ] 判题为每个用例记录 编译/执行/对比 的详细信息
- [ ] 支持 hustoj udp 判题
- [ ] 支持 hustoj redis 判题

## 差异

编译与运行的镜像一般采用下面两种方式构建

- 语言官方维护的 debian 镜像
- 尽可能的使用 debian bookworm 环境，安装 debian 源自带的对应语言的编译器

部分语言的编译与执行同 hustoj 存在一定的差异，其中差异较大的为

| 语言        | hustoj 处理方式           | pigeon-oj-judge 处理方式      |
| ----------- | ------------------------- | ----------------------------- |
| C#          | msc 编译，mono 执行       | dotnet cli 编译，直接执行产物 |
| Objective-C | gcc 编译，带 GNUstep 环境 | clang 编译                    |
