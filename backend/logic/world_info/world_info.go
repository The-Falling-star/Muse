package world_info

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/database"
)

// 默认用户ID，待认证功能完成后替换
const defaultUserID = 1

type worldInfoImpl struct {
	worldInfoRepo *database.WorldInfoRepo
}

func newWorldInfo() *worldInfoImpl {
	return &worldInfoImpl{
		worldInfoRepo: &database.WorldInfoRepo{},
	}
}

func (w *worldInfoImpl) ListWorldInfos(ctx context.Context, req *pb.ListWorldInfosRequest) (*pb.ListWorldInfosResponse, error) {
	// 从数据库获取世界书列表（不分页）
	worldInfos, _, err := w.worldInfoRepo.List(defaultUserID, 1, 1000)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 转换为pb格式
	pbWorldInfos := make([]*pb.WorldInfo, 0, len(worldInfos))
	for _, worldInfo := range worldInfos {
		pbWorldInfos = append(pbWorldInfos, convert.WorldInfoEntityToPb(worldInfo))
	}

	return &pb.ListWorldInfosResponse{
		WorldInfos: pbWorldInfos,
	}, nil
}

func (w *worldInfoImpl) GetWorldInfo(ctx context.Context, req *pb.GetWorldInfoRequest) (*pb.GetWorldInfoResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的世界书ID"))
	}

	// 从数据库获取世界书
	worldInfo, err := w.worldInfoRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if worldInfo == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("世界书不存在"))
	}

	return &pb.GetWorldInfoResponse{
		WorldInfo: convert.WorldInfoEntityToPb(worldInfo),
	}, nil
}

func (w *worldInfoImpl) CreateWorldInfo(ctx context.Context, req *pb.CreateWorldInfoRequest) (*pb.CreateWorldInfoResponse, error) {
	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("世界书名称不能为空"))
	}

	// 构建世界书实体
	worldInfo := &entity.WorldInfo{
		UserID:      defaultUserID,
		Name:        name,
		Description: req.GetDescription(),
		IsGlobal:    req.GetIsGlobal(),
	}

	// 保存到数据库
	if err := w.worldInfoRepo.Create(worldInfo); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取完整数据（包含关联）
	fullWorldInfo, err := w.worldInfoRepo.GetByID(worldInfo.ID, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.CreateWorldInfoResponse{
		WorldInfo: convert.WorldInfoEntityToPb(fullWorldInfo),
	}, nil
}

func (w *worldInfoImpl) UpdateWorldInfo(ctx context.Context, req *pb.UpdateWorldInfoRequest) (*pb.UpdateWorldInfoResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的世界书ID"))
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("世界书名称不能为空"))
	}

	// 获取当前世界书
	worldInfo, err := w.worldInfoRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if worldInfo == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("世界书不存在"))
	}

	// 更新世界书字段
	worldInfo.Name = name
	if req.Description != nil {
		worldInfo.Description = *req.Description
	}
	worldInfo.IsGlobal = req.IsGlobal

	// 更新数据库
	if err := w.worldInfoRepo.Update(worldInfo); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取更新后的世界书（包含关联数据）
	updatedWorldInfo, err := w.worldInfoRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdateWorldInfoResponse{
		WorldInfo: convert.WorldInfoEntityToPb(updatedWorldInfo),
	}, nil
}

func (w *worldInfoImpl) DeleteWorldInfo(ctx context.Context, req *pb.DeleteWorldInfoRequest) (*pb.DeleteWorldInfoResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的世界书ID"))
	}

	// 删除世界书（会级联删除关联的条目）
	if err := w.worldInfoRepo.Delete(id, defaultUserID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.DeleteWorldInfoResponse{}, nil
}

func (w *worldInfoImpl) ListWorldInfoEntries(ctx context.Context, req *pb.ListWorldInfoEntriesRequest) (*pb.ListWorldInfoEntriesResponse, error) {
	worldInfoID := int(req.GetWorldInfoId())
	if worldInfoID <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的世界书ID"))
	}

	// 从数据库获取条目列表
	entries, err := w.worldInfoRepo.ListEntries(worldInfoID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 转换为pb格式
	pbEntries := make([]*pb.WorldInfoEntry, 0, len(entries))
	for _, entry := range entries {
		pbEntries = append(pbEntries, convert.WorldInfoEntryEntityToPb(entry))
	}

	return &pb.ListWorldInfoEntriesResponse{
		Entries: pbEntries,
	}, nil
}

func (w *worldInfoImpl) AddWorldInfoEntry(ctx context.Context, req *pb.AddWorldInfoEntryRequest) (*pb.AddWorldInfoEntryResponse, error) {
	// 参数校验
	worldInfoID := int(req.GetWorldInfoId())
	if worldInfoID <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的世界书ID"))
	}

	keysList := strings.TrimSpace(req.GetKeysList())
	if keysList == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("关键词列表不能为空"))
	}

	content := strings.TrimSpace(req.GetContent())
	if content == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("内容不能为空"))
	}

	// 构建条目实体
	entry := &entity.WorldInfoEntry{
		WorldInfoID:    worldInfoID,
		UID:            strings.TrimSpace(req.GetUid()),
		KeysList:       keysList,
		SecondaryKeys:  req.GetSecondaryKeys(),
		Content:        content,
		Comment:        req.GetComment(),
		IsEnabled:      req.GetIsEnabled(),
		Constant:       req.GetConstant(),
		Selective:      req.GetSelective(),
		InsertionOrder: int(req.GetInsertionOrder()),
		Position:       req.GetPosition(),
		Depth:          int(req.GetDepth()),
		SortOrder:      int(req.GetSortOrder()),
	}

	// 保存到数据库
	if err := w.worldInfoRepo.CreateEntry(entry); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取完整数据
	fullEntry, err := w.worldInfoRepo.GetEntryByID(entry.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.AddWorldInfoEntryResponse{
		Entry: convert.WorldInfoEntryEntityToPb(fullEntry),
	}, nil
}

func (w *worldInfoImpl) UpdateWorldInfoEntry(ctx context.Context, req *pb.UpdateWorldInfoEntryRequest) (*pb.UpdateWorldInfoEntryResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的条目ID"))
	}

	// 参数校验
	keysList := strings.TrimSpace(req.GetKeysList())
	if keysList == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("关键词列表不能为空"))
	}

	content := strings.TrimSpace(req.GetContent())
	if content == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("内容不能为空"))
	}

	// 获取当前条目
	entry, err := w.worldInfoRepo.GetEntryByID(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if entry == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("条目不存在"))
	}

	// 更新条目字段
	entry.UID = strings.TrimSpace(req.GetUid())
	entry.KeysList = keysList
	entry.SecondaryKeys = req.GetSecondaryKeys()
	entry.Content = content
	entry.Comment = req.GetComment()
	entry.IsEnabled = req.GetIsEnabled()
	entry.Constant = req.GetConstant()
	entry.Selective = req.GetSelective()
	entry.InsertionOrder = int(req.GetInsertionOrder())
	entry.Position = req.GetPosition()
	entry.Depth = int(req.GetDepth())
	entry.SortOrder = int(req.GetSortOrder())

	// 更新数据库
	if err := w.worldInfoRepo.UpdateEntry(entry); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取更新后的条目
	updatedEntry, err := w.worldInfoRepo.GetEntryByID(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdateWorldInfoEntryResponse{
		Entry: convert.WorldInfoEntryEntityToPb(updatedEntry),
	}, nil
}

func (w *worldInfoImpl) DeleteWorldInfoEntry(ctx context.Context, req *pb.DeleteWorldInfoEntryRequest) (*pb.DeleteWorldInfoEntryResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的条目ID"))
	}

	// 删除条目
	if err := w.worldInfoRepo.DeleteEntry(id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.DeleteWorldInfoEntryResponse{}, nil
}

func (w *worldInfoImpl) UpdateWorldInfoEntriesOrder(ctx context.Context, req *pb.UpdateWorldInfoEntriesOrderRequest) (*pb.UpdateWorldInfoEntriesOrderResponse, error) {
	worldInfoID := int(req.GetWorldInfoId())
	if worldInfoID <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的世界书ID"))
	}

	if len(req.EntryIds) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("排序数据不能为空"))
	}

	// 构建排序映射
	entryOrders := make(map[int]int)
	for i, entryID := range req.EntryIds {
		entryOrders[int(entryID)] = i
	}

	// 更新排序
	if err := w.worldInfoRepo.UpdateEntriesOrder(worldInfoID, entryOrders); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdateWorldInfoEntriesOrderResponse{}, nil
}

func (w *worldInfoImpl) ImportWorldInfo(ctx context.Context, req *pb.ImportWorldInfoRequest) (*pb.ImportWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoImpl) ExportWorldInfo(ctx context.Context, req *pb.ExportWorldInfoRequest) (*pb.ExportWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}
