// Package files 提供文件服务功能
package file

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/common/jwt"
	"github.com/ling/muse/config"
	pb "github.com/ling/muse/gen/muse"
	log "github.com/sirupsen/logrus"
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

	relativePath := filepath.Join(fileType.String(), fileName)
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
	userDir := getUserDir(userId, pb.FileType_UploadFileTypeUnspecified)

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
func getUserDir(userId int, fileType pb.FileType) string {
	cfg := config.Get()
	basePath := cfg.File.UploadPath
	if fileType != pb.FileType_UploadFileTypeUnspecified {
		return filepath.Join(basePath, fmt.Sprintf("%d", userId), fileType.String())
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
