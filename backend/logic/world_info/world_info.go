package world_info

import (
	"context"
	"encoding/json"
	"strings"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/constant"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/common/jwt"
	"github.com/ling/muse/entity"
	"github.com/ling/muse/entity/sillytavern"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/cache"
	"github.com/ling/muse/repo/database"
)

type worldInfoImpl struct {
	worldInfoRepo database.WorldInfoRepository
}

func newWorldInfo() *worldInfoImpl {
	return &worldInfoImpl{
		worldInfoRepo: database.NewWorldInfoRepo(),
	}
}

func (w *worldInfoImpl) ListWorldInfos(ctx context.Context, req *pb.ListWorldInfosRequest) (*pb.ListWorldInfosResponse, error) {
	// 获取分页参数
	userId := jwt.GetUserId(ctx)
	page, pageSize := constant.NormalizePagination(int(req.GetPage()), int(req.GetPageSize()))

	// 从数据库获取世界书列表（支持分页）
	worldInfos, total, err := w.worldInfoRepo.List(ctx, userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbWorldInfos := make([]*pb.WorldInfo, 0, len(worldInfos))
	for _, worldInfo := range worldInfos {
		pbWorldInfos = append(pbWorldInfos, convert.WorldInfoEntityToPb(worldInfo))
	}

	return &pb.ListWorldInfosResponse{
		WorldInfos: pbWorldInfos,
		Total:      total,
		Page:       int32(page),
		PageSize:   int32(pageSize),
	}, nil
}

func (w *worldInfoImpl) GetWorldInfo(ctx context.Context, req *pb.GetWorldInfoRequest) (*pb.GetWorldInfoResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoID)
	}

	// 从数据库获取世界书
	userId := jwt.GetUserId(ctx)
	worldInfo, err := w.worldInfoRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	if worldInfo == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.WorldInfoNotFound)
	}

	return &pb.GetWorldInfoResponse{
		WorldInfo: convert.WorldInfoEntityToPb(worldInfo),
	}, nil
}

func (w *worldInfoImpl) CreateWorldInfo(ctx context.Context, req *pb.CreateWorldInfoRequest) (*pb.CreateWorldInfoResponse, error) {
	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyWorldInfoName)
	}

	// 构建世界书实体
	userId := jwt.GetUserId(ctx)
	worldInfo := &entity.WorldInfo{
		UserID:      userId,
		Name:        name,
		Description: req.GetDescription(),
		IsGlobal:    req.GetIsGlobal(),
	}

	// 保存到数据库
	if err := w.worldInfoRepo.Create(ctx, worldInfo); err != nil {
		return nil, err
	}

	// 重新获取完整数据（包含关联）
	fullWorldInfo, err := w.worldInfoRepo.GetByID(ctx, worldInfo.ID, userId)
	if err != nil {
		return nil, err
	}

	return &pb.CreateWorldInfoResponse{
		WorldInfo: convert.WorldInfoEntityToPb(fullWorldInfo),
	}, nil
}

func (w *worldInfoImpl) UpdateWorldInfo(ctx context.Context, req *pb.UpdateWorldInfoRequest) (*pb.UpdateWorldInfoResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoID)
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyWorldInfoName)
	}

	// 获取当前世界书
	userId := jwt.GetUserId(ctx)
	worldInfo, err := w.worldInfoRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	if worldInfo == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.WorldInfoNotFound)
	}

	// 更新世界书字段
	worldInfo.Name = name
	if req.Description != nil {
		worldInfo.Description = *req.Description
	}
	worldInfo.IsGlobal = req.IsGlobal

	// 更新数据库
	if err := w.worldInfoRepo.Update(ctx, worldInfo); err != nil {
		return nil, err
	}

	// 使相关缓存失效
	cache.InvalidateCacheByWorldInfo(int64(id))

	// 重新获取更新后的世界书（包含关联数据）
	updatedWorldInfo, err := w.worldInfoRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateWorldInfoResponse{
		WorldInfo: convert.WorldInfoEntityToPb(updatedWorldInfo),
	}, nil
}

func (w *worldInfoImpl) DeleteWorldInfo(ctx context.Context, req *pb.DeleteWorldInfoRequest) (*pb.DeleteWorldInfoResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoID)
	}

	// 使相关缓存失效（在删除之前）
	cache.InvalidateCacheByWorldInfo(int64(id))

	// 删除世界书（会级联删除关联的条目）
	userId := jwt.GetUserId(ctx)
	if err := w.worldInfoRepo.Delete(ctx, id, userId); err != nil {
		return nil, err
	}

	return &pb.DeleteWorldInfoResponse{}, nil
}

func (w *worldInfoImpl) ListWorldInfoEntries(ctx context.Context, req *pb.ListWorldInfoEntriesRequest) (*pb.ListWorldInfoEntriesResponse, error) {
	worldInfoID := int(req.GetWorldInfoId())
	if worldInfoID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoID)
	}

	// 从数据库获取条目列表
	entries, err := w.worldInfoRepo.ListEntries(ctx, worldInfoID)
	if err != nil {
		return nil, err
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
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoID)
	}

	// 构建条目实体
	entry := &entity.WorldInfoEntry{
		WorldInfoID:    worldInfoID,
		UID:            strings.TrimSpace(req.GetUid()),
		Keys:           req.KeysList,
		SecondaryKeys:  req.GetSecondaryKeys(),
		Content:        req.Content,
		Comment:        req.GetComment(),
		IsEnabled:      req.GetIsEnabled(),
		Constant:       req.GetConstant(),
		Selective:      req.GetSelective(),
		InsertionOrder: int(req.GetInsertionOrder()),
		Position:       req.GetPosition(),
		Depth:          int(req.GetDepth()),
		Role:           req.GetRole(),
		SortOrder:      int(req.GetSortOrder()),
	}

	// 保存到数据库
	if err := w.worldInfoRepo.CreateEntry(ctx, entry); err != nil {
		return nil, err
	}

	// 重新获取完整数据
	fullEntry, err := w.worldInfoRepo.GetEntryByID(ctx, entry.ID)
	if err != nil {
		return nil, err
	}

	return &pb.AddWorldInfoEntryResponse{
		Entry: convert.WorldInfoEntryEntityToPb(fullEntry),
	}, nil
}

func (w *worldInfoImpl) UpdateWorldInfoEntry(ctx context.Context, req *pb.UpdateWorldInfoEntryRequest) (*pb.UpdateWorldInfoEntryResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoEntryID)
	}

	// 获取当前条目
	entry, err := w.worldInfoRepo.GetEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.WorldInfoEntryNotFound)
	}

	// 更新条目字段
	entry.UID = strings.TrimSpace(req.GetUid())
	entry.Keys = req.KeysList
	entry.SecondaryKeys = req.GetSecondaryKeys()
	entry.Content = req.Content
	entry.Comment = req.GetComment()
	entry.IsEnabled = req.GetIsEnabled()
	entry.Constant = req.GetConstant()
	entry.Selective = req.GetSelective()
	entry.InsertionOrder = int(req.GetInsertionOrder())
	entry.Position = req.GetPosition()
	entry.Depth = int(req.GetDepth())
	entry.Role = req.GetRole()
	entry.SortOrder = int(req.GetSortOrder())

	// 更新数据库
	if err := w.worldInfoRepo.UpdateEntry(ctx, entry); err != nil {
		return nil, err
	}

	// 使相关缓存失效（条目修改会影响使用该世界书的缓存）
	cache.InvalidateCacheByWorldInfo(int64(entry.WorldInfoID))

	// 重新获取更新后的条目
	updatedEntry, err := w.worldInfoRepo.GetEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateWorldInfoEntryResponse{
		Entry: convert.WorldInfoEntryEntityToPb(updatedEntry),
	}, nil
}

func (w *worldInfoImpl) DeleteWorldInfoEntry(ctx context.Context, req *pb.DeleteWorldInfoEntryRequest) (*pb.DeleteWorldInfoEntryResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoEntryID)
	}

	// 删除条目
	if err := w.worldInfoRepo.DeleteEntry(ctx, id); err != nil {
		return nil, err
	}

	return &pb.DeleteWorldInfoEntryResponse{}, nil
}

func (w *worldInfoImpl) UpdateWorldInfoEntriesOrder(ctx context.Context, req *pb.UpdateWorldInfoEntriesOrderRequest) (*pb.UpdateWorldInfoEntriesOrderResponse, error) {
	worldInfoID := int(req.GetWorldInfoId())
	if worldInfoID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoID)
	}

	if len(req.EntryIds) == 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptySortData)
	}

	// 构建排序映射
	entryOrders := make(map[int]int)
	for i, entryID := range req.EntryIds {
		entryOrders[int(entryID)] = i
	}

	// 更新排序
	if err := w.worldInfoRepo.UpdateEntriesOrder(ctx, worldInfoID, entryOrders); err != nil {
		return nil, err
	}

	return &pb.UpdateWorldInfoEntriesOrderResponse{}, nil
}

func (w *worldInfoImpl) ImportWorldInfo(ctx context.Context, req *pb.ImportWorldInfoRequest) (*pb.ImportWorldInfoResponse, error) {
	// 参数校验
	if len(req.GetFileContent()) == 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyFileContent)
	}

	// 解析 JSON 文件内容
	var stWorldBook sillytavern.WorldBook
	if err := json.Unmarshal(req.GetFileContent(), &stWorldBook); err != nil {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoFile)
	}

	// 校验文件格式（必须包含 entries）
	if stWorldBook.Entries == nil {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoFile)
	}

	// 转换为 Muse 实体
	worldInfo := convert.STWorldInfoToEntity(&stWorldBook)
	userId := jwt.GetUserId(ctx)
	worldInfo.UserID = userId
	// 从文件名获取世界书名称（去掉 .json 后缀）
	if worldInfo.Name == "" {
		worldInfo.Name = strings.TrimSuffix(req.GetFileName(), ".json")
	}
	if worldInfo.Name == "" {
		worldInfo.Name = "Imported World Info"
	}

	// 保存世界书到数据库
	if err := w.worldInfoRepo.Create(ctx, worldInfo); err != nil {
		return nil, err
	}

	// 导入条目
	if len(stWorldBook.Entries) > 0 {
		entries := make([]*entity.WorldInfoEntry, 0, len(stWorldBook.Entries))
		sortOrder := 0
		for _, stEntry := range stWorldBook.Entries {
			entry := convert.STBookEntryToEntity(&stEntry)
			entry.WorldInfoID = worldInfo.ID
			entry.SortOrder = sortOrder
			entries = append(entries, entry)
			sortOrder++
		}

		// 批量创建条目
		if err := w.worldInfoRepo.BatchCreateEntries(ctx, entries); err != nil {
			// 条目创建失败，但世界书已创建，记录日志但不返回错误
			// 可以考虑在未来添加事务支持
		}
	}

	// 重新获取完整的世界书（包含条目）
	fullWorldInfo, err := w.worldInfoRepo.GetByID(ctx, worldInfo.ID, userId)
	if err != nil {
		return nil, err
	}

	return &pb.ImportWorldInfoResponse{
		WorldInfo: convert.WorldInfoEntityToPb(fullWorldInfo),
	}, nil
}

func (w *worldInfoImpl) ExportWorldInfo(ctx context.Context, req *pb.ExportWorldInfoRequest) (*pb.ExportWorldInfoResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidWorldInfoID)
	}

	// 从数据库获取世界书（包含所有条目）
	userId := jwt.GetUserId(ctx)
	worldInfo, err := w.worldInfoRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	if worldInfo == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.WorldInfoNotFound)
	}

	// 转换为 SillyTavern 格式
	stWorldBook := convert.EntityToSTWorldInfo(worldInfo)

	// 序列化为 JSON
	fileContent, jsonErr := json.MarshalIndent(stWorldBook, "", "  ")
	if jsonErr != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "序列化世界书失败: %v", jsonErr)
	}

	// 生成文件名
	fileName := worldInfo.Name + ".json"

	return &pb.ExportWorldInfoResponse{
		FileContent: fileContent,
		FileName:    fileName,
	}, nil
}
