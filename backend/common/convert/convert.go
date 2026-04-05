package convert

import (
	"fmt"
	"strings"

	"github.com/ling/muse/common/constant"
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
		KeysList:       entry.Keys,
		SecondaryKeys:  entry.SecondaryKeys,
		Content:        entry.Content,
		Comment:        entry.Comment,
		IsEnabled:      entry.IsEnabled,
		Constant:       entry.Constant,
		Selective:      entry.Selective,
		InsertionOrder: int32(entry.InsertionOrder),
		Position:       entry.Position,
		Depth:          int32(entry.Depth),
		Role:           entry.Role,
		SortOrder:      int32(entry.SortOrder),
		CreatedAt:      entry.CreatedAt.Unix(),
		UpdatedAt:      entry.UpdatedAt.Unix(),
	}
}

// PresetEntityToPb 将预设实体类转换为pb
func PresetEntityToPb(preset *entity.Preset) (*pb.Preset, []*pb.PromptItem, []*pb.RegexRule) {
	promptItems := make([]*pb.PromptItem, len(preset.PromptItems))
	for i, item := range preset.PromptItems {
		promptItems[i] = PromptItemEntityToPb(&item)
	}

	regexs := make([]*pb.RegexRule, len(preset.RegexRules))
	for i, rule := range preset.RegexRules {
		regexs[i] = RegexRuleEntityToPb(&rule)
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
	}, promptItems, regexs
}

// PromptItemEntityToPb 将提示项实体类转换为pb
func PromptItemEntityToPb(item *entity.PromptItem) *pb.PromptItem {
	var pre *int32 = nil
	if item.Pre != nil {
		val := int32(*item.Pre)
		pre = &val
	}
	var next *int32 = nil
	if item.Next != nil {
		val := int32(*item.Next)
		next = &val
	}
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
		CreatedAt:         item.CreatedAt.Unix(),
		UpdatedAt:         item.UpdatedAt.Unix(),
		Pre:               pre,
		Next:              next,
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
		CreatedAt: swipe.CreatedAt.Unix(),
	}
}

// PromptItemPbToEntity 将提示项pb转换为实体类
func PromptItemPbToEntity(item *pb.PromptItem) *entity.PromptItem {
	var pre *int = nil
	if item.Pre != nil {
		val := int(*item.Pre)
		pre = &val
	}

	var next *int = nil
	if item.Next != nil {
		val := int(*item.Next)
		next = &val
	}

	return &entity.PromptItem{
		ID:                int(item.Id),
		PresetID:          int(item.PresetId),
		Identifier:        item.Identifier,
		Name:              item.Name,
		Content:           item.Content,
		Role:              item.Role,
		IsEnabled:         item.IsEnabled,
		InjectionPosition: item.InjectionPosition,
		InjectionDepth:    int(item.InjectionDepth),
		ForbidOverrides:   item.ForbidOverrides,
		Pre:               pre,
		Next:              next,
	}
}

// ============ SillyTavern 预设转换 ============

// STPresetToEntity 将 SillyTavern OpenAI 预设转换为 Muse 预设实体
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

// ConvertSTIdentifier 将 SillyTavern 标识符字符串转换为 PromptItemIdentifier 枚举
func ConvertSTIdentifier(identifier string) pb.PromptItemIdentifier {
	switch identifier {
	case "main":
		return pb.PromptItemIdentifier_Main
	case "worldInfoBefore":
		return pb.PromptItemIdentifier_WorldInfoBefore
	case "worldInfoAfter":
		return pb.PromptItemIdentifier_WorldInfoAfter
	case "personaDescription":
		return pb.PromptItemIdentifier_PersonaDescription
	case "charDescription":
		return pb.PromptItemIdentifier_CharDescription
	case "charPersonality":
		return pb.PromptItemIdentifier_CharPersonality
	case "scenario":
		return pb.PromptItemIdentifier_Scenario
	case "nsfw":
		return pb.PromptItemIdentifier_Nsfw
	case "jailbreak":
		return pb.PromptItemIdentifier_Jailbreak
	case "dialogueExamples":
		return pb.PromptItemIdentifier_DialogueExamples
	case "chatHistory":
		return pb.PromptItemIdentifier_ChatHistory
	default:
		// 默认返回未指定
		return pb.PromptItemIdentifier_PromptItemIdentifierUnspecified
	}
}

// STPromptToEntity 将 SillyTavern 提示项转换为 Muse 提示项实体
func STPromptToEntity(presetID int, stPrompt *sillytavern.PresetPromptItem) *entity.PromptItem {
	// 转换角色
	role := ConvertSTRole(stPrompt.Role)

	// 转换标识符
	identifier := ConvertSTIdentifier(stPrompt.Identifier)
	// 转换注入位置
	injectionPosition := pb.InjectionPosition_Relative
	if stPrompt.InjectionPosition == constant.STInjectPosAbsolute {
		injectionPosition = pb.InjectionPosition_Absolute
	}

	return &entity.PromptItem{
		PresetID:          presetID,
		Identifier:        identifier,
		Name:              stPrompt.Name,
		Content:           stPrompt.Content,
		Role:              role,
		IsEnabled:         stPrompt.Enabled, // 注意这里是从排序时得来的值
		InjectionPosition: injectionPosition,
		InjectionDepth:    stPrompt.InjectionDepth,
		ForbidOverrides:   stPrompt.ForbidOverrides,
	}
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
		Theme:           user.Theme,
		Language:        user.Language,
		SendOnEnter:     user.SendOnEnter,
		ShowTimestamps:  user.ShowTimestamps,
		Provider:        user.Provider,
		Model:           user.Model,
		BaseUrl:         user.BaseURL,
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

// APIConfigEntityToPb 将API配置实体类转换为pb（不包含API Key）
func APIConfigEntityToPb(config *entity.APIConfig) *pb.APIConfig {
	return &pb.APIConfig{
		Id:        int32(config.ID),
		UserId:    int32(config.UserID),
		Provider:  config.Provider,
		ApiKey:    config.APIKey,
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
		Provider:  apiConfig.Provider,
		ApiKey:    apiKey,
		Model:     apiConfig.Model,
		IsActive:  apiConfig.IsActive,
		CreatedAt: apiConfig.CreatedAt.Unix(),
		UpdatedAt: apiConfig.UpdatedAt.Unix(),
	}
}

// ============ SillyTavern 正则规则转换 ============

// STRegexToEntity 将 SillyTavern 正则脚本转换为 Muse 实体
func STRegexToEntity(stScript *sillytavern.RegexScript, userID int) *entity.RegexRule {
	rule := &entity.RegexRule{
		UserID:          userID,
		Name:            stScript.ScriptName,
		FindPattern:     stScript.FindRegex,
		ReplacePattern:  stScript.ReplaceString,
		IsEnabled:       !stScript.Disabled,
		RunOnEdit:       stScript.RunOnEdit,
		SubstituteRegex: stScript.SubstituteRegex > 0,
		MinDepth:        stScript.MinDepth,
		MaxDepth:        stScript.MaxDepth,
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
func EntityToSTRegex(rule *entity.RegexRule) *sillytavern.RegexScript {
	stScript := &sillytavern.RegexScript{
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
func STWorldInfoToEntity(stWorldBook *sillytavern.WorldBook) *entity.WorldInfo {
	worldInfo := &entity.WorldInfo{
		Name:        stWorldBook.Name,
		Description: stWorldBook.Description,
		IsGlobal:    false, // 默认非全局
		Entries:     make([]entity.WorldInfoEntry, len(stWorldBook.Entries)),
	}
	// TODO 转换Entries
	return worldInfo
}

// STBookEntryToEntity 将 SillyTavern 世界书条目转换为 Muse 条目实体
func STBookEntryToEntity(stEntry *sillytavern.BookEntry) *entity.WorldInfoEntry {
	// 转换位置
	position := pb.EntryPosition(stEntry.Position + 1)
	return &entity.WorldInfoEntry{
		UID:            fmt.Sprintf("%d", stEntry.UID),
		Keys:           stEntry.Key,
		SecondaryKeys:  stEntry.KeySecondary,
		Content:        stEntry.Content,
		Comment:        stEntry.Comment,
		IsEnabled:      !stEntry.Disable,
		Constant:       stEntry.Constant,
		Selective:      stEntry.Selective,
		InsertionOrder: stEntry.Order,
		Position:       position,
		Depth:          stEntry.Depth,
		Role:           pb.Role_System, // 从SillyTavern导入时默认为系统角色
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
	// 尝试解析 UID 为整数，如果失败则使用索引
	uid := index
	if entry.UID != "" {
		if _, err := fmt.Sscanf(entry.UID, "%d", &uid); err != nil {
			uid = index
		}
	}

	// 转换位置
	// Muse pb.EntryPosition: 0=EntryPositionUnspecified, 1=BeforeChar, 2=AfterChar, 3=BeforeExample, 4=AfterExample, 5=AtDepth
	// SillyTavern BookEntry: 0=before_char, 1=after_char, 2=before_desc, 3=after_desc, 4=at_depth
	// 需要减1来匹配SillyTavern的枚举值
	position := int(entry.Position) - 1
	if position < 0 {
		position = 0
	}

	return sillytavern.BookEntry{
		UID:            uid,
		Key:            entry.Keys,
		KeySecondary:   entry.SecondaryKeys,
		Content:        entry.Content,
		Comment:        entry.Comment,
		Disable:        !entry.IsEnabled,
		Constant:       entry.Constant,
		Selective:      entry.Selective,
		Order:          entry.InsertionOrder,
		Position:       position,
		Depth:          entry.Depth,
		DisplayIndex:   entry.SortOrder,
		Probability:    100,  // 默认值
		UseProbability: true, // 默认值
		GroupWeight:    100,  // 默认值
	}
}

// STCharacterBookToEntity 将 SillyTavern 角色卡世界书转换为 Muse 世界书实体
func STCharacterBookToEntity(stCharacterBook *sillytavern.CharacterBook) *entity.WorldInfo {
	worldInfo := &entity.WorldInfo{
		Name:        stCharacterBook.Name,
		Description: "",
		IsGlobal:    false, // 角色卡世界书默认非全局
		Entries:     make([]entity.WorldInfoEntry, len(stCharacterBook.Entries)),
	}
	for i, stEntry := range stCharacterBook.Entries {
		worldInfo.Entries[i] = *STCharacterBookEntryToEntity(&stEntry)
		worldInfo.Entries[i].SortOrder = i
	}
	return worldInfo
}

// BookPositionToPb 将 SillyTavern 位置字符串转换为 pb.EntryPosition
// SillyTavern: "before_char", "after_char", "before_desc", "after_desc", "at_depth"
func BookPositionToPb(position int) pb.EntryPosition {
	switch position {
	case 0, 2:
		return pb.EntryPosition_BeforeChar
	case 1, 3:
		return pb.EntryPosition_AfterChar
	case 4:
		return pb.EntryPosition_AtDepth
	case 5:
		return pb.EntryPosition_BeforeExample
	case 6:
		return pb.EntryPosition_AfterExample
	default:
		return pb.EntryPosition_BeforeChar // 默认值
	}
}

// PbBookPositionToString 将 pb.EntryPosition 转换为 SillyTavern 位置字符串
func PbBookPositionToString(position pb.EntryPosition) int {
	switch position {
	case pb.EntryPosition_BeforeChar:
		return 0
	case pb.EntryPosition_AfterChar:
		return 1
	case pb.EntryPosition_BeforeExample:
		return 5
	case pb.EntryPosition_AfterExample:
		return 6
	case pb.EntryPosition_AtDepth:
		return 4
	default:
		return 0
	}
}

// STCharacterBookEntryToEntity 将 SillyTavern 角色卡世界书条目转换为 Muse 条目实体
func STCharacterBookEntryToEntity(stEntry *sillytavern.CharacterBookEntry) *entity.WorldInfoEntry {
	entry := &entity.WorldInfoEntry{
		UID:            fmt.Sprintf("%d", stEntry.ID),
		Keys:           stEntry.Keys,
		SecondaryKeys:  stEntry.SecondaryKeys,
		Content:        stEntry.Content,
		Comment:        stEntry.Comment,
		IsEnabled:      stEntry.Enabled,
		Constant:       stEntry.Constant,
		Selective:      stEntry.Selective,
		InsertionOrder: stEntry.InsertionOrder,
		Position:       BookPositionToPb(stEntry.Extensions.Position),
		Depth:          stEntry.Extensions.Depth,
		Role:           pb.Role_System, // 从SillyTavern导入时默认为系统角色
	}
	return entry
}

// EntityToSTCharacterBook 将 Muse 世界书实体转换为 SillyTavern 角色卡世界书格式
func EntityToSTCharacterBook(worldInfo *entity.WorldInfo) *sillytavern.CharacterBook {
	stCharacterBook := &sillytavern.CharacterBook{
		Name:    worldInfo.Name,
		Entries: make([]sillytavern.CharacterBookEntry, 0),
	}

	// 转换每个条目
	for _, entry := range worldInfo.Entries {
		stEntry := EntityToSTCharacterBookEntry(&entry)
		stCharacterBook.Entries = append(stCharacterBook.Entries, stEntry)
	}

	return stCharacterBook
}

// EntityToSTCharacterBookEntry 将 Muse 条目实体转换为 SillyTavern 角色卡世界书条目格式
func EntityToSTCharacterBookEntry(entry *entity.WorldInfoEntry) sillytavern.CharacterBookEntry {
	// 解析 UID 为整数 ID
	id := 0
	if entry.UID != "" {
		_, _ = fmt.Sscanf(entry.UID, "%d", &id)
	}

	return sillytavern.CharacterBookEntry{
		ID:             id,
		Keys:           entry.Keys,
		SecondaryKeys:  entry.SecondaryKeys,
		Content:        entry.Content,
		Comment:        entry.Comment,
		Enabled:        entry.IsEnabled,
		InsertionOrder: entry.InsertionOrder,
		Selective:      entry.Selective,
		Constant:       entry.Constant,
		UseRegex:       false, // 默认值
		Extensions: sillytavern.CharacterBookEntryExtensions{
			Depth:    entry.Depth,
			Position: PbBookPositionToString(entry.Position),
			// TODO 支持更多选项
		},
	}
}

// ============ SillyTavern 角色卡转换 ============

// STCharacterCardToEntity 将 SillyTavern 角色卡转换为 Muse 角色卡实体
// userID: 用户ID
func STCharacterCardToEntity(stCard *sillytavern.CharacterCard, userID int) *entity.Character {
	character := &entity.Character{}
	character.UserID = userID
	character.CreatorNotes = stCard.Data.CreatorNotes

	// 优先使用 data 中的字段，如果为空则使用顶层字段
	name := stCard.Data.Name
	if name == "" {
		name = stCard.Name
	}
	character.Name = name

	description := stCard.Data.Description
	if description == "" {
		description = stCard.Description
	}
	personality := stCard.Data.Personality
	if personality == "" {
		personality = stCard.Personality
	}
	character.Description = description + "\n" + personality

	exampleDialogue := stCard.Data.MesExample
	if exampleDialogue == "" {
		exampleDialogue = stCard.MesExample
	}
	// 将单个示例对话转换为数组
	var exampleDialogues []string
	if exampleDialogue != "" {
		exampleDialogues = []string{exampleDialogue}
	}
	character.ExampleDialogue = exampleDialogues

	character.Avatar = stCard.Avatar

	firstMessage := stCard.Data.FirstMes
	if firstMessage == "" {
		firstMessage = stCard.FirstMes
	}
	firstMessages := make([]string, 0, len(stCard.Data.AlternateGreetings)+1)
	if firstMessage != "" {
		firstMessages = append(firstMessages, stCard.Data.FirstMes)
	}
	firstMessages = append(firstMessages, stCard.Data.AlternateGreetings...)
	character.FirstMessage = firstMessages

	worldBook := STCharacterBookToEntity(&stCard.Data.CharacterBook)
	worldBook.UserID = userID
	character.WorldInfo = *worldBook

	regexs := make([]entity.RegexRule, len(stCard.Data.Extensions.RegexScripts))
	for i, regex := range stCard.Data.Extensions.RegexScripts {
		entityRegex := STRegexToEntity(&regex, userID)
		entityRegex.CharacterID = character.ID
		entityRegex.PresetID = 0
		entityRegex.SortOrder = i
		regexs[i] = *entityRegex
	}
	character.RegexRules = regexs

	return character
}

// EntityToSTCharacterCard 将 Muse 角色卡实体转换为 SillyTavern 角色卡格式
func EntityToSTCharacterCard(character *entity.Character) *sillytavern.CharacterCard {
	stCard := &sillytavern.CharacterCard{
		Spec:        "chara_card_v3",
		SpecVersion: "3.0",
		Name:        character.Name,
		Description: character.Description,
		Avatar:      character.Avatar,
		Data: sillytavern.CharacterData{
			Name:        character.Name,
			Description: character.Description,
			// 将数组转换为单个字符串（取第一个元素）
			MesExample: func() string {
				if len(character.ExampleDialogue) > 0 {
					return character.ExampleDialogue[0]
				}
				return ""
			}(),
			CreatorNotes: character.CreatorNotes,
		},
	}

	if len(character.FirstMessage) > 0 {
		stCard.Data.FirstMes = character.FirstMessage[0]
		stCard.Data.AlternateGreetings = character.FirstMessage[1:]
	}

	// 如果有关联的世界书，转换它
	if character.WorldInfo.ID > 0 {
		stCharacterBook := EntityToSTCharacterBook(&character.WorldInfo)
		stCard.Data.CharacterBook = *stCharacterBook
	}

	// 关联正则
	if len(character.RegexRules) > 0 {
		stCard.Data.Extensions.RegexScripts = make([]sillytavern.RegexScript, len(character.RegexRules))
		for i, regex := range character.RegexRules {
			stCard.Data.Extensions.RegexScripts[i] = *EntityToSTRegex(&regex)
		}
	}

	return stCard
}
