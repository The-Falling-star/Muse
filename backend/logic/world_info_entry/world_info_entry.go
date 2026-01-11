package world_info_entry

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type worldInfoEntryImpl struct {
}

func newWorldInfoEntry() *worldInfoEntryImpl {
	return &worldInfoEntryImpl{}
}

func (w *worldInfoEntryImpl) ListWorldInfoEntries(ctx context.Context, req *pb.ListWorldInfoEntriesRequest) (*pb.ListWorldInfoEntriesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoEntryImpl) AddWorldInfoEntry(ctx context.Context, req *pb.AddWorldInfoEntryRequest) (*pb.AddWorldInfoEntryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoEntryImpl) UpdateWorldInfoEntry(ctx context.Context, req *pb.UpdateWorldInfoEntryRequest) (*pb.UpdateWorldInfoEntryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoEntryImpl) DeleteWorldInfoEntry(ctx context.Context, req *pb.DeleteWorldInfoEntryRequest) (*pb.DeleteWorldInfoEntryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoEntryImpl) UpdateWorldInfoEntriesOrder(ctx context.Context, req *pb.UpdateWorldInfoEntriesOrderRequest) (*pb.UpdateWorldInfoEntriesOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}
