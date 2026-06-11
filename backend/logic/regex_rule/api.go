package regex_rule

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// RegexRule 定义了正则规则服务的接口
type RegexRule interface {
	// ListRegexRules 获取正则规则列表（全量拉取全局+预设+角色正则）
	ListRegexRules(ctx context.Context, req *pb.ListRegexRulesReq) (*pb.ListRegexRulesRsp, error)
	// AddRegexRule 添加新的正则规则
	AddRegexRule(ctx context.Context, req *pb.AddRegexRuleReq) (*pb.AddRegexRuleRsp, error)
	// UpdateRegexRule 更新指定正则规则
	UpdateRegexRule(ctx context.Context, req *pb.UpdateRegexRuleReq) (*pb.UpdateRegexRuleRsp, error)
	// DeleteRegexRule 删除指定正则规则
	DeleteRegexRule(ctx context.Context, req *pb.DeleteRegexRuleReq) (*pb.DeleteRegexRuleRsp, error)
	// UpdateRegexRulesOrder 更新正则规则的排序
	UpdateRegexRulesOrder(ctx context.Context, req *pb.UpdateRegexRulesOrderReq) (*pb.UpdateRegexRulesOrderRsp, error)
	// ImportRegexRules 导入正则规则
	ImportRegexRules(ctx context.Context, req *pb.ImportRegexRulesReq) (*pb.ImportRegexRulesRsp, error)
	// ExportRegexRules 导出正则规则
	ExportRegexRules(ctx context.Context, req *pb.ExportRegexRulesReq) (*pb.ExportRegexRulesRsp, error)
}

// NewRegexRule 创建一个新的RegexRule实例
func NewRegexRule() RegexRule {
	return newRegexRule()
}
