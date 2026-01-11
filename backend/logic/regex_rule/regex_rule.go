package regex_rule

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type regexRuleImpl struct {
}

func newRegexRule() *regexRuleImpl {
	return &regexRuleImpl{}
}

func (r *regexRuleImpl) ListRegexRules(ctx context.Context, req *pb.ListRegexRulesRequest) (*pb.ListRegexRulesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *regexRuleImpl) AddRegexRule(ctx context.Context, req *pb.AddRegexRuleRequest) (*pb.AddRegexRuleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *regexRuleImpl) UpdateRegexRule(ctx context.Context, req *pb.UpdateRegexRuleRequest) (*pb.UpdateRegexRuleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *regexRuleImpl) DeleteRegexRule(ctx context.Context, req *pb.DeleteRegexRuleRequest) (*pb.DeleteRegexRuleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *regexRuleImpl) UpdateRegexRulesOrder(ctx context.Context, req *pb.UpdateRegexRulesOrderRequest) (*pb.UpdateRegexRulesOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}
