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
	"github.com/ling/muse/common/constant"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/common/fileutil"
	"github.com/ling/muse/common/jwt"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"github.com/ling/muse/entity/sillytavern"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/file"
	"github.com/ling/muse/repo/cache"
	"github.com/ling/muse/repo/database"
	log "github.com/sirupsen/logrus"
)

type characterImpl struct {
	charRepo      database.CharacterRepository
	chatRepo      database.ChatRepository
	regexRuleRepo database.RegexRuleRepository
}

func newCharacter() *characterImpl {
	return &characterImpl{
		charRepo:      database.NewCharacterRepo(),
		chatRepo:      database.NewChatRepo(),
		regexRuleRepo: database.NewRegexRuleRepo(),
	}
}

func (c *characterImpl) ListCharacters(ctx context.Context, req *pb.ListCharactersReq) (*pb.ListCharactersRsp, error) {
	// 获取并规范化分页参数
	page, pageSize := constant.NormalizePagination(int(req.GetPage()), int(req.GetPageSize()))

	// 从数据库获取角色列表
	userId := jwt.GetUserId(ctx)
	characters, total, err := c.charRepo.List(ctx, userId, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbCharacters := make([]*pb.Character, 0, len(characters))
	for _, character := range characters {
		pbCharacters = append(pbCharacters, convert.CharaEntityToPb(character))
	}

	return &pb.ListCharactersRsp{
		Characters: pbCharacters,
		Total:      int32(total),
	}, nil
}

func (c *characterImpl) GetCharacter(ctx context.Context, req *pb.GetCharacterReq) (*pb.GetCharacterRsp, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidCharacterID)
	}

	// 从数据库获取角色
	userId := jwt.GetUserId(ctx)
	character, err := c.charRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.CharacterNotFound)
	}

	return &pb.GetCharacterRsp{
		Character: convert.CharaEntityToPb(character),
	}, nil
}

func (c *characterImpl) CreateCharacter(ctx context.Context, req *pb.CreateCharacterReq) (*pb.CreateCharacterRsp, error) {
	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyCharacterName)
	}

	// 构建角色实体
	userId := jwt.GetUserId(ctx)
	character := &entity.Character{
		UserID:          userId,
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
	if err := c.charRepo.Create(ctx, character); err != nil {
		return nil, err
	}

	return &pb.CreateCharacterRsp{
		Character: convert.CharaEntityToPb(character),
	}, nil
}

func (c *characterImpl) UpdateCharacter(ctx context.Context, req *pb.UpdateCharacterReq) (*pb.UpdateCharacterRsp, error) {
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
	userId := jwt.GetUserId(ctx)
	character, err := c.charRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, errs.NewStandardf(errs.Code(err), "获取角色失败: %v", err)
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
		character.FirstMessage = req.FirstMessage
	}
	if req.ExampleDialogue != nil {
		character.ExampleDialogue = req.ExampleDialogue
	}
	if req.CreatorNotes != nil {
		character.CreatorNotes = *req.CreatorNotes
	}
	if req.WorldInfoId != nil {
		character.WorldInfoID = int(*req.WorldInfoId)
	}

	// 更新数据库
	if err := c.charRepo.Update(ctx, character); err != nil {
		return nil, err
	}

	// 使相关缓存失效
	cache.InvalidateCacheByCharacter(int64(id))

	// 重新获取更新后的角色（包含关联数据）
	updatedCharacter, err := c.charRepo.GetByID(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateCharacterRsp{
		Character: convert.CharaEntityToPb(updatedCharacter),
	}, nil
}

func (c *characterImpl) DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterReq) (*pb.DeleteCharacterRsp, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidCharacterID)
	}

	// 使相关缓存失效（在删除之前）
	cache.InvalidateCacheByCharacter(int64(id))

	// 删除角色
	userId := jwt.GetUserId(ctx)
	if err := c.charRepo.Delete(ctx, id, userId); err != nil {
		return nil, err
	}
	return &pb.DeleteCharacterRsp{}, nil
}

func (c *characterImpl) ImportCharacter(ctx context.Context, req *pb.ImportCharacterReq) (*pb.ImportCharacterRsp, error) {
	fileContent := req.GetFileContent()
	fileName := req.GetFileName()

	var card *sillytavern.CharacterCard
	var err error

	// 根据文件扩展名判断格式
	ext := strings.ToLower(filepath.Ext(fileName))
	var pngData []byte
	var avatarData []byte // 用于存储解码后的头像图片数据
	switch ext {
	case ".png":
		// 从PNG图片中解析角色卡
		card, pngData, err = readFromPNG(fileContent)
		if err != nil {
			return nil, errs.NewStandardf(connect.CodeInvalidArgument, "无法解析PNG文件: %v", err)
		}
		// 如果角色卡没有头像且是PNG文件，使用PNG图片作为头像
		if (card.Avatar == "" || card.Avatar == "none") && len(pngData) > 0 {
			avatarData = pngData
		}
	case ".json":
		// 直接解析JSON格式
		card = &sillytavern.CharacterCard{}
		if err = json.Unmarshal(fileContent, card); err != nil {
			return nil, errs.NewStandardf(connect.CodeInvalidArgument, "无法解析JSON文件: %v", err)
		}
	default:
		return nil, errs.NewStandardf(connect.CodeInvalidArgument, "不支持的文件格式: %s", ext)
	}

	if card == nil {
		return nil, errs.NewStandardf(connect.CodeInvalidArgument, "图片不包含任何相关数据: %v", err)
	}

	// 解析头像数据（如果不是PNG文件或没有从PNG获取头像数据）
	if len(avatarData) == 0 && card.Avatar != "" && card.Avatar != "none" {
		// 尝试从Avatar字段解析base64数据
		avatarStr := card.Avatar
		// 处理data URI格式
		if strings.HasPrefix(avatarStr, "data:image/") {
			// 提取base64部分
			parts := strings.SplitN(avatarStr, ",", 2)
			if len(parts) == 2 {
				avatarStr = parts[1]
			}
		}
		// 解码base64
		decoded, decodeErr := base64.StdEncoding.DecodeString(avatarStr)
		if decodeErr == nil && len(decoded) > 0 {
			avatarData = decoded
		}
	}

	userId := jwt.GetUserId(ctx)
	character := convert.STCharacterCardToEntity(card, userId)
	// 清空Avatar字段，后面会设置为文件路径
	character.Avatar = ""

	// 压缩世界书数据
	var buf bytes.Buffer
	zip := gzip.NewWriter(&buf)
	defer zip.Close()
	worldBookJson, _ := json.Marshal(character.WorldInfo)
	if _, err = zip.Write(worldBookJson); err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "压缩世界书失败： %v", err)
	}
	character.WorldInfoBackup = buf.Bytes()

	if err = c.charRepo.Create(ctx, character); err != nil {
		return nil, err
	}
	log.Infof("角色创建成功，ID: %d, 名称: %s", character.ID, character.Name)
	log.Debugf("角色世界书创建成功: %d", len(character.WorldInfo.Entries))

	// 保存头像到文件系统
	if len(avatarData) > 0 {
		avatarPath, saveErr := fileutil.SaveAvatarFile(ctx, userId, character.ID, character.Name, avatarData, pb.FileType_CharAvatar)
		if saveErr != nil {
			return nil, errs.NewStandardf(connect.CodeInternal, "保存头像文件失败: %v", saveErr)
		}
		// 更新角色的Avatar字段
		character.Avatar = avatarPath
		if updateErr := c.charRepo.Update(ctx, character); updateErr != nil {
			if err = fileutil.DeleteAvatarFile(fileutil.BuildAvatarURL(config.Get().File.UploadPath,
				userId, character.ID, character.Name, pb.FileType_CharAvatar)); err != nil {
				log.Warnf("删除头像文件失败: %v", err)
			}
			return nil, errs.NewStandardf(connect.CodeInternal, "更新角色头像路径失败: %v", err)
		}
	}

	// 构建响应前先打印调试信息
	pbChar := convert.CharaEntityToPb(character)
	resp := &pb.ImportCharacterRsp{
		Character: pbChar,
	}
	return resp, nil
}

func (c *characterImpl) ExportCharacter(ctx context.Context, req *pb.ExportCharacterReq) (*pb.ExportCharacterRsp, error) {
	userId := jwt.GetUserId(ctx)
	character, err := c.charRepo.GetByID(ctx, int(req.GetId()), userId)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.CharacterNotFound)
	}

	// 将角色数据转换为CharacterCard
	card := convert.EntityToSTCharacterCard(character)

	// 从文件系统读取头像PNG数据
	var avatarPNG []byte
	if character.Avatar != "" {
		avatarPNG, err = file.ReadAvatarFile(ctx, userId, character.Avatar)
		if err != nil {
			log.Warnf("读取头像文件失败: %v, 使用默认头像", err)
		}
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
	return &pb.ExportCharacterRsp{
		FileContent: pngData,
		FileName:    fileName,
	}, nil
}

func (c *characterImpl) RestoreCharacterWorldInfo(ctx context.Context, req *pb.RestoreCharacterWorldInfoReq) (*pb.RestoreCharacterWorldInfoRsp, error) {
	//TODO implement me
	panic("implement me")
}

// readFromPNG 从PNG文件字节数据中读取角色卡信息
// 返回角色卡数据、干净的PNG图片数据（去掉嵌入的角色卡元数据）和错误
func readFromPNG(data []byte) (*sillytavern.CharacterCard, []byte, error) {
	// 使用库解析PNG
	pmp := png.NewPngMediaParser()
	mc, err := pmp.ParseBytes(data)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid PNG files: signature mismatch: %v", err)
	}

	// 类型断言为ChunkSlice
	cs, ok := mc.(*png.ChunkSlice)
	if !ok {
		return nil, nil, fmt.Errorf("invalid PNG files: not a chunk slice")
	}

	// 获取所有Chunk
	chunks := cs.Chunks()

	// 在tEXt块中查找角色卡数据
	var charaData, ccv3Data string
	for _, chunk := range chunks {
		log.Debugf("chunk块的类型: %s", chunk.Type)
		if chunk.Type != "tEXt" {
			continue
		}

		keyword, text := parseTextChunk(chunk.Data)
		log.Debugf("chunk的关键字: %s", keyword)
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
		return nil, nil, fmt.Errorf("no character card data found")
	}

	// Base64解码
	decoded, err := base64.StdEncoding.DecodeString(textData)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid base64 data: %v", err)
	}

	// JSON解析
	var card sillytavern.CharacterCard
	if err = json.Unmarshal(decoded, &card); err != nil {
		return nil, nil, fmt.Errorf("invalid JSON data: %v", err)
	}

	// 过滤掉包含角色卡数据的tEXt块，构建干净的PNG图片数据
	var filteredChunks []*png.Chunk
	for _, chunk := range chunks {
		// 只过滤掉chara和ccv3的tEXt块
		if chunk.Type == "tEXt" {
			keyword, _ := parseTextChunk(chunk.Data)
			keywordLower := strings.ToLower(keyword)
			if keywordLower == "chara" || keywordLower == "ccv3" {
				continue
			}
		}
		filteredChunks = append(filteredChunks, chunk)
	}

	// 重新构建干净的PNG
	cleanPNG, buildErr := buildPNGFromChunks(filteredChunks)
	if buildErr != nil {
		return nil, nil, fmt.Errorf("failed to build clean PNG: %v", buildErr)
	}

	return &card, cleanPNG, nil
}

// buildPNGFromChunks 从chunks构建PNG数据
func buildPNGFromChunks(chunks []*png.Chunk) ([]byte, error) {
	newCs := png.NewChunkSlice(chunks)
	var buf bytes.Buffer
	if err := newCs.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteToPNG 将角色卡数据写入PNG图片
func WriteToPNG(originalPNG []byte, card *sillytavern.CharacterCard) ([]byte, error) {
	// 使用库解析PNG
	pmp := png.NewPngMediaParser()
	mc, err := pmp.ParseBytes(originalPNG)
	if err != nil {
		return nil, fmt.Errorf("invalid PNG files: signature mismatch: %v", err)
	}

	// 类型断言为ChunkSlice
	cs, ok := mc.(*png.ChunkSlice)
	if !ok {
		return nil, fmt.Errorf("invalid PNG files: not a chunk slice")
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
	card, _, err := readFromPNG(data)
	return err == nil && card != nil && card.Name != ""
}
