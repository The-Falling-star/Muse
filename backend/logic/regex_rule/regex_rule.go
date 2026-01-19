package regex_rule

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/entity"
	"github.com/ling/muse/entity/sillytavern"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/database"
)

type regexRuleImpl struct {
	regexRuleRepo *database.RegexRuleRepo
}

func newRegexRule() *regexRuleImpl {
	return &regexRuleImpl{
		regexRuleRepo: &database.RegexRuleRepo{},
	}
}

func (r *regexRuleImpl) ListRegexRules(ctx context.Context, req *pb.ListRegexRulesRequest) (*pb.ListRegexRulesResponse, error) {
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	// 从数据库获取正则规则列表
	rules, err := r.regexRuleRepo.List(presetID)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbRules := make([]*pb.RegexRule, 0, len(rules))
	for _, rule := range rules {
		pbRules = append(pbRules, convert.RegexRuleEntityToPb(rule))
	}

	return &pb.ListRegexRulesResponse{
		Rules: pbRules,
	}, nil
}

func (r *regexRuleImpl) AddRegexRule(ctx context.Context, req *pb.AddRegexRuleRequest) (*pb.AddRegexRuleResponse, error) {
	// 参数校验
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyRegexRuleName)
	}

	findPattern := strings.TrimSpace(req.GetFindPattern())
	if findPattern == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyFindPattern)
	}

	affectFlags := req.GetAffectFlags()
	if affectFlags == nil {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyAffectFlags)
	}

	// 构建正则规则实体
	rule := &entity.RegexRule{
		PresetID:                presetID,
		Name:                    name,
		FindPattern:             findPattern,
		ReplacePattern:          req.GetReplacePattern(),
		MinDepth:                int(req.GetMinDepth()),
		MaxDepth:                int(req.GetMaxDepth()),
		IsEnabled:               req.IsEnabled,
		RunOnEdit:               req.RunOnEdit,
		SubstituteRegex:         req.SubstituteRegex,
		AffectFlagsUserInput:    affectFlags.UserInput,
		AffectFlagsAIOutput:     affectFlags.AiOutput,
		AffectFlagsSlashCommand: affectFlags.SlashCommand,
		AffectFlagsWorldInfo:    affectFlags.WorldInfo,
		AffectFlagsPrompt:       affectFlags.Prompt,
		SortOrder:               int(req.SortOrder),
	}

	// 保存到数据库
	if err := r.regexRuleRepo.Create(rule); err != nil {
		return nil, err
	}

	// 重新获取完整数据
	fullRule, err := r.regexRuleRepo.GetByID(rule.ID)
	if err != nil {
		return nil, err
	}

	return &pb.AddRegexRuleResponse{
		Rule: convert.RegexRuleEntityToPb(fullRule),
	}, nil
}

func (r *regexRuleImpl) UpdateRegexRule(ctx context.Context, req *pb.UpdateRegexRuleRequest) (*pb.UpdateRegexRuleResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidRegexRuleID)
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyRegexRuleName)
	}

	findPattern := strings.TrimSpace(req.GetFindPattern())
	if findPattern == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyFindPattern)
	}

	affectFlags := req.GetAffectFlags()
	if affectFlags == nil {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyAffectFlags)
	}

	// 获取当前规则
	rule, err := r.regexRuleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.RegexRuleNotFound)
	}

	// 更新规则字段
	rule.Name = name
	rule.FindPattern = findPattern
	rule.ReplacePattern = req.GetReplacePattern()
	rule.MinDepth = int(req.GetMinDepth())
	rule.MaxDepth = int(req.GetMaxDepth())
	rule.IsEnabled = req.IsEnabled
	rule.RunOnEdit = req.RunOnEdit
	rule.SubstituteRegex = req.SubstituteRegex
	rule.AffectFlagsUserInput = affectFlags.UserInput
	rule.AffectFlagsAIOutput = affectFlags.AiOutput
	rule.AffectFlagsSlashCommand = affectFlags.SlashCommand
	rule.AffectFlagsWorldInfo = affectFlags.WorldInfo
	rule.AffectFlagsPrompt = affectFlags.Prompt
	rule.SortOrder = int(req.SortOrder)

	// 更新数据库
	if err := r.regexRuleRepo.Update(rule); err != nil {
		return nil, err
	}

	// 重新获取更新后的规则
	updatedRule, err := r.regexRuleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateRegexRuleResponse{
		Rule: convert.RegexRuleEntityToPb(updatedRule),
	}, nil
}

func (r *regexRuleImpl) DeleteRegexRule(ctx context.Context, req *pb.DeleteRegexRuleRequest) (*pb.DeleteRegexRuleResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidRegexRuleID)
	}

	// 删除规则
	if err := r.regexRuleRepo.Delete(id); err != nil {
		return nil, err
	}

	return &pb.DeleteRegexRuleResponse{}, nil
}

func (r *regexRuleImpl) UpdateRegexRulesOrder(ctx context.Context, req *pb.UpdateRegexRulesOrderRequest) (*pb.UpdateRegexRulesOrderResponse, error) {
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPresetID)
	}

	if len(req.RuleIds) == 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptySortData)
	}

	// 构建排序映射
	ruleOrders := make(map[int]int)
	for i, ruleID := range req.RuleIds {
		ruleOrders[int(ruleID)] = i
	}

	// 更新排序
	if err := r.regexRuleRepo.UpdateRulesOrder(presetID, ruleOrders); err != nil {
		return nil, err
	}

	return &pb.UpdateRegexRulesOrderResponse{}, nil
}

func (r *regexRuleImpl) ImportRegexRules(ctx context.Context, req *pb.ImportRegexRulesRequest) (*pb.ImportRegexRulesResponse, error) {
	fileContent := req.GetFileContent()

	// 移除 UTF-8 BOM（如果存在）
	if len(fileContent) >= 3 && fileContent[0] == 0xEF && fileContent[1] == 0xBB && fileContent[2] == 0xBF {
		fileContent = fileContent[3:]
	}

	// 尝试解析为数组或单个对象
	var stScripts []sillytavern.RegexScript
	if err := json.Unmarshal(fileContent, &stScripts); err != nil {
		// 尝试解析为单个对象
		var singleScript sillytavern.RegexScript
		if err = json.Unmarshal(fileContent, &singleScript); err != nil {
			return nil, errs.NewStandardf(connect.CodeInvalidArgument, "无效的正则规则文件: %v", err)
		}
		stScripts = []sillytavern.RegexScript{singleScript}
	}

	if len(stScripts) == 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, "正则规则文件为空")
	}

	// 导入的都是全局正则（preset_id = 0）
	const presetID = 0

	// 获取当前最大排序号
	maxOrder, err := r.regexRuleRepo.GetMaxSortOrder(presetID)
	if err != nil {
		return nil, err
	}

	// 转换为 Muse 正则规则实体
	rules := make([]*entity.RegexRule, 0, len(stScripts))
	for i, stScript := range stScripts {
		if stScript.ScriptName == "" {
			continue // 跳过没有名称的脚本
		}

		rule := convert.STRegexToEntity(&stScript, presetID, 0, maxOrder+i+1)
		rules = append(rules, rule)
	}

	if len(rules) == 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, "没有有效的正则规则可导入")
	}

	// 批量创建
	if err = r.regexRuleRepo.BatchCreate(rules); err != nil {
		return nil, err
	}

	return &pb.ImportRegexRulesResponse{}, nil
}

func (r *regexRuleImpl) ExportRegexRules(ctx context.Context, req *pb.ExportRegexRulesRequest) (*pb.ExportRegexRulesResponse, error) {
	presetID := int(req.GetPresetId()) // 0 表示全局正则

	// 获取正则规则列表
	rules, err := r.regexRuleRepo.List(presetID)
	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return nil, errs.NewStandard(connect.CodeNotFound, "没有可导出的正则规则")
	}

	// 转换为 SillyTavern 格式
	stScripts := make([]sillytavern.RegexScript, 0, len(rules))
	for _, rule := range rules {
		stScripts = append(stScripts, convert.EntityToSTRegex(rule))
	}

	// 序列化为 JSON
	fileContent, jsonErr := json.MarshalIndent(stScripts, "", "  ")
	if jsonErr != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "序列化正则规则失败: %v", jsonErr)
	}

	// 生成文件名
	fileName := "regex-rules.json"
	if presetID > 0 {
		fileName = fmt.Sprintf("regex-rules-preset-%d.json", presetID)
	} else {
		fileName = "regex-rules-global.json"
	}

	return &pb.ExportRegexRulesResponse{
		FileContent: fileContent,
		FileName:    fileName,
	}, nil
}
