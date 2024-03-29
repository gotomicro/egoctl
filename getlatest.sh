#!/bin/bash
RED='\033[1;31m'
GREEN='\033[1;32m'
CYAN='\033[0;36m'
NC='\033[0m'

# 使用方法：./getlatest.sh
#         os=linux ./getlatest.sh
#         os=osx ./getlatest.sh
echo -e "${CYAN}[egoctl-tools] 下载 protoc、google/protobuf、protoc-gen-go、egoctl，根据网络情况，整个过程可能耗时 1-2 分钟左右${NC}"
echo -e "${CYAN}[egoctl-tools] 安装过程中可能需要输入 root 密码\n${NC}"

unameOut=$(uname -s)
arch=$(uname -m)
# 如果未设置os，则自动查询os，可选值linux\osx
if [[ -z ${os+x} ]]; then
  case "${unameOut}" in
      Linux*)     os=linux;;
      Darwin*)    os=osx;;
      *)          exit 1;;
  esac
fi
# 根据os可选值linux\osx，设置unameOut、GOOS等变量
case "${os}" in
  linux*)     os=linux && unameOut="Linux" && export GOOS=linux;;
  osx*)    os=osx && unameOut="Darwin" && export GOOS=darwin;;
  *)          exit 1;;
esac

echo unameOut:${unameOut}
echo os:${os}
echo arch:${arch}
echo GOOS:${GOOS}

protocVersion=3.17.3
protocGenGoVersion=1.27.1
protocGenGoGRPCVersion=1.39.1
protocGenOpenapiv2Version=2.6.0
protocGenGoErrorsVersion=1.1.1
protocGenGoTestVersion=1.1.1
protocGenGoHttpVersion=1.1.1
egoctlVersion=1.0.6
githubUrl=https://github.com

# 初始化目标目录
targetDir=""
if [[ ${unameOut} == "Linux" ]]; then
  targetDir=$HOME/.cache/prototool/${unameOut}/x86_64/protobuf/${protocVersion}
else
  targetDir=$HOME/Library/Caches/prototool/${unameOut}/x86_64/protobuf/${protocVersion}
fi
mkdir -p ${targetDir}/bin
mkdir -p ${targetDir}/include


# 初始化download目录
tmpDir=/tmp/.egoctl-scripts
rm -rf ${tmpDir}/*
mkdir -p ${tmpDir}
cd ${tmpDir}

function down_protoc() {
    echo -e "${CYAN}[egoctl-tools] 需要下载 protoc-${protocVersion} 并移动至 ${targetDir} 下么？${NC} (y/n)"
    read download
    if [ "$download" != "${download#[Yy]}"  ] ;then
        echo -e "${CYAN}[egoctl-tools] 下载并配置 protoc、google/protobuf 中...${NC}"
        protocTmp=${tmpDir}/protoc-${protocVersion}
        protoZip=protoc-${protocVersion}-${os}-x86_64.zip
        wget ${githubUrl}/protocolbuffers/protobuf/releases/download/v${protocVersion}/${protoZip} --show-progress
        unzip -q -o $protoZip -d ${protocTmp}

        # 复制protoc到目标目录，并软链到/usr/local/bin下
        cp -r ${protocTmp}/bin/* ${targetDir}/bin/
        sudo ln -sf ${targetDir}/bin/protoc /usr/local/bin/protoc

        # 复制google/protobuf到目标目录，并软链到/usr/local/include下
        cp -r ${protocTmp}/include/* ${targetDir}/include/
        sudo ln -sf ${targetDir}/include/google /usr/local/include/google
        echo -e "${GREEN}[egoctl-tools] 下载并配置 protoc、google/protobuf 成功！${NC}"
    fi
    echo -e "\n"
}

function down_egoctl() {
    echo -e "${CYAN}[egoctl-tools] 下载并配置 $1 中...${NC}"
    wget -O ${tmpDir}/$1-${unameOut}.tar.gz $2 -q --show-progress
    tar -C ${tmpDir} -xvf ${tmpDir}/$1-${unameOut}.tar.gz
    chmod +x ${tmpDir}/$1
    sudo cp -f ${tmpDir}/$1 /usr/local/bin/$1
    echo -e "${GREEN}[egoctl-tools] 下载并配置 $1 成功！${NC}"
    echo -e "\n"
}

# https://github.com/protocolbuffers/protobuf-go
function down_protoc_gen_go() {
  echo -e "${CYAN}[egoctl-tools] 下载并配置 protoc-gen-go 中...${NC}"
  git clone --quiet --depth 1 --branch v${protocGenGoVersion} ${githubUrl}/protocolbuffers/protobuf-go ${tmpDir}/protobuf > /dev/null 2>&1
  cd ${tmpDir}/protobuf/cmd/protoc-gen-go && go build
  sudo cp -f ${tmpDir}/protobuf/cmd/protoc-gen-go/protoc-gen-go /usr/local/bin/protoc-gen-go
  echo -e "${GREEN}[egoctl-tools] 下载并配置 protoc-gen-go 成功！${NC}"
  echo -e "\n"
}

function down_protoc_gen_go_grpc() {
  echo -e "${CYAN}[egoctl-tools] 下载并配置 protoc-gen-go-grpc 中...${NC}"
  git clone --quiet --depth 1 --branch v${protocGenGoGRPCVersion} ${githubUrl}/grpc/grpc-go ${tmpDir}/grpc-go > /dev/null 2>&1
  cd ${tmpDir}/grpc-go/cmd/protoc-gen-go-grpc && go build
  sudo cp -f ${tmpDir}/grpc-go/cmd/protoc-gen-go-grpc/protoc-gen-go-grpc /usr/local/bin/protoc-gen-go-grpc
  echo -e "${GREEN}[egoctl-tools] 下载并配置 protoc-gen-go-grpc 成功！${NC}"
  echo -e "\n"
}

function down_protoc_gen_openapiv2() {
  echo -e "${CYAN}[egoctl-tools] 下载并配置 protoc-gen-openapiv2 中...${NC}"
  git clone --quiet --depth 1 --branch v${protocGenOpenapiv2Version} ${githubUrl}/grpc-ecosystem/grpc-gateway ${tmpDir}/grpc-gateway > /dev/null 2>&1
  cd ${tmpDir}/grpc-gateway/protoc-gen-openapiv2 && go build
  sudo cp -f ${tmpDir}/grpc-gateway/protoc-gen-openapiv2/protoc-gen-openapiv2 /usr/local/bin/protoc-gen-openapiv2
  echo -e "${GREEN}[egoctl-tools] 下载并配置 protoc-gen-openapiv2 成功！${NC}"
  echo -e "\n"
}

# 下载protoc
down_protoc

# 下载protoc-gen-go
down_protoc_gen_go

# 下载protoc-gen-go-grpc
down_protoc_gen_go_grpc

# 下载protoc-gen-openapiv2
down_protoc_gen_openapiv2

# 下载protoc-gen-go-errors
down_egoctl protoc-gen-go-errors ${githubUrl}/gotomicro/ego/releases/download/v${protocGenGoErrorsVersion}/protoc-gen-go-errors-${protocGenGoErrorsVersion}-${unameOut}-${arch}.tar.gz

# 下载protoc-gen-go-test
down_egoctl protoc-gen-go-test ${githubUrl}/gotomicro/ego/releases/download/v${protocGenGoTestVersion}/protoc-gen-go-test-${protocGenGoTestVersion}-${unameOut}-${arch}.tar.gz

# 下载protoc-gen-go-http
down_egoctl protoc-gen-go-http ${githubUrl}/gotomicro/ego/releases/download/v${protocGenGoHttpVersion}/protoc-gen-go-http-${protocGenGoHttpVersion}-${unameOut}-${arch}.tar.gz

# 下载egoctl
down_egoctl egoctl ${githubUrl}/gotomicro/egoctl/releases/download/v${egoctlVersion}/egoctl-${egoctlVersion}-${unameOut}-${arch}.tar.gz

echo -e "${GREEN}[egoctl-tools] 配置 protoc protoc-gen-go、protoc-gen-go-grpc、protoc-gen-openapiv2、protoc-gen-go-errors、protoc-gen-go-http、egoctl 成功!${NC}"
which protoc protoc-gen-go protoc-gen-go-grpc protoc-gen-openapiv2 protoc-gen-go-errors protoc-gen-go-http egoctl

echo "protoc version:" $(/usr/local/bin/protoc --version)
echo "protoc-gen-go version:" $(/usr/local/bin/protoc-gen-go -version)
echo "protoc-gen-go-grpc version:" $(/usr/local/bin/protoc-gen-go-grpc -version)
echo "protoc-gen-openapiv2 version:" $(/usr/local/bin/protoc-gen-openapiv2 -version)
echo "protoc-gen-go-errors version:" $(/usr/local/bin/protoc-gen-go-errors -version)
echo "protoc-gen-go-http version:" $(/usr/local/bin/protoc-gen-go-http -version)
echo "egoctl version:" $(/usr/local/bin/egoctl version | grep buildGitVersion)

echo -e "\n"
exit 0
