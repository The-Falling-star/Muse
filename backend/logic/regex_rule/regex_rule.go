package regex_rule

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
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	// 从数据库获取正则规则列表
	rules, err := r.regexRuleRepo.List(presetID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("规则名称不能为空"))
	}

	findPattern := strings.TrimSpace(req.GetFindPattern())
	if findPattern == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("查找模式不能为空"))
	}

	affectFlags := req.GetAffectFlags()
	if affectFlags == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("影响标志不能为空"))
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
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取完整数据
	fullRule, err := r.regexRuleRepo.GetByID(rule.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.AddRegexRuleResponse{
		Rule: convert.RegexRuleEntityToPb(fullRule),
	}, nil
}

func (r *regexRuleImpl) UpdateRegexRule(ctx context.Context, req *pb.UpdateRegexRuleRequest) (*pb.UpdateRegexRuleResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的规则ID"))
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("规则名称不能为空"))
	}

	findPattern := strings.TrimSpace(req.GetFindPattern())
	if findPattern == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("查找模式不能为空"))
	}

	affectFlags := req.GetAffectFlags()
	if affectFlags == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("影响标志不能为空"))
	}

	// 获取当前规则
	rule, err := r.regexRuleRepo.GetByID(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if rule == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("规则不存在"))
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
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取更新后的规则
	updatedRule, err := r.regexRuleRepo.GetByID(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdateRegexRuleResponse{
		Rule: convert.RegexRuleEntityToPb(updatedRule),
	}, nil
}

func (r *regexRuleImpl) DeleteRegexRule(ctx context.Context, req *pb.DeleteRegexRuleRequest) (*pb.DeleteRegexRuleResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的规则ID"))
	}

	// 删除规则
	if err := r.regexRuleRepo.Delete(id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.DeleteRegexRuleResponse{}, nil
}

func (r *regexRuleImpl) UpdateRegexRulesOrder(ctx context.Context, req *pb.UpdateRegexRulesOrderRequest) (*pb.UpdateRegexRulesOrderResponse, error) {
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	if len(req.RuleIds) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("排序数据不能为空"))
	}

	// 构建排序映射
	ruleOrders := make(map[int]int)
	for i, ruleID := range req.RuleIds {
		ruleOrders[int(ruleID)] = i
	}

	// 更新排序
	if err := r.regexRuleRepo.UpdateRulesOrder(presetID, ruleOrders); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdateRegexRulesOrderResponse{}, nil
}
