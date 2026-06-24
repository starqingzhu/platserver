#!/bin/bash
# filepath: platserver/scripts/build.sh
# 根据当前系统自动编译对应的可执行文件

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}平台服务器编译工具${NC}"
echo -e "${GREEN}========================================${NC}"

echo -e "${BLUE}请选择部署环境（20秒内选择，超时默认本地环境）:${NC}"
echo "1) 本地开发环境 (.devops.yaml)"
echo "2) 测试环境 (.devops_test.yaml)"
echo "3) 生产环境 (.devops_production.yaml)"
echo "4) 提审环境 (.devops_inter.yaml)"

if read -t 20 -p "请输入选项 [1-4, 默认1]: " env_choice; then
    echo ""
else
    echo ""
    echo -e "${YELLOW}超时，使用默认: 本地开发环境${NC}"
    env_choice=1
fi

case ${env_choice:-1} in
    1) CONFIG_FILE=".devops.yaml";             ROBOT_NOTICE_FILE="RobotNotice.yaml";            ENV_NAME="本地开发" ;;
    2) CONFIG_FILE=".devops_test.yaml";        ROBOT_NOTICE_FILE="RobotNotice_test.yaml";        ENV_NAME="测试" ;;
    3) CONFIG_FILE=".devops_production.yaml";  ROBOT_NOTICE_FILE="RobotNotice_production.yaml";  ENV_NAME="生产" ;;
    4) CONFIG_FILE=".devops_inter.yaml";       ROBOT_NOTICE_FILE="RobotNotice_inter.yaml";       ENV_NAME="提审" ;;
    *) CONFIG_FILE=".devops.yaml";             ROBOT_NOTICE_FILE="RobotNotice.yaml";            ENV_NAME="本地开发" ;;
esac

echo -e "${GREEN}已选择: ${ENV_NAME}环境 (${CONFIG_FILE})${NC}"
echo ""

mkdir -p ../bin ../config

OS=$(uname -s)
ARCH=$(uname -m)

case $ARCH in
    x86_64)   ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l)   ARCH="arm" ;;
    i386|i686) ARCH="386" ;;
esac

case $OS in
    Linux*)
        GOOS="linux"
        OUTPUT="../bin/platserver"
        echo "检测到 Linux 系统，编译 Linux 可执行文件..."
        ;;
    Darwin*)
        GOOS="darwin"
        OUTPUT="../bin/platserver"
        echo "检测到 macOS 系统，编译 macOS 可执行文件..."
        ;;
    CYGWIN*|MINGW*|MSYS*)
        GOOS="windows"
        OUTPUT="../bin/platserver.exe"
        echo "检测到 Windows 系统，编译 Windows 可执行文件..."
        ;;
    *)
        echo "未知操作系统: $OS，默认编译 Linux 可执行文件"
        GOOS="linux"
        OUTPUT="../bin/platserver"
        ;;
esac

GIT_COMMIT_SHA1=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date '+%Y-%m-%d %H:%M:%S')

echo "编译参数: GOOS=$GOOS GOARCH=$ARCH"
echo "版本信息: commit=$GIT_COMMIT_SHA1, date=$BUILD_DATE"

GOOS=$GOOS GOARCH=$ARCH go build \
    -ldflags "-X 'main.gitCommitSha1=$GIT_COMMIT_SHA1' -X 'main.date=$BUILD_DATE'" \
    -o $OUTPUT \
    -v ../cmd/main.go

echo ""
echo -e "${YELLOW}复制配置文件...${NC}"

if [ -f "../conf/${CONFIG_FILE}" ]; then
    cp -f "../conf/${CONFIG_FILE}" "../bin/.devops.yaml"
    echo "已复制: ${CONFIG_FILE} -> bin/.devops.yaml"
else
    echo -e "${RED}错误: 配置文件 ../conf/${CONFIG_FILE} 不存在！${NC}"
    exit 1
fi

if [ -d "../config" ]; then
    cp -rf ../config/. "../bin/config/"
    echo "已复制: config/ -> bin/config/"
else
    echo -e "${RED}错误: config 目录不存在！${NC}"
    exit 1
fi

if [ -f "../conf/${ROBOT_NOTICE_FILE}" ]; then
    cp -f "../conf/${ROBOT_NOTICE_FILE}" "../config/RobotNotice.yaml"
    echo "已复制: conf/${ROBOT_NOTICE_FILE} -> /config/RobotNotice.yaml"
else
    echo -e "${RED}错误: 配置文件 conf/${ROBOT_NOTICE_FILE} 不存在！${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo "编译完成！"
echo "操作系统: $OS ($GOOS)"
echo "架构: $ARCH"
echo "Git版本: $GIT_COMMIT_SHA1"
echo "编译时间: $BUILD_DATE"
echo "可执行文件: $OUTPUT"
echo "输出目录: $(cd ../bin && pwd)"
echo "部署环境: ${ENV_NAME}"
echo -e "${GREEN}========================================${NC}"
