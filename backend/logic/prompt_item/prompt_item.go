package prompt_item

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type promptItemImpl struct {
}

func newPromptItem() *promptItemImpl {
	return &promptItemImpl{}
}

func (p *promptItemImpl) ListPromptItems(ctx context.Context, req *pb.ListPromptItemsRequest) (*pb.ListPromptItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *promptItemImpl) AddPromptItem(ctx context.Context, req *pb.AddPromptItemRequest) (*pb.AddPromptItemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *promptItemImpl) UpdatePromptItem(ctx context.Context, req *pb.UpdatePromptItemRequest) (*pb.UpdatePromptItemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *promptItemImpl) DeletePromptItem(ctx context.Context, req *pb.DeletePromptItemRequest) (*pb.DeletePromptItemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *promptItemImpl) UpdatePromptItemsOrder(ctx context.Context, req *pb.UpdatePromptItemsOrderRequest) (*pb.UpdatePromptItemsOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}
