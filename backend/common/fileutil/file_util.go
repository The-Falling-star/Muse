package fileutil

import (
	"bytes"
	"context"
	"fmt"
	imgpng "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/ling/muse/config"
	pb "github.com/ling/muse/gen/muse"
	log "github.com/sirupsen/logrus"
)

// 头像压缩配置
const (
	avatarMaxWidth  = 512 // 头像最大宽度
	avatarMaxHeight = 512 // 头像最大高度
)

// SaveAvatarFile 保存头像文件（自动压缩为PNG格式）
func SaveAvatarFile(_ context.Context, userId int, characterId int, characterName string, imageData []byte, fileType pb.FileType) (
	string, error) {
	cfg := config.Get()
	avatarDir := filepath.Join(cfg.File.UploadPath, fmt.Sprintf("%d", userId), fileType.String())

	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return "", fmt.Errorf("创建头像目录失败: %w", err)
	}

	// 压缩图片
	compressedData, err := compressImage(imageData, avatarMaxWidth, avatarMaxHeight)
	if err != nil {
		return "", fmt.Errorf("压缩图片失败: %w", err)
	}

	// 文件名格式: {character_id}_{character_name}.png
	fileName := fmt.Sprintf("%d_%s.png", characterId, SanitizeFileName(characterName))
	filePath := filepath.Join(avatarDir, fileName)

	if err = os.WriteFile(filePath, compressedData, 0644); err != nil {
		return "", fmt.Errorf("保存头像文件失败: %w", err)
	}

	// 返回相对路径
	relativePath := filepath.Join(fileType.String(), fileName)
	log.Infof("头像保存成功: userId=%d, characterId=%d, path=%s, originalSize=%d, compressedSize=%d",
		userId, characterId, relativePath, len(imageData), len(compressedData))

	return relativePath, nil
}

// compressImage 压缩图片为PNG格式，限制最大尺寸
func compressImage(data []byte, maxWidth, maxHeight int) ([]byte, error) {
	// 解码图片
	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}

	// 获取原始尺寸
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// 如果图片尺寸超过限制，则缩放
	if origWidth > maxWidth || origHeight > maxHeight {
		img = imaging.Resize(img, maxWidth, maxHeight, imaging.Lanczos)
		log.Debugf("图片缩放: %dx%d -> %dx%d", origWidth, origHeight,
			img.Bounds().Dx(), img.Bounds().Dy())
	}

	// 编码为PNG
	var buf bytes.Buffer
	if err = imgpng.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("编码PNG失败: %w", err)
	}

	return buf.Bytes(), nil
}

// SanitizeFileName 清理文件名，移除不安全字符
func SanitizeFileName(name string) string {
	// 移除路径分隔符和其他不安全字符
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	result := replacer.Replace(name)
	// 限制长度
	if len(result) > 50 {
		result = result[:50]
	}
	return result
}

// DeleteAvatarFile 删除头像文件
func DeleteAvatarFile(avatarPath string) error {
	if avatarPath == "" {
		return nil
	}
	if err := os.Remove(avatarPath); err != nil {
		return fmt.Errorf("删除头像文件失败: %w", err)
	}
	return nil
}

// BuildAvatarURL 构建头像URL
func BuildAvatarURL(baseUrl string, userId int, characterId int, characterName string, fileType pb.FileType) string {
	return fmt.Sprintf("%s/%d/%s/%d_%s.png", baseUrl, userId, fileType.String(), characterId,
		SanitizeFileName(characterName))
}
