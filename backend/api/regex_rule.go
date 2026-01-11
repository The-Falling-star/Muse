package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/regex_rule"
)

// RegexRuleServer 正则规则服务
type RegexRuleServer struct {
	regexRule regex_rule.RegexRule
}

// NewRegexRuleServer 创建一个新的RegexRuleServer实例
func NewRegexRuleServer() *RegexRuleServer {
	return &RegexRuleServer{
		regexRule: regex_rule.NewRegexRule(),
	}
}

// ListRegexRules 获取正则规则列表
func (r *RegexRuleServer) ListRegexRules(ctx context.Context, req *connect.Request[pb.ListRegexRulesRequest]) (
	*connect.Response[pb.ListRegexRulesResponse], error) {
	resp, err := r.regexRule.ListRegexRules(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// AddRegexRule 添加新的正则规则
func (r *RegexRuleServer) AddRegexRule(ctx context.Context, req *connect.Request[pb.AddRegexRuleRequest]) (
	*connect.Response[pb.AddRegexRuleResponse], error) {
	resp, err := r.regexRule.AddRegexRule(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateRegexRule 更新指定正则规则
func (r *RegexRuleServer) UpdateRegexRule(ctx context.Context, req *connect.Request[pb.UpdateRegexRuleRequest]) (
	*connect.Response[pb.UpdateRegexRuleResponse], error) {
	resp, err := r.regexRule.UpdateRegexRule(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteRegexRule 删除指定正则规则
func (r *RegexRuleServer) DeleteRegexRule(ctx context.Context, req *connect.Request[pb.DeleteRegexRuleRequest]) (
	*connect.Response[pb.DeleteRegexRuleResponse], error) {
	resp, err := r.regexRule.DeleteRegexRule(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateRegexRulesOrder 更新正则规则的排序
func (r *RegexRuleServer) UpdateRegexRulesOrder(ctx context.Context,
	req *connect.Request[pb.UpdateRegexRulesOrderRequest]) (
	*connect.Response[pb.UpdateRegexRulesOrderResponse], error) {
	resp, err := r.regexRule.UpdateRegexRulesOrder(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
