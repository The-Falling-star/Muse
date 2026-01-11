package world_info_entry

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type WorldInfoEntry interface {
	ListWorldInfoEntries(ctx context.Context, req *pb.ListWorldInfoEntriesRequest) (*pb.ListWorldInfoEntriesResponse, error)
	AddWorldInfoEntry(ctx context.Context, req *pb.AddWorldInfoEntryRequest) (*pb.AddWorldInfoEntryResponse, error)
	UpdateWorldInfoEntry(ctx context.Context, req *pb.UpdateWorldInfoEntryRequest) (*pb.UpdateWorldInfoEntryResponse, error)
	DeleteWorldInfoEntry(ctx context.Context, req *pb.DeleteWorldInfoEntryRequest) (*pb.DeleteWorldInfoEntryResponse, error)
	UpdateWorldInfoEntriesOrder(ctx context.Context, req *pb.UpdateWorldInfoEntriesOrderRequest) (*pb.UpdateWorldInfoEntriesOrderResponse, error)
}

func NewWorldInfoEntry() WorldInfoEntry {
	return newWorldInfoEntry()
}
