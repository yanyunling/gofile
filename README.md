# 网页版简单文件服务器

## 介绍

- 超级管理员账号：admin
- 超级管理员默认密码：Admin123

## 文件权限说明

### 1.公开文件

1. 下载：所有用户均可
2. 上传/创建目录/删除：超级管理员

### 2.保护文件

1. 下载：登录用户均可
2. 上传/创建目录/删除：超级管理员

### 3.私有文件

1. 下载：登录用户自己，超级管理员
2. 上传/创建目录/删除：登录用户自己（如有更新权限），超级管理员

## 软件架构

- go1.25
- vue3
- sqlite

## 编译步骤

- 安装 go 环境
- 安装 nodejs 环境
- 进入目录，添加权限并执行 **build.bat** 或 **build.sh**

## 命令行参数

- `-p`：监听端口号。默认值：**9300**
- `-log`：日志目录，存放近 30 天的日志。默认值：**./logs**
- `-data`：数据目录，存放数据库及文件资源。默认值：**./data**
- `-pass`：超级管理员密码。默认值：**Admin123**

## docker 部署

> [https://hub.docker.com/r/streamerzero/gofile](https://hub.docker.com/r/streamerzero/gofile)

- 启动命令：`docker run -d --name gofile --restart always -p 9300:9300 -v /etc/localtime:/etc/localtime:ro -v /home/docker/gofile:/gofile/data -e pass="Admin123" streamerzero/gofile`

### 环境变量

- `-pass`：超级管理员密码。默认值：**Admin123**

### 目录

- 数据文件：**/gofile/data**
- 日志文件：**/gofile/logs**
