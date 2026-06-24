#!/bin/bash
# filepath: platserver/scripts/package.sh
# 打包脚本：编译并打包服务器程序及配置文件

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_NAME="platserver"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}平台服务器打包工具${NC}"
echo -e "${GREEN}========================================${NC}"

echo -e "${BLUE}请选择编译平台（60秒内选择，超时默认Linux amd64）:${NC}"
echo "1) Linux (amd64)"
echo "2) Linux (arm64)"
echo "3) Windows (amd64)"
echo "4) macOS (amd64)"
echo "5) macOS (arm64)"

if read -t 60 -p "请输入选项 [1-5, 默认1]: " platform_choice; then echo ""; else echo ""; platform_choice=1; fi

case ${platform_choice:-1} in
    1) TARGET_OS="linux";   TARGET_ARCH="amd64"; BINARY_NAME="platserver" ;;
    2) TARGET_OS="linux";   TARGET_ARCH="arm64"; BINARY_NAME="platserver" ;;
    3) TARGET_OS="windows"; TARGET_ARCH="amd64"; BINARY_NAME="platserver.exe" ;;
    4) TARGET_OS="darwin";  TARGET_ARCH="amd64"; BINARY_NAME="platserver" ;;
    5) TARGET_OS="darwin";  TARGET_ARCH="arm64"; BINARY_NAME="platserver" ;;
    *) TARGET_OS="linux";   TARGET_ARCH="amd64"; BINARY_NAME="platserver" ;;
esac

echo ""
echo -e "${BLUE}请选择部署环境（60秒内选择，超时默认本地环境）:${NC}"
echo "1) 本地开发环境 (.devops.yaml)"
echo "2) 测试环境 (.devops_test.yaml)"
echo "3) 生产环境 (.devops_production.yaml)"
echo "4) 提审环境 (.devops_inter.yaml)"

if read -t 60 -p "请输入选项 [1-4, 默认1]: " env_choice; then echo ""; else echo ""; env_choice=1; fi

case ${env_choice:-1} in
    1) CONFIG_FILE=".devops.yaml";             ROBOT_NOTICE_FILE="RobotNotice.yaml";            ENV_NAME="本地开发" ;;
    2) CONFIG_FILE=".devops_test.yaml";        ROBOT_NOTICE_FILE="RobotNotice_test.yaml";        ENV_NAME="测试" ;;
    3) CONFIG_FILE=".devops_production.yaml";  ROBOT_NOTICE_FILE="RobotNotice_production.yaml";  ENV_NAME="生产" ;;
    4) CONFIG_FILE=".devops_inter.yaml";       ROBOT_NOTICE_FILE="RobotNotice_inter.yaml";       ENV_NAME="提审" ;;
    *) CONFIG_FILE=".devops.yaml";             ROBOT_NOTICE_FILE="RobotNotice.yaml";            ENV_NAME="本地开发" ;;
esac

VERSION=${1:-$(date +"%Y%m%d_%H%M%S")}
PACKAGE_DIR="../release"
TEMP_DIR="${PACKAGE_DIR}/${PROJECT_NAME}_${VERSION}"

echo ""
echo -e "${GREEN}打包配置:${NC}"
echo "  项目: ${PROJECT_NAME}"
echo "  版本: ${VERSION}"
echo "  平台: ${TARGET_OS}/${TARGET_ARCH}"
echo "  环境: ${ENV_NAME} (${CONFIG_FILE})"

# 1. 准备目录
echo ""
echo -e "${YELLOW}[1/6] 准备打包环境...${NC}"
mkdir -p "${PACKAGE_DIR}"
find "${PACKAGE_DIR}" -maxdepth 1 -name "${PROJECT_NAME}_*" -type d -exec rm -rf {} + 2>/dev/null || true
mkdir -p "${TEMP_DIR}/bin" "${TEMP_DIR}/logs" "${TEMP_DIR}/scripts" "${TEMP_DIR}/config"

# 2. 编译
echo -e "${YELLOW}[2/6] 编译程序...${NC}"
cd "$(dirname "$0")"
GOOS=$TARGET_OS GOARCH=$TARGET_ARCH go build \
    -ldflags "-s -w -X 'main.gitCommitSha1=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)' -X 'main.date=$(date +%Y-%m-%d\ %H:%M:%S)'" \
    -o "${TEMP_DIR}/bin/${BINARY_NAME}" ../cmd/main.go
echo -e "${GREEN}编译成功: ${BINARY_NAME}${NC}"

# 3. 复制配置
echo -e "${YELLOW}[3/6] 复制配置文件...${NC}"
if [ -f "../conf/${CONFIG_FILE}" ]; then
    cp -f "../conf/${CONFIG_FILE}" "${TEMP_DIR}/bin/.devops.yaml"
    echo "已复制: ${CONFIG_FILE} -> bin/.devops.yaml"
else
    echo -e "${RED}错误: 配置文件 ${CONFIG_FILE} 不存在！${NC}"
    exit 1
fi

if [ -d "../config" ]; then
    cp -rf ../config/. "${TEMP_DIR}/config/"
    echo "已复制: config/ -> config/"
else
    echo -e "${RED}错误: config 目录不存在！${NC}"
    exit 1
fi

if [ -f "../conf/${ROBOT_NOTICE_FILE}" ]; then
    cp -f "../conf/${ROBOT_NOTICE_FILE}" "${TEMP_DIR}/config/RobotNotice.yaml"
    echo "已复制: conf/${ROBOT_NOTICE_FILE} -> config/RobotNotice.yaml"
else
    echo -e "${RED}错误: 配置文件 conf/${ROBOT_NOTICE_FILE} 不存在！${NC}"
    exit 1
fi

# 4. 复制脚本
echo -e "${YELLOW}[4/6] 复制脚本文件...${NC}"
cp -f start.sh stop.sh restart.sh status.sh "${TEMP_DIR}/scripts/"
chmod +x "${TEMP_DIR}/scripts/"*.sh
echo "已复制: 启动/停止/重启/状态脚本"

# 5. 版本信息
echo -e "${YELLOW}[5/6] 生成版本信息...${NC}"
cat > "${TEMP_DIR}/VERSION" << EOF
Project:    ${PROJECT_NAME}
Version:    ${VERSION}
Build Time: $(date '+%Y-%m-%d %H:%M:%S')
Build OS:   ${TARGET_OS}/${TARGET_ARCH}
Git Branch: $(git branch --show-current 2>/dev/null || echo unknown)
Git Commit: $(git rev-parse --short HEAD 2>/dev/null || echo unknown)
EOF

# 6. 打包压缩
echo -e "${YELLOW}[6/6] 打包压缩...${NC}"
PACKAGE_NAME="${PROJECT_NAME}_${VERSION}_${TARGET_OS}_${TARGET_ARCH}.tar.gz"
cd "${PACKAGE_DIR}"
tar -czf "${PACKAGE_NAME}" "$(basename ${TEMP_DIR})"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}打包完成！${NC}"
echo "包名称: ${PACKAGE_NAME}"
echo "包路径: $(pwd)/${PACKAGE_NAME}"
echo "包大小: $(du -h "${PACKAGE_NAME}" | cut -f1)"
echo ""
echo -e "${YELLOW}部署步骤：${NC}"
echo "1. 上传 ${PACKAGE_NAME} 到服务器"
echo "2. tar -xzf ${PACKAGE_NAME}"
echo "3. cd $(basename ${TEMP_DIR})"
echo "4. vi bin/.devops.yaml  # 确认 wechatWebhookURL"
echo "5. chmod +x bin/${BINARY_NAME} scripts/*.sh"
echo "6. cd scripts && ./start.sh"
echo -e "${GREEN}========================================${NC}"
