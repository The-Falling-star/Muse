#!/usr/bin/env bash
# 生成 JWT Secret 和 API 加密密钥
# 使用 openssl 生成 256-bit 随机密钥，输出为 URL-safe Base64

set -euo pipefail

gen_key() {
  openssl rand -base64 32 | tr '+/' '-_' | tr -d '='
}

JWT_KEY=$(gen_key)
API_KEY=$(gen_key)

echo "请将以下密钥替换到 config.yaml 中："
echo ""
echo "auth:"
echo "  jwt_secret: $JWT_KEY"
echo ""
echo "api_encrypt:"
echo "  encryption_key: $API_KEY"
echo ""
echo "对应 config.yaml 中的字段："
echo "  auth.jwt_secret"
echo "  api_encrypt.encryption_key"
