package character

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	png "github.com/dsoprea/go-png-image-structure/v2"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/entity"
	"github.com/ling/muse/entity/sillytavern"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/database"
)

type characterImpl struct {
	charaRepo *database.CharacterRepo
}

func newCharacter() *characterImpl {
	return &characterImpl{
		charaRepo: &database.CharacterRepo{},
	}
}

func (c *characterImpl) ListCharacters(ctx context.Context, req *pb.ListCharactersRequest) (*pb.ListCharactersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) GetCharacter(ctx context.Context, req *pb.GetCharacterRequest) (*pb.GetCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) CreateCharacter(ctx context.Context, req *pb.CreateCharacterRequest) (*pb.CreateCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) UpdateCharacter(ctx context.Context, req *pb.UpdateCharacterRequest) (*pb.UpdateCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterRequest) (*pb.DeleteCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) ImportCharacter(ctx context.Context, req *pb.ImportCharacterRequest) (*pb.ImportCharacterResponse, error) {
	fileContent := req.GetFileContent()
	fileName := req.GetFileName()

	var card *sillytavern.CharacterCard
	var err error

	// 根据文件扩展名判断格式
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".png":
		// 从PNG图片中解析角色卡
		card, err = readFromPNG(fileContent)
		if err != nil {
			return nil, err
		}
	case ".json":
		// 直接解析JSON格式
		card = &sillytavern.CharacterCard{}
		if err := json.Unmarshal(fileContent, card); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	// 将角色卡数据转换为entity.Character

	worldBook := entity.WorldInfo{
		Name:        card.CharacterBook.Name,
		Description: card.CharacterBook.Description,
		IsGlobal:    false,
		Entries:     nil,
	}

	for i, entry := range card.CharacterBook.Entries {
		worldBook.Entries = append(worldBook.Entries, entity.WorldInfoEntry{
			KeysList:       strings.Join(entry.Keys, ","),
			SecondaryKeys:  strings.Join(entry.SecondaryKeys, ","),
			Content:        entry.Content,
			Comment:        entry.Comment,
			IsEnabled:      entry.Enabled,
			Constant:       entry.Constant,
			Selective:      entry.Selective,
			InsertionOrder: entry.InsertionOrder,
			// Position:       entry.Position, TODO
			Depth:     entry.Depth,
			SortOrder: i,
		})
	}

	character := &entity.Character{
		Name:            card.Name,
		Description:     card.Description,
		FirstMessage:    card.FirstMes,
		ExampleDialogue: card.MesExample,
		CreatorNotes:    card.CreatorNotes,
		WorldInfo:       worldBook,
	}

	// 如果有图片数据（从V3 assets中提取），设置头像
	if len(card.Assets) > 0 {
		for _, asset := range card.Assets {
			if asset.Type == "icon" && asset.URI != "" {
				// 如果是data URI，直接使用；否则可能需要下载
				character.Avatar = asset.URI
				break
			}
		}
	}

	// 如果PNG本身作为头像，将其转为base64 data URI
	if character.Avatar == "" && ext == ".png" {
		character.Avatar = "data:image/png;base64," + base64.StdEncoding.EncodeToString(fileContent)
	}

	if err = c.charaRepo.InsertChara(character); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.ImportCharacterResponse{
		Character: convert.CharaEntityToPb(character),
	}, nil
}

func (c *characterImpl) ExportCharacter(ctx context.Context, req *pb.ExportCharacterRequest) (*pb.ExportCharacterResponse, error) {
	// TODO: 从数据库获取角色数据
	// character, err := c.repo.GetCharacter(ctx, req.GetId())
	// if err != nil {
	//     return nil, err
	// }

	// 目前先返回未实现错误，等数据库层完成后再补充
	// 下面是导出逻辑的示例代码，展示如何使用 charcard.WriteToPNG

	/*
		// 将角色数据转换为CharacterCard
		card := &charcard.SillyTavernCharacterCard{
			Name:        character.Name,
			Description: derefString(character.Description),
			Personality: derefString(character.Personality),
			Scenario:    derefString(character.Scenario),
			FirstMes:    derefString(character.FirstMessage),
			MesExample:  derefString(character.ExampleDialogue),
			CreatorNotes: derefString(character.CreatorNotes),
			SystemPrompt: derefString(character.SystemPrompt),
		}

		// 获取原始头像PNG数据
		// 如果头像是data URI格式，需要解码
		var avatarPNG []byte
		if character.Avatar != nil && strings.HasPrefix(*character.Avatar, "data:image/png;base64,") {
			avatarPNG, _ = base64.StdEncoding.DecodeString(strings.TrimPrefix(*character.Avatar, "data:image/png;base64,"))
		}

		// 如果没有头像，使用默认空白PNG
		if len(avatarPNG) == 0 {
			avatarPNG = defaultPNG // 需要提供一个默认PNG
		}

		// 将角色卡数据写入PNG
		pngData, err := charcard.WriteToPNG(avatarPNG, card)
		if err != nil {
			return nil, err
		}

		fileName := fmt.Sprintf("%s.png", character.Name)
		return &pb.ExportCharacterResponse{
			FileContent: pngData,
			FileName:    fileName,
		}, nil
	*/

	panic("implement me: waiting for database layer")
}

func (c *characterImpl) RestoreCharacterWorldInfo(ctx context.Context, req *pb.RestoreCharacterWorldInfoRequest) (*pb.RestoreCharacterWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

// readFromPNG 从PNG文件字节数据中读取角色卡信息
func readFromPNG(data []byte) (*sillytavern.CharacterCard, error) {
	// 使用库解析PNG
	pmp := png.NewPngMediaParser()
	mc, err := pmp.ParseBytes(data)
	if err != nil {
		return nil, fmt.Errorf("invalid PNG file: signature mismatch: %v", err)
	}

	// 类型断言为ChunkSlice
	cs, ok := mc.(*png.ChunkSlice)
	if !ok {
		return nil, fmt.Errorf("invalid PNG file: not a chunk slice")
	}

	// 获取所有Chunk
	chunks := cs.Chunks()

	// 在tEXt块中查找角色卡数据
	var charaData, ccv3Data string
	for _, chunk := range chunks {
		if chunk.Type != "tEXt" {
			continue
		}

		keyword, text := parseTextChunk(chunk.Data)
		switch strings.ToLower(keyword) {
		case "ccv3":
			ccv3Data = text
		case "chara":
			charaData = text
		}
	}

	// 优先使用V3格式（ccv3），否则使用V2格式（chara）
	textData := ccv3Data
	if textData == "" {
		textData = charaData
	}
	if textData == "" {
		return nil, fmt.Errorf("no character card data found")
	}

	// Base64解码
	decoded, err := base64.StdEncoding.DecodeString(textData)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 data: %v", err)
	}

	// JSON解析
	var card sillytavern.CharacterCard
	if err := json.Unmarshal(decoded, &card); err != nil {
		return nil, fmt.Errorf("invalid JSON data: %v", err)
	}

	return &card, nil
}

// WriteToPNG 将角色卡数据写入PNG图片
func WriteToPNG(originalPNG []byte, card *sillytavern.CharacterCard) ([]byte, error) {
	// 使用库解析PNG
	pmp := png.NewPngMediaParser()
	mc, err := pmp.ParseBytes(originalPNG)
	if err != nil {
		return nil, fmt.Errorf("invalid PNG file: signature mismatch: %v", err)
	}

	// 类型断言为ChunkSlice
	cs, ok := mc.(*png.ChunkSlice)
	if !ok {
		return nil, fmt.Errorf("invalid PNG file: not a chunk slice")
	}

	// 获取所有Chunk并过滤掉旧的角色卡数据
	chunks := cs.Chunks()
	var filteredChunks []*png.Chunk

	for _, chunk := range chunks {
		// 过滤掉旧的chara和ccv3 tEXt块
		if chunk.Type == "tEXt" {
			keyword, _ := parseTextChunk(chunk.Data)
			if strings.ToLower(keyword) == "chara" || strings.ToLower(keyword) == "ccv3" {
				continue
			}
		}
		filteredChunks = append(filteredChunks, chunk)
	}

	// 将角色卡序列化为JSON
	jsonData, err := json.Marshal(card)
	if err != nil {
		return nil, err
	}

	// Base64编码
	base64Data := base64.StdEncoding.EncodeToString(jsonData)

	// 创建新的tEXt块数据：keyword + NULL + text
	textData := append([]byte("chara"), 0)
	textData = append(textData, []byte(base64Data)...)

	// 创建新的tEXt Chunk
	newChunk := &png.Chunk{
		Type:   "tEXt",
		Data:   textData,
		Length: uint32(len(textData)),
	}
	// 更新CRC校验
	newChunk.UpdateCrc32()

	// 将新Chunk插入到IHDR之后、IDAT之前
	var resultChunks []*png.Chunk
	inserted := false
	for _, chunk := range filteredChunks {
		if !inserted && (chunk.Type == "IDAT" || chunk.Type == "IEND") {
			resultChunks = append(resultChunks, newChunk)
			inserted = true
		}
		resultChunks = append(resultChunks, chunk)
	}

	// 如果还没插入（极端情况），添加到末尾
	if !inserted {
		resultChunks = append(resultChunks, newChunk)
	}

	// 使用库重新构建PNG
	newCs := png.NewChunkSlice(resultChunks)

	// 写入到buffer
	var buf bytes.Buffer
	if err := newCs.WriteTo(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// parseTextChunk 解析tEXt块的数据
func parseTextChunk(data []byte) (keyword, text string) {
	for i, b := range data {
		if b == 0 { // 找到NULL分隔符
			return string(data[:i]), string(data[i+1:])
		}
	}
	// 没有找到NULL分隔符，整个数据作为keyword
	return string(data), ""
}

// HasCharacterData 检查PNG数据中是否包含角色卡信息
func HasCharacterData(data []byte) bool {
	card, err := readFromPNG(data)
	return err == nil && card != nil && card.Name != ""
}
