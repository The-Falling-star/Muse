package preset

import (
	"context"
	"encoding/json"
	"fmt"
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
	log "github.com/sirupsen/logrus"
)

type presetImpl struct {
	presetRepo    database.PresetRepository
	regexRuleRepo database.RegexRuleRepository
}

func newPreset() *presetImpl {
	return &presetImpl{
		presetRepo:    database.NewPresetRepo(),
		regexRuleRepo: database.NewRegexRuleRepo(),
	}
}

func (p *presetImpl) ListPresets(ctx context.Context, req *pb.ListPresetsRequest) (*pb.ListPresetsResponse, error) {
	// 获取并规范化分页参数
	page, pageSize := constant.NormalizePagination(int(req.GetPage()), int(req.GetPageSize()))

	// 从数据库获取预设列表（不加载关联的PromptItems和RegexRules）
	userId := jwt.GetUserId(ctx)
	presets, presetProLen, total, err := p.presetRepo.List(ctx, userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbPresets := make([]*pb.PresetWithPromptLen, 0, len(presets))
	for i := range presets {
		pbPreset, _, _ := convert.PresetEntityToPb(presets[i])
		pbPresets = append(pbPresets, &pb.PresetWithPromptLen{
			Preset:    pbPreset,
			PromptLen: int32(presetProLen[i]),
		})
	}

	return &pb.ListPresetsResponse{
		Presets:  pbPresets,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

func (p *presetImpl) GetPreset(ctx context.Context, req *pb.GetPresetRequest) (*pb.GetPresetResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	// 从数据库获取预设
	userId := jwt.GetUserId(ctx)
	preset, err := p.presetRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	if preset == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.PresetNotFound)
	}

	pbPreset, pbPromptItem, pbRegex := convert.PresetEntityToPb(preset)
	return &pb.GetPresetResponse{
		Preset: &pb.PresetWithAll{
			Preset:      pbPreset,
			PromptItems: pbPromptItem,
			RegexRules:  pbRegex,
		},
	}, nil
}

func (p *presetImpl) CreatePreset(ctx context.Context, req *pb.CreatePresetRequest) (*pb.CreatePresetResponse, error) {
	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPresetName)
	}

	// 构建预设实体
	userId := jwt.GetUserId(ctx)
	preset := &entity.Preset{
		UserID:           userId,
		Name:             name,
		Temperature:      req.GetTemperature(),
		TopP:             req.GetTopP(),
		TopK:             int(req.GetTopK()),
		MaxTokens:        int(req.GetMaxTokens()),
		FrequencyPenalty: req.GetFrequencyPenalty(),
		PresencePenalty:  req.GetPresencePenalty(),
		Version:          1,
	}

	// 保存到数据库
	if err := p.presetRepo.Create(ctx, preset); err != nil {
		return nil, err
	}

	// 创建关联的提示项
	for _, itemReq := range req.GetPromptItems() {
		item := &entity.PromptItem{
			PresetID:          preset.ID,
			Identifier:        itemReq.GetIdentifier(),
			Name:              strings.TrimSpace(itemReq.GetName()),
			Content:           itemReq.GetContent(),
			Role:              itemReq.GetRole(),
			IsEnabled:         itemReq.GetIsEnabled(),
			InjectionPosition: itemReq.GetInjectionPosition(),
			InjectionDepth:    int(itemReq.GetInjectionDepth()),
			ForbidOverrides:   itemReq.GetForbidOverrides(),
		}
		// 创建提示项，忽略单个失败继续创建其他项
		_ = p.presetRepo.CreatePromptItem(ctx, item)
	}

	return &pb.CreatePresetResponse{}, nil
}

func (p *presetImpl) UpdatePreset(ctx context.Context, req *pb.UpdatePresetRequest) (*pb.UpdatePresetResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPresetName)
	}

	// 获取当前预设
	userId := jwt.GetUserId(ctx)
	preset, err := p.presetRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	if preset == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.PresetNotFound)
	}

	// 检查版本号
	if int64(preset.Version) != req.GetVersion() {
		return nil, errs.NewStandard(connect.CodeAborted, errs.DataConflict)
	}

	// 更新预设字段
	preset.Name = name
	preset.Temperature = req.Temperature
	preset.TopP = req.TopP
	preset.TopK = int(req.TopK)
	preset.MaxTokens = int(req.MaxTokens)
	preset.FrequencyPenalty = req.FrequencyPenalty
	preset.PresencePenalty = req.PresencePenalty

	// 更新数据库
	if err = p.presetRepo.Update(ctx, preset); err != nil {
		return nil, err
	}

	// 使相关缓存失效
	cache.InvalidateCacheByPreset(int64(id))

	return &pb.UpdatePresetResponse{}, nil
}

func (p *presetImpl) DeletePreset(ctx context.Context, req *pb.DeletePresetRequest) (*pb.DeletePresetResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	// 使相关缓存失效（在删除之前）
	cache.InvalidateCacheByPreset(int64(id))

	// 删除预设（会级联删除关联的提示项）
	userId := jwt.GetUserId(ctx)
	if err := p.presetRepo.Delete(ctx, id, userId); err != nil {
		return nil, err
	}

	return &pb.DeletePresetResponse{}, nil
}

func (p *presetImpl) ListPromptItems(ctx context.Context, req *pb.ListPromptItemsRequest) (*pb.ListPromptItemsResponse, error) {
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	// 从数据库获取提示项列表
	items, err := p.presetRepo.ListPromptItems(ctx, presetID)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbItems := make([]*pb.PromptItem, 0, len(items))
	for _, item := range items {
		pbItems = append(pbItems, convert.PromptItemEntityToPb(item))
	}

	return &pb.ListPromptItemsResponse{
		Items: pbItems,
	}, nil
}

func (p *presetImpl) AddPromptItem(ctx context.Context, req *pb.AddPromptItemRequest) (*pb.AddPromptItemResponse, error) {
	// 参数校验
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPromptItemName)
	}

	// 构建提示项实体
	item := &entity.PromptItem{
		PresetID:          presetID,
		Identifier:        req.GetIdentifier(),
		Name:              name,
		Content:           req.GetContent(),
		Role:              req.GetRole(),
		IsEnabled:         req.GetIsEnabled(),
		InjectionPosition: req.GetInjectionPosition(),
		InjectionDepth:    int(req.GetInjectionDepth()),
		ForbidOverrides:   req.GetForbidOverrides(),
	}

	// 保存到数据库
	if err := p.presetRepo.CreatePromptItem(ctx, item); err != nil {
		return nil, err
	}

	// 重新获取完整数据
	fullItem, err := p.presetRepo.GetPromptItemByID(ctx, item.ID)
	if err != nil {
		return nil, err
	}

	return &pb.AddPromptItemResponse{
		Item: convert.PromptItemEntityToPb(fullItem),
	}, nil
}

func (p *presetImpl) UpdatePromptItem(ctx context.Context, req *pb.UpdatePromptItemRequest) (*pb.UpdatePromptItemResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPromptItemID)
	}

	// 参数校验
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPromptItemName)
	}

	// 获取当前提示项
	item, err := p.presetRepo.GetPromptItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.PromptItemNotFound)
	}

	// 更新提示项字段
	item.Identifier = req.GetIdentifier()
	item.Name = name
	if req.Content != nil {
		item.Content = *req.Content
	}
	item.Role = req.Role
	item.IsEnabled = req.IsEnabled
	item.InjectionPosition = req.InjectionPosition
	item.InjectionDepth = int(req.InjectionDepth)
	item.ForbidOverrides = req.ForbidOverrides

	// 更新数据库
	if err := p.presetRepo.UpdatePromptItem(ctx, item); err != nil {
		return nil, err
	}

	// 使相关缓存失效（提示项修改会影响使用该预设的缓存）
	cache.InvalidateCacheByPreset(int64(item.PresetID))

	// 重新获取更新后的提示项
	updatedItem, err := p.presetRepo.GetPromptItemByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePromptItemResponse{
		Item: convert.PromptItemEntityToPb(updatedItem),
	}, nil
}

func (p *presetImpl) DeletePromptItem(ctx context.Context, req *pb.DeletePromptItemRequest) (*pb.DeletePromptItemResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPromptItemID)
	}

	// 删除提示项
	if err := p.presetRepo.DeletePromptItem(ctx, id); err != nil {
		return nil, err
	}

	return &pb.DeletePromptItemResponse{}, nil
}

func (p *presetImpl) UpdatePromptItemsOrder(ctx context.Context, req *pb.UpdatePromptItemsOrderRequest) (*pb.UpdatePromptItemsOrderResponse, error) {
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	if req.SourceId <= 0 || req.DesId <= 0 || req.SortOperation == pb.SortOperation_OrderOperationUnspecified {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptySortData)
	}

	userID := jwt.GetUserId(ctx)
	// 更新排序
	if err := p.presetRepo.UpdatePromptItemsOrder(ctx, presetID, userID, int(req.SourceId), int(req.DesId), req.SortOperation); err != nil {
		return nil, err
	}

	return &pb.UpdatePromptItemsOrderResponse{}, nil
}

func (p *presetImpl) SetActivePreset(ctx context.Context, req *pb.SetActivePresetRequest) (*pb.SetActivePresetResponse, error) {
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}
	if err := p.presetRepo.SetActivePreset(ctx, jwt.GetUserId(ctx), presetID); err != nil {
		return nil, err
	}
	return &pb.SetActivePresetResponse{}, nil
}

func (p *presetImpl) ImportPreset(ctx context.Context, req *pb.ImportPresetRequest) (*pb.ImportPresetResponse, error) {
	fileContent := req.GetFileContent()
	fileName := strings.TrimSpace(req.GetFileName())

	// 移除 UTF-8 BOM（如果存在）
	if len(fileContent) >= 3 && fileContent[0] == 0xEF && fileContent[1] == 0xBB && fileContent[2] == 0xBF {
		fileContent = fileContent[3:]
	}

	// 解析 SillyTavern 预设文件
	stPreset := &sillytavern.OpenAIPreset{}
	if err := json.Unmarshal(fileContent, stPreset); err != nil {
		return nil, errs.NewStandardf(connect.CodeInvalidArgument, "无效的预设文件: %v", err.Error())
	}

	// 生成预设名称：优先使用文件名，否则使用默认名称
	presetName := "导入的预设"
	if fileName != "" {
		// 移除文件扩展名
		presetName = strings.TrimSuffix(fileName, ".json")
	}

	// 使用 convert 包转换为 Muse 预设实体
	userId := jwt.GetUserId(ctx)
	preset := convert.STPresetToEntity(stPreset, userId, presetName)

	// 保存预设到数据库
	if err := p.presetRepo.Create(ctx, preset); err != nil {
		return nil, err
	}

	// 转换并保存提示项
	if len(stPreset.Prompts) > 0 {
		promptItems := orderPrompt(stPreset.Prompts, stPreset.PromptOrder)
		items := make([]*entity.PromptItem, 0, len(promptItems))
		for _, stPrompt := range promptItems {
			// 使用 convert 包转换提示项
			item := convert.STPromptToEntity(preset.ID, &stPrompt)
			items = append(items, item)
		}
		if err := p.presetRepo.BatchCreatePromptItem(ctx, items); err != nil {
			log.Errorf("failed to batch create prompt items: %v", err)
			return nil, err
		}
		log.Infof("batch created %d prompt items", len(items))
		// 排序提示项
		var pre *int = nil
		for i := 0; i < len(items); i++ {
			items[i].Pre = pre
			pre = &items[i].ID
			if i < len(items)-1 {
				items[i].Next = &items[i+1].ID
			}
			preLog := "nil"
			if items[i].Pre != nil {
				preLog = fmt.Sprintf("%d", *items[i].Pre)
			}
			nextLog := "nil"
			if items[i].Next != nil {
				nextLog = fmt.Sprintf("%d", *items[i].Next)
			}
			log.Debugf("当前提示项 ID: %d, Pre: %s, Next: %s", items[i].ID, preLog, nextLog)
		}
		if err := p.presetRepo.BatchUpdatePromptItemOrder(ctx, items); err != nil {
			log.Errorf("failed to batch update prompt item order: %v", err)
			return nil, err
		}
	}

	// TODO 导入预设内嵌的正则脚本
	if len(stPreset.Extensions.RegexScripts) > 0 {
		regexRules := make([]*entity.RegexRule, 0, len(stPreset.Extensions.RegexScripts))
		for _, stScript := range stPreset.Extensions.RegexScripts {
			// 使用 convert 包转换正则脚本，关联到新创建的预设（characterID=0 表示非角色范围正则）
			rule := convert.STRegexToEntity(&stScript, userId)
			regexRules = append(regexRules, rule)
		}

		// 批量创建正则规则
		if err := p.regexRuleRepo.BatchCreate(ctx, regexRules); err != nil {
			log.Errorf("failed to batch create regex rules: %v", err)
			return nil, errs.NewStandardf(errs.Code(err), "导入预设时批量创建正则规则失败: %v", err.Error())
		}
	}

	// 重新获取完整预设数据（包含关联的提示项）
	fullPreset, err := p.presetRepo.GetByID(ctx, preset.ID, userId)
	if err != nil {
		return nil, err
	}

	pbPreset, pbPromptItem, pbRegex := convert.PresetEntityToPb(fullPreset)
	return &pb.ImportPresetResponse{
		Preset: &pb.PresetWithAll{
			Preset:      pbPreset,
			PromptItems: pbPromptItem,
			RegexRules:  pbRegex,
		},
	}, nil
}

func (p *presetImpl) ExportPreset(ctx context.Context, req *pb.ExportPresetRequest) (*pb.ExportPresetResponse, error) {
	//TODO implement me
	panic("implement me")
}

// orderPrompt 构建提示词顺序映射
func orderPrompt(prompt []sillytavern.PresetPromptItem,
	promptOrder []sillytavern.PromptOrderItem) []sillytavern.PresetPromptItem {

	promptMap := make(map[string]*sillytavern.PresetPromptItem, len(prompt))
	for _, item := range prompt {
		promptMap[item.Identifier] = &item
	}
	orderPrompts := make([]sillytavern.PresetPromptItem, 0, len(prompt))

	const defaultCharacterID = 100001

	for _, order := range promptOrder {
		if order.CharacterID != defaultCharacterID {
			continue
		}
		for _, item := range order.Order {
			if promptItem, exist := promptMap[item.Identifier]; exist {
				promptItem.Enabled = item.Enabled // 注意这里特别坑, Enable放在了顺序这里, 所以得手动赋值到外层结构体
				orderPrompts = append(orderPrompts, *promptItem)
			}
		}
		break
	}
	// 如果没有找到默认角色，使用第一个配置
	if len(orderPrompts) == 0 && len(promptOrder) > 0 {
		log.Info("没有找到默认角色，使用第一个配置")
		for _, item := range promptOrder[0].Order {
			if promptItem, exist := promptMap[item.Identifier]; exist {
				promptItem.Enabled = item.Enabled // 注意这里特别坑, Enable放在了顺序这里, 所以得手动赋值到外层结构体
				orderPrompts = append(orderPrompts, *promptItem)
			}
		}
	}
	return orderPrompts
}
