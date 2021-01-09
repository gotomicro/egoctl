# egoctl
## Requirements

- Go version >= 1.13.

## Installation


## 快速上手

```bash
egoctl -h # 查看使用帮助
```

### 快速生成代码

- 初始化目录和配置文件
```bash
# 创建demo目录
mkdir -p ~/demo

# 下载egoctl.toml样例配置
cd ~/demo
wget https://github.com/gotomicro/egoctl-tmpls/-/raw/master/example/ego.go

go mod init demo
egoctl gen code 
```