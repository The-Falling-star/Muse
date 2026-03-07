// Package file 提供文件服务功能
package file

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image/png"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	"github.com/disintegration/imaging"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/common/jwt"
	"github.com/ling/muse/config"
	pb "github.com/ling/muse/gen/muse"
	log "github.com/sirupsen/logrus"
)

// 文件类型目录常量
const (
	FileTypeAvatar = "avatar" // 头像文件目录名
)

// 头像压缩配置
const (
	avatarMaxWidth  = 512 // 头像最大宽度
	avatarMaxHeight = 512 // 头像最大高度
)

// File 定义了文件服务的接口
type File interface {
	// UploadFile 上传文件
	UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error)
	// DownloadFile 下载文件
	DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error)
}

type fileImpl struct{}

// NewFile 创建一个新的File实例
func NewFile() File {
	return &fileImpl{}
}

// UploadFile 上传文件
func (f *fileImpl) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	fileContent := req.GetFileContent()
	fileName := req.GetFileName()
	fileType := req.GetFileType()

	if len(fileContent) == 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, "文件内容不能为空")
	}

	cfg := config.Get()
	maxSize := int64(cfg.File.MaxUploadSize) * 1024 * 1024
	if int64(len(fileContent)) > maxSize {
		return nil, errs.NewStandardf(connect.CodeInvalidArgument, "文件大小超过限制，最大允许%dMB", cfg.File.MaxUploadSize)
	}

	userId := jwt.GetUserId(ctx)
	userDir := getUserDir(userId, fileType)

	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "创建目录失败: %v", err)
	}

	filePath := filepath.Join(userDir, fileName)
	if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "保存文件失败: %v", err)
	}

	relativePath := filepath.Join(fileType, fileName)
	log.Infof("文件上传成功: userId=%d, path=%s", userId, relativePath)

	return &pb.UploadFileResponse{
		FilePath: relativePath,
	}, nil
}

// DownloadFile 下载文件
func (f *fileImpl) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	filePath := req.GetFilePath()
	if filePath == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, "文件路径不能为空")
	}

	userId := jwt.GetUserId(ctx)
	userDir := getUserDir(userId, "")

	// 安全校验：确保文件路径在用户目录下
	absUserDir, err := filepath.Abs(userDir)
	if err != nil {
		return nil, errs.NewStandard(connect.CodeInternal, "获取用户目录绝对路径失败")
	}

	// 构建完整文件路径
	fullPath := filepath.Join(absUserDir, filePath)
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, "无效的文件路径")
	}

	// 安全检查：确保解析后的路径在用户目录内
	if !strings.HasPrefix(absFilePath, absUserDir+string(filepath.Separator)) {
		log.Warnf("非法文件访问尝试: userId=%d, path=%s", userId, filePath)
		return nil, errs.NewStandard(connect.CodePermissionDenied, "无权访问该文件")
	}

	// 读取文件
	fileContent, err := os.ReadFile(absFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errs.NewStandard(connect.CodeNotFound, "文件不存在")
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "读取文件失败: %v", err)
	}

	// 检测MIME类型
	contentType := detectContentType(fileContent, filepath.Base(filePath))
	fileName := filepath.Base(filePath)

	log.Debugf("文件下载成功: userId=%d, path=%s, size=%d", userId, filePath, len(fileContent))

	return &pb.DownloadFileResponse{
		FileContent: fileContent,
		ContentType: contentType,
		FileName:    fileName,
	}, nil
}

// getUserDir 获取用户的文件目录
func getUserDir(userId int, fileType string) string {
	cfg := config.Get()
	basePath := cfg.File.UploadPath
	if fileType != "" {
		return filepath.Join(basePath, fmt.Sprintf("%d", userId), fileType)
	}
	return filepath.Join(basePath, fmt.Sprintf("%d", userId))
}

// detectContentType 检测文件内容类型
func detectContentType(data []byte, fileName string) string {
	// 首先尝试从文件内容检测
	if len(data) > 512 {
		contentType := http.DetectContentType(data[:512])
		if contentType != "application/octet-stream" {
			return contentType
		}
	} else if len(data) > 0 {
		contentType := http.DetectContentType(data)
		if contentType != "application/octet-stream" {
			return contentType
		}
	}

	// 然后尝试从文件扩展名推断
	ext := filepath.Ext(fileName)
	if ext != "" {
		contentType := mime.TypeByExtension(ext)
		if contentType != "" {
			return contentType
		}
	}

	return "application/octet-stream"
}

// ComputeFileHash 计算文件内容的SHA256哈希值
func ComputeFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// SaveAvatarFile 保存头像文件（自动压缩为PNG格式）
func SaveAvatarFile(ctx context.Context, userId int, characterId int, characterName string, imageData []byte) (string, error) {
	cfg := config.Get()
	avatarDir := filepath.Join(cfg.File.UploadPath, fmt.Sprintf("%d", userId), FileTypeAvatar)

	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return "", fmt.Errorf("创建头像目录失败: %w", err)
	}

	// 压缩图片
	compressedData, err := compressImage(imageData, avatarMaxWidth, avatarMaxHeight)
	if err != nil {
		return "", fmt.Errorf("压缩图片失败: %w", err)
	}

	// 文件名格式: {character_id}_{character_name}.png
	fileName := fmt.Sprintf("%d_%s.png", characterId, sanitizeFileName(characterName))
	filePath := filepath.Join(avatarDir, fileName)

	if err = os.WriteFile(filePath, compressedData, 0644); err != nil {
		return "", fmt.Errorf("保存头像文件失败: %w", err)
	}

	// 返回相对路径
	relativePath := filepath.Join(FileTypeAvatar, fileName)
	log.Infof("头像保存成功: userId=%d, characterId=%d, path=%s, originalSize=%d, compressedSize=%d",
		userId, characterId, relativePath, len(imageData), len(compressedData))

	return relativePath, nil
}

// ReadAvatarFile 读取头像文件
func ReadAvatarFile(ctx context.Context, userId int, avatarPath string) ([]byte, error) {
	if avatarPath == "" {
		return nil, nil
	}

	cfg := config.Get()
	userDir := filepath.Join(cfg.File.UploadPath, fmt.Sprintf("%d", userId))
	fullPath := filepath.Join(userDir, avatarPath)

	// 安全检查：确保路径在用户目录下
	absUserDir, err := filepath.Abs(userDir)
	if err != nil {
		return nil, fmt.Errorf("获取用户目录绝对路径失败: %w", err)
	}

	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无效的文件路径: %w", err)
	}

	if !strings.HasPrefix(absFilePath, absUserDir+string(filepath.Separator)) {
		return nil, fmt.Errorf("非法文件路径访问")
	}

	data, err := os.ReadFile(absFilePath)
	if err != nil {
		return nil, fmt.Errorf("读取头像文件失败: %w", err)
	}

	return data, nil
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

// sanitizeFileName 清理文件名，移除不安全字符
func sanitizeFileName(name string) string {
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
	if err = png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("编码PNG失败: %w", err)
	}

	return buf.Bytes(), nil
}

// BuildAvatarURL 构建头像URL
func BuildAvatarURL(baseUrl string, userId int, characterId int, characterName string) string {
	return fmt.Sprintf("%s/%d/%s/%d_%s.png", baseUrl, userId, FileTypeAvatar, characterId,
		sanitizeFileName(characterName))
}
