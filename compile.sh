#!/bin/bash
set -e

# Muse 编译脚本 - 将前端资源嵌入到 Go 二进制文件中

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
FRONTEND_DIR="$PROJECT_DIR/frontend"
BACKEND_DIR="$PROJECT_DIR/backend"
STATIC_DIR="$BACKEND_DIR/static"

echo "=== Muse 编译脚本 ==="

# 1. 构建前端
echo "[1/3] 构建前端..."
if [ -d "$FRONTEND_DIR" ] && [ -f "$FRONTEND_DIR/package.json" ]; then
  cd "$FRONTEND_DIR"
  npm run build
  echo "  前端构建完成"
else
  echo "  跳过: 未找到前端项目"
  exit 1
fi

# 2. 复制前端产物到 backend/static
echo "[2/3] 复制前端产物..."
rm -rf "$STATIC_DIR"
cp -r "$FRONTEND_DIR/dist" "$STATIC_DIR"
echo "  已复制到 $STATIC_DIR"

# 3. 编译后端
echo "[3/3] 编译后端..."
cd "$BACKEND_DIR"
go build -tags embed -o ../muse .
echo "编译完成: $PROJECT_DIR/muse"
