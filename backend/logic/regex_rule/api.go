package regex_rule

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// RegexRule 定义了正则规则服务的接口
type RegexRule interface {
	// ListRegexRules 获取正则规则列表
	ListRegexRules(ctx context.Context, req *pb.ListRegexRulesRequest) (*pb.ListRegexRulesResponse, error)
	// AddRegexRule 添加新的正则规则
	AddRegexRule(ctx context.Context, req *pb.AddRegexRuleRequest) (*pb.AddRegexRuleResponse, error)
	// UpdateRegexRule 更新指定正则规则
	UpdateRegexRule(ctx context.Context, req *pb.UpdateRegexRuleRequest) (*pb.UpdateRegexRuleResponse, error)
	// DeleteRegexRule 删除指定正则规则
	DeleteRegexRule(ctx context.Context, req *pb.DeleteRegexRuleRequest) (*pb.DeleteRegexRuleResponse, error)
	// UpdateRegexRulesOrder 更新正则规则的排序
	UpdateRegexRulesOrder(ctx context.Context, req *pb.UpdateRegexRulesOrderRequest) (*pb.UpdateRegexRulesOrderResponse, error)
	// ImportRegexRules 导入正则规则
	ImportRegexRules(ctx context.Context, req *pb.ImportRegexRulesRequest) (*pb.ImportRegexRulesResponse, error)
	// ExportRegexRules 导出正则规则
	ExportRegexRules(ctx context.Context, req *pb.ExportRegexRulesRequest) (*pb.ExportRegexRulesResponse, error)
}

// NewRegexRule 创建一个新的RegexRule实例
func NewRegexRule() RegexRule {
	return newRegexRule()
}
