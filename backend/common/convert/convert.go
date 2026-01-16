package convert

import (
	"strings"

	"github.com/ling/muse/entity"
	"github.com/ling/muse/entity/sillytavern"
	pb "github.com/ling/muse/gen/muse"
)

// CharaEntityToPb 将角色卡实体类转换为pb
func CharaEntityToPb(character *entity.Character) *pb.Character {
	worldInfoEntries := make([]*pb.WorldInfoEntry, len(character.WorldInfo.Entries))
	for _, entry := range character.WorldInfo.Entries {
		worldInfoEntries = append(worldInfoEntries, WorldInfoEntryEntityToPb(&entry))
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

// RegexRuleEntityToPb 将正则规则实体类转换为pb
func RegexRuleEntityToPb(rule *entity.RegexRule) *pb.RegexRule {
	return &pb.RegexRule{
		Id:              int32(rule.ID),
		PresetId:        int32(rule.PresetID),
		Name:            rule.Name,
		FindPattern:     rule.FindPattern,
		ReplacePattern:  rule.ReplacePattern,
		IsEnabled:       rule.IsEnabled,
		RunOnEdit:       rule.RunOnEdit,
		SubstituteRegex: rule.SubstituteRegex,
		MinDepth:        int32(rule.MinDepth),
		MaxDepth:        int32(rule.MaxDepth),
		AffectFlags: &pb.RegexAffectFlags{
			UserInput:    rule.AffectFlagsUserInput,
			AiOutput:     rule.AffectFlagsAIOutput,
			SlashCommand: rule.AffectFlagsSlashCommand,
			WorldInfo:    rule.AffectFlagsWorldInfo,
			Prompt:       rule.AffectFlagsPrompt,
		},
		SortOrder: int32(rule.SortOrder),
		CreatedAt: rule.CreatedAt.Unix(),
		UpdatedAt: rule.UpdatedAt.Unix(),
	}
}

// WorldInfoEntityToPb 将世界书实体类转换为pb
func WorldInfoEntityToPb(worldInfo *entity.WorldInfo) *pb.WorldInfo {
	entries := make([]*pb.WorldInfoEntry, len(worldInfo.Entries))
	for i, entry := range worldInfo.Entries {
		entries[i] = WorldInfoEntryEntityToPb(&entry)
	}
	return &pb.WorldInfo{
		Id:          int32(worldInfo.ID),
		UserId:      int32(worldInfo.UserID),
		Name:        worldInfo.Name,
		Description: worldInfo.Description,
		IsGlobal:    worldInfo.IsGlobal,
		CreatedAt:   worldInfo.CreatedAt.Unix(),
		UpdatedAt:   worldInfo.UpdatedAt.Unix(),
		Entries:     entries,
	}
}

// WorldInfoEntryEntityToPb 将世界书条目实体类转换为pb
func WorldInfoEntryEntityToPb(entry *entity.WorldInfoEntry) *pb.WorldInfoEntry {
	return &pb.WorldInfoEntry{
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
	}
}

// PresetEntityToPb 将预设实体类转换为pb
func PresetEntityToPb(preset *entity.Preset) *pb.Preset {
	promptItems := make([]*pb.PromptItem, len(preset.PromptItems))
	for i, item := range preset.PromptItems {
		promptItems[i] = PromptItemEntityToPb(&item)
	}
	return &pb.Preset{
		Id:               int32(preset.ID),
		UserId:           int32(preset.UserID),
		Name:             preset.Name,
		Temperature:      preset.Temperature,
		TopP:             preset.TopP,
		TopK:             int32(preset.TopK),
		MaxTokens:        int32(preset.MaxTokens),
		FrequencyPenalty: preset.FrequencyPenalty,
		PresencePenalty:  preset.PresencePenalty,
		Version:          int64(preset.Version),
		CreatedAt:        preset.CreatedAt.Unix(),
		UpdatedAt:        preset.UpdatedAt.Unix(),
		PromptItems:      promptItems,
	}
}

// PromptItemEntityToPb 将提示项实体类转换为pb
func PromptItemEntityToPb(item *entity.PromptItem) *pb.PromptItem {
	return &pb.PromptItem{
		Id:                int32(item.ID),
		PresetId:          int32(item.PresetID),
		Identifier:        item.Identifier,
		Name:              item.Name,
		Content:           item.Content,
		Role:              item.Role,
		IsEnabled:         item.IsEnabled,
		InjectionPosition: item.InjectionPosition,
		InjectionDepth:    int32(item.InjectionDepth),
		ForbidOverrides:   item.ForbidOverrides,
		SortOrder:         int32(item.SortOrder),
		CreatedAt:         item.CreatedAt.Unix(),
		UpdatedAt:         item.UpdatedAt.Unix(),
	}
}

// SessionEntityToPb 将聊天会话实体类转换为pb
func SessionEntityToPb(session *entity.ChatSession) *pb.ChatSession {
	messages := make([]*pb.Message, len(session.Messages))
	for i, msg := range session.Messages {
		messages[i] = MessageEntityToPb(&msg)
	}

	var character *pb.Character
	if session.Character != nil {
		character = CharaEntityToPb(session.Character)
	}

	return &pb.ChatSession{
		Id:          int32(session.ID),
		UserId:      int32(session.UserID),
		CharacterId: int32(session.CharacterID),
		Name:        session.Name,
		Version:     int64(session.Version),
		CreatedAt:   session.CreatedAt.Unix(),
		UpdatedAt:   session.UpdatedAt.Unix(),
		Character:   character,
		Messages:    messages,
	}
}

// MessageEntityToPb 将消息实体类转换为pb
func MessageEntityToPb(message *entity.Message) *pb.Message {
	swipes := make([]*pb.MessageSwipe, len(message.Swipes))
	for i, swipe := range message.Swipes {
		swipes[i] = MessageSwipeEntityToPb(&swipe)
	}
	return &pb.Message{
		Id:               int32(message.ID),
		SessionId:        int32(message.SessionID),
		Role:             message.Role,
		ActiveSwipeIndex: int32(message.ActiveSwipeIndex),
		SortOrder:        int32(message.SortOrder),
		CreatedAt:        message.CreatedAt.Unix(),
		UpdatedAt:        message.UpdatedAt.Unix(),
		Swipes:           swipes,
	}
}

// MessageSwipeEntityToPb 将消息Swipe实体类转换为pb
func MessageSwipeEntityToPb(swipe *entity.MessageSwipe) *pb.MessageSwipe {
	return &pb.MessageSwipe{
		Id:        int32(swipe.ID),
		MessageId: int32(swipe.MessageID),
		Content:   swipe.Content,
		SortOrder: int32(swipe.SortOrder),
		CreatedAt: swipe.CreatedAt.Unix(),
	}
}

// ============ SillyTavern 预设转换 ============

// STPresetToEntity 将 SillyTavern OpenAI 预设转换为 Muse 预设实体
// userID: 用户ID
// presetName: 预设名称
func STPresetToEntity(stPreset *sillytavern.OpenAIPreset, userID int, presetName string) *entity.Preset {
	preset := &entity.Preset{
		UserID:           userID,
		Name:             presetName,
		Temperature:      stPreset.Temperature,
		TopP:             stPreset.TopP,
		TopK:             stPreset.TopK,
		MaxTokens:        stPreset.OpenAIMaxTokens,
		FrequencyPenalty: stPreset.FrequencyPenalty,
		PresencePenalty:  stPreset.PresencePenalty,
		Version:          1,
	}

	// 设置默认值（如果未设置）
	if preset.Temperature == 0 {
		preset.Temperature = 1.0
	}
	if preset.TopP == 0 {
		preset.TopP = 1.0
	}
	if preset.MaxTokens == 0 {
		preset.MaxTokens = 300
	}

	return preset
}

// STPromptToEntity 将 SillyTavern 提示项转换为 Muse 提示项实体
// presetID: 关联的预设ID
// stPrompt: SillyTavern 提示项
// sortOrder: 排序顺序
// promptOrderMap: 提示词顺序映射（用于确定启用状态）
func STPromptToEntity(presetID int, stPrompt *sillytavern.PresetPromptItem, sortOrder int, promptOrderMap map[string]sillytavern.PromptOrderIdentifier) *entity.PromptItem {
	// 转换角色
	role := ConvertSTRole(stPrompt.Role)

	// 从 prompt_order 获取启用状态，默认启用
	isEnabled := true
	if orderInfo, ok := promptOrderMap[stPrompt.Identifier]; ok {
		isEnabled = orderInfo.Enabled
	}

	// 转换注入位置
	injectionPosition := pb.InjectionPosition_Relative
	if stPrompt.InjectionPosition == 1 {
		injectionPosition = pb.InjectionPosition_Absolute
	}

	return &entity.PromptItem{
		PresetID:          presetID,
		Identifier:        stPrompt.Identifier,
		Name:              stPrompt.Name,
		Content:           stPrompt.Content,
		Role:              role,
		IsEnabled:         isEnabled,
		InjectionPosition: injectionPosition,
		InjectionDepth:    stPrompt.InjectionDepth,
		ForbidOverrides:   stPrompt.ForbidOverrides,
		SortOrder:         sortOrder,
	}
}

// BuildPromptOrderMap 构建提示词顺序映射
// SillyTavern 的 prompt_order 包含多个角色的顺序配置，这里取默认角色（100000）的配置
func BuildPromptOrderMap(promptOrder []sillytavern.PromptOrderItem) map[string]sillytavern.PromptOrderIdentifier {
	result := make(map[string]sillytavern.PromptOrderIdentifier)

	for _, order := range promptOrder {
		// 优先使用默认角色（100000）的配置
		if order.CharacterID == 100000 {
			for _, item := range order.Order {
				result[item.Identifier] = item
			}
			break
		}
	}

	// 如果没有找到默认角色，使用第一个配置
	if len(result) == 0 && len(promptOrder) > 0 {
		for _, item := range promptOrder[0].Order {
			result[item.Identifier] = item
		}
	}

	return result
}

// ConvertSTRole 将 SillyTavern 角色字符串转换为 pb.Role
func ConvertSTRole(role string) pb.Role {
	switch strings.ToLower(role) {
	case "system":
		return pb.Role_System
	case "user":
		return pb.Role_User
	case "assistant":
		return pb.Role_Assistant
	default:
		return pb.Role_System // 默认为系统角色
	}
}
