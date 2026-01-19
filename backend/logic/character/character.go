package character

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	png "github.com/dsoprea/go-png-image-structure/v2"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/entity"
	"github.com/ling/muse/entity/sillytavern"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/database"
)

// 默认用户ID，待认证功能完成后替换
const defaultUserID = 1

type characterImpl struct {
	charaRepo     *database.CharacterRepo
	regexRuleRepo *database.RegexRuleRepo
}

func newCharacter() *characterImpl {
	return &characterImpl{
		charaRepo:     &database.CharacterRepo{},
		regexRuleRepo: &database.RegexRuleRepo{},
	}
}

func (c *characterImpl) ListCharacters(ctx context.Context, req *pb.ListCharactersRequest) (*pb.ListCharactersResponse, error) {
	// 获取分页参数
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 从数据库获取角色列表
	characters, total, err := c.charaRepo.List(defaultUserID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbCharacters := make([]*pb.Character, 0, len(characters))
	for _, character := range characters {
		pbCharacters = append(pbCharacters, convert.CharaEntityToPb(character))
	}

	return &pb.ListCharactersResponse{
		Characters: pbCharacters,
		Total:      int32(total),
	}, nil
}

func (c *characterImpl) GetCharacter(ctx context.Context, req *pb.GetCharacterRequest) (*pb.GetCharacterResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidCharacterID)
	}

	// 从数据库获取角色
	character, err := c.charaRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.CharacterNotFound)
	}

	return &pb.GetCharacterResponse{
		Character: convert.CharaEntityToPb(character),
	}, nil
}

func (c *characterImpl) CreateCharacter(ctx context.Context, req *pb.CreateCharacterRequest) (*pb.CreateCharacterResponse, error) {
	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyCharacterName)
	}

	// 构建角色实体
	character := &entity.Character{
		UserID:          defaultUserID,
		Name:            name,
		Avatar:          req.GetAvatar(),
		Description:     req.GetDescription(),
		FirstMessage:    req.GetFirstMessage(),
		ExampleDialogue: req.GetExampleDialogue(),
		CreatorNotes:    req.GetCreatorNotes(),
		Version:         1,
	}

	// 设置关联的世界信息ID
	if req.WorldInfoId != nil {
		character.WorldInfoID = int(*req.WorldInfoId)
	}

	// 保存到数据库
	if err := c.charaRepo.Create(character); err != nil {
		return nil, err
	}

	return &pb.CreateCharacterResponse{
		Character: convert.CharaEntityToPb(character),
	}, nil
}

func (c *characterImpl) UpdateCharacter(ctx context.Context, req *pb.UpdateCharacterRequest) (*pb.UpdateCharacterResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidCharacterID)
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyCharacterName)
	}

	// 获取当前角色
	character, err := c.charaRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.CharacterNotFound)
	}

	// 检查版本号
	if int64(character.Version) != req.GetVersion() {
		return nil, errs.NewStandard(connect.CodeAborted, errs.DataConflict)
	}

	// 更新角色字段
	character.Name = name
	if req.Avatar != nil {
		character.Avatar = *req.Avatar
	}
	if req.Description != nil {
		character.Description = *req.Description
	}
	if req.FirstMessage != nil {
		character.FirstMessage = *req.FirstMessage
	}
	if req.ExampleDialogue != nil {
		character.ExampleDialogue = *req.ExampleDialogue
	}
	if req.CreatorNotes != nil {
		character.CreatorNotes = *req.CreatorNotes
	}
	if req.WorldInfoId != nil {
		character.WorldInfoID = int(*req.WorldInfoId)
	}

	// 更新数据库
	if err := c.charaRepo.Update(character); err != nil {
		return nil, err
	}

	// 重新获取更新后的角色（包含关联数据）
	updatedCharacter, err := c.charaRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateCharacterResponse{
		Character: convert.CharaEntityToPb(updatedCharacter),
	}, nil
}

func (c *characterImpl) DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterRequest) (*pb.DeleteCharacterResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidCharacterID)
	}

	// 删除角色
	if err := c.charaRepo.Delete(id, defaultUserID); err != nil {
		return nil, err
	}

	return &pb.DeleteCharacterResponse{}, nil
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
			return nil, errs.NewStandardf(connect.CodeInvalidArgument, "无法解析PNG文件: %v", err)
		}
	case ".json":
		// 直接解析JSON格式
		card = &sillytavern.CharacterCard{}
		if err := json.Unmarshal(fileContent, card); err != nil {
			return nil, errs.NewStandardf(connect.CodeInvalidArgument, "无法解析JSON文件: %v", err)
		}
	default:
		return nil, errs.NewStandardf(connect.CodeInvalidArgument, "不支持的文件格式: %s", ext)
	}

	worldBook := getWorldBookFromCard(card)
	var buf bytes.Buffer
	zip := gzip.NewWriter(&buf)
	defer zip.Close()
	worldBookJson, _ := json.Marshal(worldBook)
	if _, err = zip.Write(worldBookJson); err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "压缩世界书失败： %v", err)
	}
	character := &entity.Character{
		UserID:          defaultUserID,
		Name:            card.Name,
		Description:     card.Description,
		FirstMessage:    card.FirstMes,
		ExampleDialogue: card.MesExample,
		CreatorNotes:    card.CreatorNotes,
		WorldInfoBackup: buf.Bytes(),
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

	if err := c.charaRepo.Create(character); err != nil {
		return nil, err
	}

	// 导入角色卡内嵌的正则脚本（Scoped Scripts）
	if card.Extensions != nil && len(card.Extensions.RegexScripts) > 0 {
		regexRules := make([]*entity.RegexRule, 0, len(card.Extensions.RegexScripts))
		for i, stScript := range card.Extensions.RegexScripts {
			// 使用 convert 包转换正则脚本，关联到新创建的角色（presetID=0 表示非预设正则）
			rule := convert.STRegexToEntity(&stScript, 0, character.ID, i)
			regexRules = append(regexRules, rule)
		}

		// 批量创建正则规则
		if err = c.regexRuleRepo.BatchCreate(regexRules); err != nil {
			// 即使正则规则创建失败，也不影响角色卡的导入
			// 可以记录日志但不返回错误
		}
	}

	return &pb.ImportCharacterResponse{
		Character: convert.CharaEntityToPb(character),
	}, nil
}

func getWorldBookFromCard(card *sillytavern.CharacterCard) entity.WorldInfo {
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
			// Position:       entry.Position, TODO: 需要将字符串位置转换为枚举
			Depth:     entry.Depth,
			SortOrder: i,
		})
	}
	return worldBook
}

func (c *characterImpl) ExportCharacter(ctx context.Context, req *pb.ExportCharacterRequest) (*pb.ExportCharacterResponse, error) {
	character, err := c.charaRepo.GetByID(int(req.GetId()), defaultUserID)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.CharacterNotFound)
	}

	// 将角色数据转换为CharacterCard
	card := &sillytavern.CharacterCard{
		Name:         character.Name,
		Description:  character.Description,
		FirstMes:     character.FirstMessage,
		MesExample:   character.ExampleDialogue,
		CreatorNotes: character.CreatorNotes,
	}

	// 获取原始头像PNG数据
	// 如果头像是data URI格式，需要解码
	var avatarPNG []byte
	if character.Avatar != "" && strings.HasPrefix(character.Avatar, "data:image/png;base64,") {
		avatarPNG, _ = base64.StdEncoding.DecodeString(strings.TrimPrefix(character.Avatar, "data:image/png;base64,"))
	}

	// 如果没有头像，使用默认空白PNG
	if len(avatarPNG) == 0 {
		// 创建一个1x1的透明PNG作为默认头像
		// 这是一个最小的有效PNG文件
		avatarPNG, _ = base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==")
	}

	// 将角色卡数据写入PNG
	pngData, pngErr := WriteToPNG(avatarPNG, card)
	if pngErr != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "写入角色卡到PNG失败: %v", pngErr)
	}

	fileName := fmt.Sprintf("%s.png", character.Name)
	return &pb.ExportCharacterResponse{
		FileContent: pngData,
		FileName:    fileName,
	}, nil
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
	if err = newCs.WriteTo(&buf); err != nil {
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
