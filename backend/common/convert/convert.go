package convert

import (
	"fmt"
	"strings"

	"github.com/ling/muse/config"
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
		CharacterId:     int32(rule.CharacterID),
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
		CandidateCount:   int32(preset.CandidateCount),
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

// ==================== 用户相关转换函数 ====================

// UserEntityToPb 将用户实体类转换为pb
func UserEntityToPb(user *entity.User) *pb.SysUser {
	return &pb.SysUser{
		Id:              int32(user.ID),
		Username:        user.Username,
		ActivePersonaId: int32(user.ActivePersonaID),
		ActivePresetId:  int32(user.ActivePresetID),
		CreatedAt:       user.CreatedAt.Unix(),
		UpdatedAt:       user.UpdatedAt.Unix(),
	}
}

// PersonaEntityToPb 将人设实体类转换为pb
func PersonaEntityToPb(persona *entity.Persona) *pb.Persona {
	return &pb.Persona{
		Id:          int32(persona.ID),
		UserId:      int32(persona.UserID),
		Name:        persona.Name,
		Avatar:      persona.Avatar,
		Description: persona.Description,
		CreatedAt:   persona.CreatedAt.Unix(),
		UpdatedAt:   persona.UpdatedAt.Unix(),
	}
}

// UserSettingEntityToPb 将用户设置实体类转换为pb
func UserSettingEntityToPb(setting *entity.UserSetting) *pb.UserSetting {
	return &pb.UserSetting{
		Id:             int32(setting.ID),
		UserId:         int32(setting.UserID),
		Theme:          pb.Theme(setting.Theme),
		Language:       setting.Language,
		SendOnEnter:    setting.SendOnEnter,
		ShowTimestamps: setting.ShowTimestamps,
		CreatedAt:      setting.CreatedAt.Unix(),
		UpdatedAt:      setting.UpdatedAt.Unix(),
	}
}

// APIConfigEntityToPb 将API配置实体类转换为pb（不包含API Key）
func APIConfigEntityToPb(config *entity.APIConfig) *pb.APIConfig {
	return &pb.APIConfig{
		Id:        int32(config.ID),
		UserId:    int32(config.UserID),
		Name:      config.Name,
		Provider:  config.Provider,
		ApiKey:    config.APIKey,
		BaseUrl:   config.BaseURL,
		Model:     config.Model,
		IsActive:  config.IsActive,
		CreatedAt: config.CreatedAt.Unix(),
		UpdatedAt: config.UpdatedAt.Unix(),
	}
}

// APIConfigEntityToPbWithKey 将API配置实体类转换为pb（包含API Key）
// 根据配置决定是否返回API Key，以及是否返回加密后的密文
func APIConfigEntityToPbWithKey(apiConfig *entity.APIConfig) *pb.APIConfig {
	apiKey := ""

	// 检查配置是否允许获取API Key
	apiEncryptCfg := config.Get().APIEncrypt
	if apiEncryptCfg.AllowGetKey {
		// 如果启用了加密，直接返回数据库中的密文（前端需要解密）
		// 如果没有启用加密，返回明文
		apiKey = apiConfig.APIKey
	}

	return &pb.APIConfig{
		Id:        int32(apiConfig.ID),
		UserId:    int32(apiConfig.UserID),
		Name:      apiConfig.Name,
		Provider:  apiConfig.Provider,
		ApiKey:    apiKey,
		BaseUrl:   apiConfig.BaseURL,
		Model:     apiConfig.Model,
		IsActive:  apiConfig.IsActive,
		CreatedAt: apiConfig.CreatedAt.Unix(),
		UpdatedAt: apiConfig.UpdatedAt.Unix(),
	}
}

// ============ SillyTavern 正则规则转换 ============

// STRegexToEntity 将 SillyTavern 正则脚本转换为 Muse 实体
// presetID: 0 表示全局正则，非 0 表示预设正则（Preset Scripts）
// characterID: 0 表示非角色范围正则，非 0 表示角色正则（Scoped Scripts）
// sortOrder: 排序顺序
func STRegexToEntity(stScript *sillytavern.RegexScript, presetID int, characterID int, sortOrder int) *entity.RegexRule {
	rule := &entity.RegexRule{
		PresetID:        presetID,
		CharacterID:     characterID,
		Name:            stScript.ScriptName,
		FindPattern:     stScript.FindRegex,
		ReplacePattern:  stScript.ReplaceString,
		IsEnabled:       !stScript.Disabled,
		RunOnEdit:       stScript.RunOnEdit,
		SubstituteRegex: stScript.SubstituteRegex > 0,
		MinDepth:        stScript.MinDepth,
		MaxDepth:        stScript.MaxDepth,
		SortOrder:       sortOrder,
	}

	// 解析 placement 数组，映射到 affect flags
	// SillyTavern placement: 0=用户输入, 1=AI输出, 2=斜杠命令, 3=世界书, 4=提示词
	for _, p := range stScript.Placement {
		switch p {
		case 0:
			rule.AffectFlagsUserInput = true
		case 1:
			rule.AffectFlagsAIOutput = true
		case 2:
			rule.AffectFlagsSlashCommand = true
		case 3:
			rule.AffectFlagsWorldInfo = true
		case 4:
			rule.AffectFlagsPrompt = true
		}
	}

	return rule
}

// EntityToSTRegex 将 Muse 实体转换为 SillyTavern 正则脚本格式
func EntityToSTRegex(rule *entity.RegexRule) sillytavern.RegexScript {
	stScript := sillytavern.RegexScript{
		ID:            fmt.Sprintf("%d", rule.ID), // 使用 ID 作为标识符
		ScriptName:    rule.Name,
		FindRegex:     rule.FindPattern,
		ReplaceString: rule.ReplacePattern,
		Disabled:      !rule.IsEnabled,
		MarkdownOnly:  false, // Muse 暂不支持此字段
		PromptOnly:    false, // Muse 暂不支持此字段
		RunOnEdit:     rule.RunOnEdit,
		MinDepth:      rule.MinDepth,
		MaxDepth:      rule.MaxDepth,
	}

	// 转换 SubstituteRegex
	if rule.SubstituteRegex {
		stScript.SubstituteRegex = 3 // 全部替换
	}

	// 构建 placement 数组
	var placements []int
	if rule.AffectFlagsUserInput {
		placements = append(placements, 0)
	}
	if rule.AffectFlagsAIOutput {
		placements = append(placements, 1)
	}
	if rule.AffectFlagsSlashCommand {
		placements = append(placements, 2)
	}
	if rule.AffectFlagsWorldInfo {
		placements = append(placements, 3)
	}
	if rule.AffectFlagsPrompt {
		placements = append(placements, 4)
	}
	stScript.Placement = placements

	return stScript
}

// ============ SillyTavern 世界书转换 ============

// STWorldInfoToEntity 将 SillyTavern 世界书转换为 Muse 世界书实体
// userID: 用户ID
// worldInfoName: 世界书名称
func STWorldInfoToEntity(stWorldBook *sillytavern.WorldBook, userID int, worldInfoName string) *entity.WorldInfo {
	worldInfo := &entity.WorldInfo{
		UserID:      userID,
		Name:        worldInfoName,
		Description: stWorldBook.Description,
		IsGlobal:    false, // 默认非全局
	}

	// 如果 WorldBook 有名称，使用它
	if stWorldBook.Name != "" {
		worldInfo.Name = stWorldBook.Name
	}

	return worldInfo
}

// STBookEntryToEntity 将 SillyTavern 世界书条目转换为 Muse 条目实体
// worldInfoID: 关联的世界书ID
// stEntry: SillyTavern 条目
// sortOrder: 排序顺序
func STBookEntryToEntity(worldInfoID int, stEntry *sillytavern.BookEntry, sortOrder int) *entity.WorldInfoEntry {
	// 将关键词数组转换为逗号分隔的字符串
	keysList := strings.Join(stEntry.Key, ",")
	secondaryKeys := strings.Join(stEntry.KeySecondary, ",")

	// 转换位置
	// SillyTavern: 0=before_char, 1=after_char, 2=before_desc, 3=after_desc, 4=at_depth
	// Muse pb.EntryPosition: 0=BeforeCharacter, 1=AfterCharacter, 2=BeforeAuthorNote, 3=AfterAuthorNote, 4=AtDepth
	position := pb.EntryPosition(stEntry.Position)

	return &entity.WorldInfoEntry{
		WorldInfoID:    worldInfoID,
		UID:            fmt.Sprintf("%d", stEntry.UID),
		KeysList:       keysList,
		SecondaryKeys:  secondaryKeys,
		Content:        stEntry.Content,
		Comment:        stEntry.Comment,
		IsEnabled:      !stEntry.Disable,
		Constant:       stEntry.Constant,
		Selective:      stEntry.Selective,
		InsertionOrder: stEntry.Order,
		Position:       position,
		Depth:          stEntry.Depth,
		SortOrder:      sortOrder,
	}
}

// EntityToSTWorldInfo 将 Muse 世界书实体转换为 SillyTavern 世界书格式
func EntityToSTWorldInfo(worldInfo *entity.WorldInfo) *sillytavern.WorldBook {
	stWorldBook := &sillytavern.WorldBook{
		Name:        worldInfo.Name,
		Description: worldInfo.Description,
		Entries:     make(map[string]sillytavern.BookEntry),
	}

	// 转换每个条目
	for i, entry := range worldInfo.Entries {
		stEntry := EntityToSTBookEntry(&entry, i)
		// 使用 UID 或 索引作为 key
		key := entry.UID
		if key == "" {
			key = fmt.Sprintf("%d", i)
		}
		stWorldBook.Entries[key] = stEntry
	}

	return stWorldBook
}

// EntityToSTBookEntry 将 Muse 条目实体转换为 SillyTavern 条目格式
func EntityToSTBookEntry(entry *entity.WorldInfoEntry, index int) sillytavern.BookEntry {
	// 将逗号分隔的字符串转换为数组
	var keys []string
	if entry.KeysList != "" {
		keys = strings.Split(entry.KeysList, ",")
		// 去除空白
		for i := range keys {
			keys[i] = strings.TrimSpace(keys[i])
		}
	}

	var secondaryKeys []string
	if entry.SecondaryKeys != "" {
		secondaryKeys = strings.Split(entry.SecondaryKeys, ",")
		for i := range secondaryKeys {
			secondaryKeys[i] = strings.TrimSpace(secondaryKeys[i])
		}
	}

	// 尝试解析 UID 为整数，如果失败则使用索引
	uid := index
	if entry.UID != "" {
		if parsedUID, err := fmt.Sscanf(entry.UID, "%d", &uid); err != nil || parsedUID != 1 {
			uid = index
		}
	}

	return sillytavern.BookEntry{
		UID:            uid,
		Key:            keys,
		KeySecondary:   secondaryKeys,
		Content:        entry.Content,
		Comment:        entry.Comment,
		Disable:        !entry.IsEnabled,
		Constant:       entry.Constant,
		Selective:      entry.Selective,
		Order:          entry.InsertionOrder,
		Position:       int(entry.Position),
		Depth:          entry.Depth,
		DisplayIndex:   entry.SortOrder,
		Probability:    100,  // 默认值
		UseProbability: true, // 默认值
		GroupWeight:    100,  // 默认值
	}
}
