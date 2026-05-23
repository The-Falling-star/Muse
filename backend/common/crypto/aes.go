package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// AESGCM AES-GCM加密器
type AESGCM struct {
	key []byte
}

// NewAESGCM 创建AES-GCM加密器
// key: 加密密钥，长度必须是16、24或32字节（对应AES-128、AES-192、AES-256）
func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("密钥长度必须是16、24或32字节")
	}

	return &AESGCM{
		key: key,
	}, nil
}

// NewAESGCMFromKey 从配置的密钥字符串创建AES-GCM加密器
// 将密钥字符串转换为固定长度的字节数组
func NewAESGCMFromKey(key string) (*AESGCM, error) {
	if key == "" {
		return nil, errors.New("密钥不能为空")
	}

	// 使用SHA256哈希将任意长度的密钥转换为32字节的固定长度
	hash := sha256.Sum256([]byte(key))
	return NewAESGCM(hash[:])
}

// Encrypt 加密明文，返回Base64编码的密文
func (a *AESGCM) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// 创建AES cipher
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据（nonce会附加到密文前面）
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密Base64编码的密文
func (a *AESGCM) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建AES cipher
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文长度不足")
	}

	// 分离nonce和实际密文
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
