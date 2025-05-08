## 安装方式
可以从[Release](https://github.com/loveRyujin/ReviewBot/Release)下载预编译的二进制文件，将其放置路径加入环境变量。
执行以下命令：
```sh
reviewbot
```  

输出以下内容，代表安装成功：
```
help code review when merging code

Usage:
  reviewbot [command]

Available Commands:
  commit      Automically generate commit message
  config      Manage configuration settings
  help        Help about any command
  review      Auto review code changes in git stage

Flags:
  -c, --config string   config file path
  -h, --help            help for reviewbot

Use "reviewbot [command] --help" for more information about a command.
```  

从源码安装：
```sh
go install github.com/loveRyujin/ReviewBot/cmd/reviewbot
```

## 功能
- 帮助生成git commit message（遵循conventional commits规范）
- 帮助进行code review，针对代码变更生成对应的建议
- 支持自定义git diff生成的差异上下文行数
- 支持选择让git diff忽略的文件
- 支持proxy配置

## 使用方法
### 生成git commit message
```sh
git add .
reviewbot commit
```
  
### 进行code review
```sh
git add .
reviewbot review
```

### 列出可选配置
```sh
reviewbot config list
```
可选配置如下:  

![config_list](./images/config_list.png)


### 更新配置
```sh
reviewbot config set ai.api_key xxxxxx
```
更新成功输出类似下面：
![config_set](./images/config_set.png)
