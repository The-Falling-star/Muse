package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/file"
)

// FileServer 文件服务
type FileServer struct {
	file file.File
}

// NewFileServer 创建一个新的FileServer实例
func NewFileServer() *FileServer {
	return &FileServer{
		file: file.NewFile(),
	}
}

// UploadFile 上传文件
func (f *FileServer) UploadFile(ctx context.Context, req *connect.Request[pb.UploadFileRequest]) (
	*connect.Response[pb.UploadFileResponse], error) {
	resp, err := f.file.UploadFile(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UploadFile", &req.Msg.FileName, resp, err)
	}
	return doResponse(ctx, "UploadFile", req.Msg, resp)
}

// DownloadFile 下载文件
func (f *FileServer) DownloadFile(ctx context.Context, req *connect.Request[pb.DownloadFileRequest]) (
	*connect.Response[pb.DownloadFileResponse], error) {
	resp, err := f.file.DownloadFile(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DownloadFile", req.Msg, resp, err)
	}
	_, _ = doResponse(ctx, "DownloadFile", req.Msg, &resp.FileName)
	return connect.NewResponse(resp), nil
}
