package convert

import (
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
)

// CharaEntityToPb 将角色卡实体类转换为pb
func CharaEntityToPb(character *entity.Character) *pb.Character {
	worldInfoEntries := make([]*pb.WorldInfoEntry, len(character.WorldInfo.Entries))
	for _, entry := range character.WorldInfo.Entries {
		worldInfoEntries = append(worldInfoEntries, &pb.WorldInfoEntry{
			Id:             int32(entry.ID),
			WorldInfoId:    int32(entry.WorldInfoID),
			Uid:            entry.UID,
			KeysList:       entry.KeysList,
			SecondaryKeys:  entry.SecondaryKeys,
			Content:        entry.Content,
			Comment:        entry.Comment,
			IsEnabled:      entry.IsEnabled,
			Constant:       entry.Constant,
			Selective:      entry.Selective,
			InsertionOrder: int32(entry.InsertionOrder),
			Position:       entry.Position,
			Depth:          int32(entry.Depth),
			SortOrder:      int32(entry.SortOrder),
			CreatedAt:      entry.CreatedAt.Unix(),
			UpdatedAt:      entry.UpdatedAt.Unix(),
		})
	}
	worldInfo := &pb.WorldInfo{
		Id:          int32(character.WorldInfo.ID),
		UserId:      int32(character.UserID),
		Name:        character.WorldInfo.Name,
		Description: character.WorldInfo.Description,
		IsGlobal:    character.WorldInfo.IsGlobal,
		CreatedAt:   character.WorldInfo.CreatedAt.Unix(),
		UpdatedAt:   character.WorldInfo.UpdatedAt.Unix(),
		Entries:     worldInfoEntries,
	}
	chara := &pb.Character{
		Id:              int32(character.ID),
		UserId:          int32(character.UserID),
		Name:            character.Name,
		Avatar:          character.Avatar,
		Description:     character.Description,
		FirstMessage:    character.FirstMessage,
		ExampleDialogue: character.ExampleDialogue,
		CreatorNotes:    character.CreatorNotes,
		WorldInfoId:     int32(character.WorldInfo.ID),
		Version:         int64(character.Version),
		CreatedAt:       character.CreatedAt.Unix(),
		UpdatedAt:       character.UpdatedAt.Unix(),
		WorldInfo:       worldInfo,
	}
	return chara
}
