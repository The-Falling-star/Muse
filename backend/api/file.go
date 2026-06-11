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
func (f *FileServer) UploadFile(ctx context.Context, req *connect.Request[pb.UploadFileReq]) (
	*connect.Response[pb.UploadFileRsp], error) {
	resp, err := f.file.UploadFile(ctx, req.Msg)
	req.Msg.FileContent = nil
	if err != nil {
		return doResponseExp(ctx, "UploadFile", req, resp, err)
	}
	return doResponse(ctx, "UploadFile", req, resp)
}

// DownloadFile 下载文件
func (f *FileServer) DownloadFile(ctx context.Context, req *connect.Request[pb.DownloadFileReq]) (
	*connect.Response[pb.DownloadFileRsp], error) {
	resp, err := f.file.DownloadFile(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DownloadFile", req.Msg, resp, err)
	}
	_, _ = doResponse(ctx, "DownloadFile", req.Msg, &resp.FileName)
	return connect.NewResponse(resp), nil
}
