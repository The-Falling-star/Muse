package chat

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/constrant"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/common/jwt"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/cache"
	"github.com/ling/muse/repo/database"
	"github.com/ling/muse/repo/model"
)

type chatImpl struct {
	chatRepo      *database.ChatRepo
	presetRepo    *database.PresetRepo
	worldInfoRepo *database.WorldInfoRepo
	regexRepo     *database.RegexRuleRepo
	charRepo      *database.CharacterRepo
	userRepo      *database.UserRepo

	// 会话缓存管理器
	cacheManager *cache.SessionCacheManager
}

func newChat() *chatImpl {
	return &chatImpl{
		chatRepo:      &database.ChatRepo{},
		presetRepo:    &database.PresetRepo{},
		worldInfoRepo: &database.WorldInfoRepo{},
		regexRepo:     &database.RegexRuleRepo{},
		charRepo:      &database.CharacterRepo{},
		userRepo:      &database.UserRepo{},
		cacheManager:  cache.GetSessionCacheManager(),
	}
}

func (c *chatImpl) ListChatSessions(ctx context.Context, req *pb.ListChatSessionsRequest) (*pb.ListChatSessionsResponse, error) {
	// 获取并规范化分页参数
	page, pageSize := constrant.NormalizePagination(int(req.GetPage()), int(req.GetPageSize()))
	characterID := int(req.GetCharacterId())

	// 从数据库获取会话列表
	userId := jwt.GetUserId(ctx)
	sessions, total, err := c.chatRepo.ListSessions(userId, characterID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbSessions := make([]*pb.ChatSession, 0, len(sessions))
	for _, session := range sessions {
		pbSessions = append(pbSessions, convert.SessionEntityToPb(session))
	}

	return &pb.ListChatSessionsResponse{
		Sessions: pbSessions,
		Total:    int32(total),
	}, nil
}

func (c *chatImpl) GetChatSession(ctx context.Context, req *pb.GetChatSessionRequest) (
	*pb.GetChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidSessionID)
	}
	userId := jwt.GetUserId(ctx)
	// 从数据库获取会话
	session, err := c.chatRepo.GetSessionByID(id, userId)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.SessionNotFound)
	}

	return &pb.GetChatSessionResponse{
		Session: convert.SessionEntityToPb(session),
	}, nil
}

func (c *chatImpl) CreateChatSession(ctx context.Context, req *pb.CreateChatSessionRequest) (*pb.CreateChatSessionResponse, error) {
	// 参数校验
	characterID := int(req.GetCharacterId())
	if characterID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidCharacterID)
	}

	userId := jwt.GetUserId(ctx)

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		name = "新会话"
	}

	// 构建会话实体
	session := &entity.ChatSession{
		UserID:      userId,
		CharacterID: characterID,
		Name:        name,
		Version:     1,
	}

	// 保存到数据库
	if err := c.chatRepo.CreateSession(session); err != nil {
		return nil, err
	}

	// 重新获取完整数据（包含关联）
	fullSession, err := c.chatRepo.GetSessionByID(session.ID, userId)
	if err != nil {
		return nil, err
	}

	return &pb.CreateChatSessionResponse{
		Session: convert.SessionEntityToPb(fullSession),
	}, nil
}

func (c *chatImpl) UpdateChatSession(ctx context.Context, req *pb.UpdateChatSessionRequest) (*pb.UpdateChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidSessionID)
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptySessionName)
	}
	userId := jwt.GetUserId(ctx)
	// 更新数据库
	_, err := c.chatRepo.UpdateSession(int(req.GetId()), userId, int(req.GetVersion()), req.GetName())
	if err != nil {
		return nil, err
	}
	// 重新获取更新后的会话（包含关联数据）
	updatedSession, err := c.chatRepo.GetSessionByID(id, userId)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateChatSessionResponse{
		Session: convert.SessionEntityToPb(updatedSession),
	}, nil
}

func (c *chatImpl) DeleteChatSession(ctx context.Context, req *pb.DeleteChatSessionRequest) (*pb.DeleteChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidSessionID)
	}
	userId := jwt.GetUserId(ctx)

	// 删除会话（会级联删除关联的消息和swipes）
	if err := c.chatRepo.DeleteSession(id, userId); err != nil {
		return nil, err
	}

	return &pb.DeleteChatSessionResponse{}, nil
}

func (c *chatImpl) SendMessage(ctx context.Context, req *pb.SendMessageRequest, stream SendMessageStream) error {
	// 参数校验
	sessionID := int(req.GetSessionId())
	if sessionID <= 0 {
		return errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidSessionID)
	}
	content := strings.TrimSpace(req.GetContent())
	if content == "" {
		return errs.NewStandard(connect.CodeInvalidArgument, "消息内容不能为空")
	}

	userID := jwt.GetUserId(ctx)

	// 获取用户信息以获取活跃预设ID
	user, err := c.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errs.NewStandard(connect.CodeNotFound, "用户不存在")
	}

	// 获取活跃的 API 配置
	apiConfig, err := c.userRepo.GetActiveAPIConfig(userID)
	if err != nil {
		return err
	}
	if apiConfig == nil {
		return errs.NewStandard(connect.CodeFailedPrecondition, "请先配置并激活 API")
	}

	// 尝试从缓存获取会话数据
	sessionCache, cacheHit := c.cacheManager.Get(int64(userID), int64(sessionID))

	// 如果缓存命中，校验版本是否有效
	if cacheHit {
		versions, verifyErr := c.getVersionInfo(userID, sessionCache)
		if verifyErr != nil {
			return verifyErr
		}
		if !sessionCache.IsVersionValid(versions) {
			// 版本不匹配，缓存失效
			cacheHit = false
			c.cacheManager.Remove(int64(userID), int64(sessionID))
		}
	}

	var session *entity.ChatSession
	var preset *entity.Preset
	var promptItems []*entity.PromptItem
	var globalWorldInfos []*entity.WorldInfo
	var charWorldInfo *entity.WorldInfo
	var regexRules []*entity.RegexRule
	var worldInfoByDepth map[int][]string

	if cacheHit {
		// 使用缓存数据
		session = &entity.ChatSession{
			ID:          int(sessionCache.SessionID),
			UserID:      int(sessionCache.UserID),
			CharacterID: int(sessionCache.CharacterID),
			Character:   sessionCache.Character,
			Messages:    make([]entity.Message, len(sessionCache.Messages)),
		}
		for i, msg := range sessionCache.Messages {
			session.Messages[i] = *msg
		}

		preset = sessionCache.Preset
		promptItems = sessionCache.PromptItems
		regexRules = sessionCache.RegexRules
		globalWorldInfos = sessionCache.WorldInfos
		worldInfoByDepth = sessionCache.WorldInfoByDepth

		// 角色卡世界书从世界书列表中查找
		if sessionCache.Character != nil && sessionCache.Character.WorldInfoID > 0 {
			for _, wi := range sessionCache.WorldInfos {
				if wi.ID == sessionCache.Character.WorldInfoID {
					charWorldInfo = wi
					break
				}
			}
		}
	} else {
		// 缓存未命中，从数据库加载数据
		session, err = c.chatRepo.GetSessionWithMessages(sessionID, userID)
		if err != nil {
			return err
		}
		if session == nil {
			return errs.NewStandard(connect.CodeNotFound, errs.SessionNotFound)
		}

		// 获取预设（包含提示项）
		if user.ActivePresetID > 0 {
			preset, err = c.presetRepo.GetByID(user.ActivePresetID, userID)
			if err != nil {
				return err
			}
		}
		if preset == nil {
			// 使用默认预设
			preset = &entity.Preset{
				Temperature:    1.0,
				TopP:           1.0,
				MaxTokens:      2048,
				CandidateCount: 1,
			}
		}

		// 加载预设的提示项
		promptItems, err = c.presetRepo.ListPromptItems(preset.ID)
		if err != nil {
			return err
		}

		// 加载全局世界书
		globalWorldInfos, err = c.worldInfoRepo.ListGlobalWorldInfosWithEntries(userID)
		if err != nil {
			return err
		}

		// 加载角色卡关联的世界书
		if session.Character != nil && session.Character.WorldInfoID > 0 {
			charWorldInfo, err = c.worldInfoRepo.GetByIDWithEntries(session.Character.WorldInfoID, userID)
			if err != nil {
				return err
			}
		}

		// 加载正则规则
		regexRules, err = c.regexRepo.ListEnabledRules(preset.ID, session.CharacterID)
		if err != nil {
			return err
		}

		// 预计算世界书按深度分组
		worldInfoByDepth = c.buildWorldInfoByDepth(globalWorldInfos, charWorldInfo)

		// 构建缓存
		sessionCache = c.buildSessionCache(session, preset, promptItems, globalWorldInfos, charWorldInfo, regexRules, worldInfoByDepth)
		c.cacheManager.Set(int64(userID), int64(sessionID), sessionCache)
	}

	// 构建消息列表
	messages := c.buildMessages(session, content, promptItems, globalWorldInfos, charWorldInfo, regexRules)

	// 保存用户消息到数据库
	maxOrder, err := c.chatRepo.GetMaxMessageSortOrder(sessionID)
	if err != nil {
		return err
	}
	userMessage := &entity.Message{
		SessionID:        sessionID,
		Role:             pb.Role_User,
		ActiveSwipeIndex: 0,
		SortOrder:        maxOrder + 1,
	}
	if err = c.chatRepo.CreateMessage(userMessage); err != nil {
		return err
	}
	userSwipe := &entity.MessageSwipe{
		MessageID: userMessage.ID,
		Content:   content,
		SortOrder: 0,
	}
	if err = c.chatRepo.CreateMessageSwipe(userSwipe); err != nil {
		return err
	}

	// 更新缓存中的消息列表
	userMessage.Swipes = []entity.MessageSwipe{*userSwipe}
	sessionCache.AddMessage(userMessage)

	// 创建 AI 消息占位
	aiMessage := &entity.Message{
		SessionID:        sessionID,
		Role:             pb.Role_Assistant,
		ActiveSwipeIndex: 0,
		SortOrder:        maxOrder + 2,
	}
	if err = c.chatRepo.CreateMessage(aiMessage); err != nil {
		return err
	}

	// 根据候选数量创建对应数量的 swipe
	aiSwipes := make([]*entity.MessageSwipe, preset.CandidateCount)
	for i := 0; i < preset.CandidateCount; i++ {
		aiSwipe := &entity.MessageSwipe{
			MessageID: aiMessage.ID,
			Content:   "",
			SortOrder: i,
		}
		if err = c.chatRepo.CreateMessageSwipe(aiSwipe); err != nil {
			return err
		}
		aiSwipes[i] = aiSwipe
	}

	// 选择 LLM 模型并调用
	llm := c.getLLM(apiConfig.Provider)

	// 流式调用大模型
	resultChan := llm.StreamGenerateContent(ctx, apiConfig.APIKey, apiConfig.Model, *preset, messages)

	// 收集每个候选的内容
	contents := make([]string, preset.CandidateCount)

	// 处理流式响应
	for result := range resultChan {
		if result.Error != nil {
			return result.Error
		}

		if result.Index < len(contents) {
			contents[result.Index] += result.Content
		}

		// 发送流式响应给前端
		if sendErr := stream.Send(&pb.SendMessageResponse{
			Index:      int32(result.Index),
			Content:    result.Content,
			IsComplete: result.Done,
		}); sendErr != nil {
			return errs.NewStandardf(connect.CodeInternal, "发送响应失败: %v", sendErr)
		}

		if result.Done {
			break
		}
	}

	// 应用正则规则到 AI 输出
	for i, content := range contents {
		contents[i] = c.applyRegexRules(content, regexRules, RegexTextTypeAIOutput)
	}

	// 更新数据库中的 AI swipe 内容
	for i, content := range contents {
		if i < len(aiSwipes) {
			aiSwipes[i].Content = content
			if err = c.chatRepo.UpdateMessageSwipe(aiSwipes[i].ID, content); err != nil {
				return err
			}
		}
	}

	// 更新缓存中的 AI 消息
	aiMessage.Swipes = make([]entity.MessageSwipe, len(aiSwipes))
	for i, swipe := range aiSwipes {
		aiMessage.Swipes[i] = *swipe
	}
	sessionCache.AddMessage(aiMessage)

	return nil
}

// getVersionInfo 获取缓存相关实体的版本信息
func (c *chatImpl) getVersionInfo(userID int, sessionCache *cache.SessionCache) (*cache.VersionInfo, *connect.Error) {
	versions := &cache.VersionInfo{
		WorldInfoVersions: make(map[int64]int),
		RegexRuleVersions: make(map[int64]int),
	}

	// 获取角色卡版本
	if sessionCache.CharacterID > 0 {
		charVersion, err := c.charRepo.GetVersion(int(sessionCache.CharacterID), userID)
		if err != nil {
			return nil, err
		}
		versions.CharacterVersion = charVersion
	}

	// 获取预设版本
	if sessionCache.Preset != nil && sessionCache.Preset.ID > 0 {
		presetVersion, err := c.presetRepo.GetVersion(sessionCache.Preset.ID, userID)
		if err != nil {
			return nil, err
		}
		versions.PresetVersion = presetVersion
	}

	// 获取世界书版本
	if len(sessionCache.WorldInfoVersions) > 0 {
		worldInfoIDs := make([]int, 0, len(sessionCache.WorldInfoVersions))
		for id := range sessionCache.WorldInfoVersions {
			worldInfoIDs = append(worldInfoIDs, int(id))
		}
		worldInfoVersions, err := c.worldInfoRepo.GetVersions(worldInfoIDs, userID)
		if err != nil {
			return nil, err
		}
		versions.WorldInfoVersions = worldInfoVersions
	}

	// 获取正则规则版本
	presetID := 0
	if sessionCache.Preset != nil {
		presetID = sessionCache.Preset.ID
	}
	regexVersions, err := c.regexRepo.GetEnabledRuleVersions(presetID, int(sessionCache.CharacterID))
	if err != nil {
		return nil, err
	}
	versions.RegexRuleVersions = regexVersions

	return versions, nil
}

// buildWorldInfoByDepth 预计算世界书按深度分组
func (c *chatImpl) buildWorldInfoByDepth(globalWorldInfos []*entity.WorldInfo, charWorldInfo *entity.WorldInfo) map[int][]string {
	worldInfoByDepth := make(map[int][]string)

	// 收集全局世界书常驻条目（按深度分组）
	for _, wi := range globalWorldInfos {
		for _, entry := range wi.Entries {
			if entry.Constant {
				worldInfoByDepth[entry.Depth] = append(worldInfoByDepth[entry.Depth], entry.Content)
			}
		}
	}

	// 收集角色世界书常驻条目（按深度分组）
	if charWorldInfo != nil {
		for _, entry := range charWorldInfo.Entries {
			if entry.Constant {
				worldInfoByDepth[entry.Depth] = append(worldInfoByDepth[entry.Depth], entry.Content)
			}
		}
	}

	return worldInfoByDepth
}

// buildSessionCache 构建会话缓存
func (c *chatImpl) buildSessionCache(
	session *entity.ChatSession,
	preset *entity.Preset,
	promptItems []*entity.PromptItem,
	globalWorldInfos []*entity.WorldInfo,
	charWorldInfo *entity.WorldInfo,
	regexRules []*entity.RegexRule,
	worldInfoByDepth map[int][]string,
) *cache.SessionCache {
	sessionCache := &cache.SessionCache{
		SessionID:         int64(session.ID),
		UserID:            int64(session.UserID),
		CharacterID:       int64(session.CharacterID),
		Character:         session.Character,
		Preset:            preset,
		PromptItems:       promptItems,
		RegexRules:        regexRules,
		WorldInfoByDepth:  worldInfoByDepth,
		WorldInfoVersions: make(map[int64]int),
		RegexRuleVersions: make(map[int64]int),
	}

	// 设置版本号
	if session.Character != nil {
		sessionCache.CharacterVersion = session.Character.Version
	}
	if preset != nil {
		sessionCache.PresetVersion = preset.Version
	}

	// 收集世界书
	allWorldInfos := make([]*entity.WorldInfo, 0, len(globalWorldInfos)+1)
	allWorldInfos = append(allWorldInfos, globalWorldInfos...)
	if charWorldInfo != nil {
		allWorldInfos = append(allWorldInfos, charWorldInfo)
	}
	sessionCache.WorldInfos = allWorldInfos

	// 收集世界书条目
	var allEntries []*entity.WorldInfoEntry
	for _, wi := range allWorldInfos {
		sessionCache.WorldInfoVersions[int64(wi.ID)] = wi.Version
		for i := range wi.Entries {
			allEntries = append(allEntries, &wi.Entries[i])
		}
	}
	sessionCache.WorldInfoEntries = allEntries

	// 收集正则规则版本
	for _, rule := range regexRules {
		sessionCache.RegexRuleVersions[int64(rule.ID)] = rule.Version
	}

	// 复制消息列表
	sessionCache.Messages = make([]*entity.Message, len(session.Messages))
	for i := range session.Messages {
		sessionCache.Messages[i] = &session.Messages[i]
	}

	return sessionCache
}

// InvalidateSessionCache 使指定会话的缓存失效
func (c *chatImpl) InvalidateSessionCache(userID, sessionID int64) {
	c.cacheManager.Remove(userID, sessionID)
}

// InvalidateCacheByCharacter 使指定角色相关的所有缓存失效
func (c *chatImpl) InvalidateCacheByCharacter(characterID int64) {
	c.cacheManager.RemoveByCharacter(characterID)
}

// InvalidateCacheByPreset 使使用指定预设的所有缓存失效
func (c *chatImpl) InvalidateCacheByPreset(presetID int64) {
	c.cacheManager.RemoveByPreset(presetID)
}

// InvalidateCacheByWorldInfo 使使用指定世界书的所有缓存失效
func (c *chatImpl) InvalidateCacheByWorldInfo(worldInfoID int64) {
	c.cacheManager.RemoveByWorldInfo(worldInfoID)
}

// InvalidateCacheByUser 使指定用户的所有缓存失效
func (c *chatImpl) InvalidateCacheByUser(userID int64) {
	c.cacheManager.RemoveByUser(userID)
}

// GetCacheStats 获取缓存统计信息
func (c *chatImpl) GetCacheStats() int {
	return c.cacheManager.Size()
}

// buildMessages 构建发送给大模型的消息列表
func (c *chatImpl) buildMessages(
	session *entity.ChatSession,
	userContent string,
	promptItems []*entity.PromptItem,
	globalWorldInfos []*entity.WorldInfo,
	charWorldInfo *entity.WorldInfo,
	regexRules []*entity.RegexRule,
) []model.Message {
	var messages []model.Message

	// 按深度分组的世界书条目
	worldInfoByDepth := make(map[int][]string)

	// 收集全局世界书条目
	for _, wi := range globalWorldInfos {
		for _, entry := range wi.Entries {
			if entry.Constant || c.matchWorldInfoKeys(entry, session, userContent) {
				worldInfoByDepth[entry.Depth] = append(worldInfoByDepth[entry.Depth], entry.Content)
			}
		}
	}

	// 收集角色世界书条目
	if charWorldInfo != nil {
		for _, entry := range charWorldInfo.Entries {
			if entry.Constant || c.matchWorldInfoKeys(entry, session, userContent) {
				worldInfoByDepth[entry.Depth] = append(worldInfoByDepth[entry.Depth], entry.Content)
			}
		}
	}

	// 分离不同位置的提示项
	var systemPrompts []string
	var beforeCharPrompts []string
	var afterCharPrompts []string
	depthPrompts := make(map[int][]string)

	for _, item := range promptItems {
		if !item.IsEnabled {
			continue
		}
		content := c.applyRegexRules(item.Content, regexRules, RegexTextTypePrompt)

		switch item.InjectionPosition {
		case pb.InjectionPosition_Relative:
			// 相对位置：根据角色添加到对应位置
			switch item.Role {
			case pb.Role_System:
				systemPrompts = append(systemPrompts, content)
			case pb.Role_User:
				beforeCharPrompts = append(beforeCharPrompts, content)
			case pb.Role_Assistant:
				afterCharPrompts = append(afterCharPrompts, content)
			}
		case pb.InjectionPosition_Absolute:
			// 绝对位置：按深度插入
			depthPrompts[item.InjectionDepth] = append(depthPrompts[item.InjectionDepth], content)
		}
	}

	// 构建系统提示（包含角色描述）
	var systemContent strings.Builder
	for _, prompt := range systemPrompts {
		systemContent.WriteString(prompt)
		systemContent.WriteString("\n\n")
	}

	// 添加角色描述
	if session.Character != nil {
		for _, prompt := range beforeCharPrompts {
			systemContent.WriteString(prompt)
			systemContent.WriteString("\n\n")
		}
		if session.Character.Description != "" {
			systemContent.WriteString(session.Character.Description)
			systemContent.WriteString("\n\n")
		}
		for _, prompt := range afterCharPrompts {
			systemContent.WriteString(prompt)
			systemContent.WriteString("\n\n")
		}
		// 添加示例对话
		if session.Character.ExampleDialogue != "" {
			systemContent.WriteString("示例对话:\n")
			systemContent.WriteString(session.Character.ExampleDialogue)
			systemContent.WriteString("\n\n")
		}
	}

	// 添加深度0的世界书和提示项到系统消息
	for _, content := range worldInfoByDepth[0] {
		systemContent.WriteString(content)
		systemContent.WriteString("\n\n")
	}
	for _, content := range depthPrompts[0] {
		systemContent.WriteString(content)
		systemContent.WriteString("\n\n")
	}

	if systemContent.Len() > 0 {
		messages = append(messages, model.Message{
			Role:    pb.Role_System,
			Content: strings.TrimSpace(systemContent.String()),
		})
	}

	// 添加角色开场白
	if session.Character != nil && len(session.Character.FirstMessage) != 0 && len(session.Messages) == 0 {
		firstMsg := c.applyRegexRules(session.Character.FirstMessage[0], regexRules, RegexTextTypePrompt)
		messages = append(messages, model.Message{
			Role:    pb.Role_Assistant,
			Content: firstMsg,
		})
	}

	// 添加历史消息（按深度插入世界书条目）
	historyMessages := c.buildHistoryMessages(session.Messages, worldInfoByDepth, depthPrompts, regexRules)
	messages = append(messages, historyMessages...)

	// 添加用户新消息
	processedUserContent := c.applyRegexRules(userContent, regexRules, RegexTextTypeUserInput)
	messages = append(messages, model.Message{
		Role:    pb.Role_User,
		Content: processedUserContent,
	})

	return messages
}

// buildHistoryMessages 构建历史消息列表，并在对应深度插入世界书条目
func (c *chatImpl) buildHistoryMessages(
	historyMsgs []entity.Message,
	worldInfoByDepth map[int][]string,
	depthPrompts map[int][]string,
	regexRules []*entity.RegexRule,
) []model.Message {
	var messages []model.Message

	// 计算每条消息的深度（从末尾开始计算）
	totalMessages := len(historyMsgs)

	for i, msg := range historyMsgs {
		depth := totalMessages - i

		// 在对应深度插入世界书条目和提示项
		if contents, ok := worldInfoByDepth[depth]; ok {
			for _, content := range contents {
				messages = append(messages, model.Message{
					Role:    pb.Role_System,
					Content: content,
				})
			}
		}
		if contents, ok := depthPrompts[depth]; ok {
			for _, content := range contents {
				messages = append(messages, model.Message{
					Role:    pb.Role_System,
					Content: content,
				})
			}
		}

		// 添加消息内容（使用当前活跃的 swipe）
		if len(msg.Swipes) > 0 {
			swipeIndex := msg.ActiveSwipeIndex
			if swipeIndex >= len(msg.Swipes) {
				swipeIndex = 0
			}
			content := msg.Swipes[swipeIndex].Content

			// 对历史消息应用正则
			if msg.Role == pb.Role_Assistant {
				content = c.applyRegexRules(content, regexRules, RegexTextTypeAIOutput)
			} else {
				content = c.applyRegexRules(content, regexRules, RegexTextTypeUserInput)
			}

			messages = append(messages, model.Message{
				Role:    msg.Role,
				Content: content,
			})
		}
	}

	return messages
}

// matchWorldInfoKeys 检查世界书条目的关键词是否匹配
func (c *chatImpl) matchWorldInfoKeys(entry entity.WorldInfoEntry, session *entity.ChatSession, userContent string) bool {
	// 将所有内容合并用于匹配
	var searchContent strings.Builder
	searchContent.WriteString(userContent)

	// 添加最近几条消息内容用于匹配
	for i := len(session.Messages) - 1; i >= 0 && i >= len(session.Messages)-5; i-- {
		msg := session.Messages[i]
		if len(msg.Swipes) > 0 {
			swipeIndex := msg.ActiveSwipeIndex
			if swipeIndex >= len(msg.Swipes) {
				swipeIndex = 0
			}
			searchContent.WriteString(" ")
			searchContent.WriteString(msg.Swipes[swipeIndex].Content)
		}
	}

	searchText := strings.ToLower(searchContent.String())

	// 分割关键词列表
	keys := strings.Split(entry.KeysList, ",")
	for _, key := range keys {
		key = strings.TrimSpace(strings.ToLower(key))
		if key != "" && strings.Contains(searchText, key) {
			// 如果有次要关键词，需要同时匹配
			if entry.Selective && entry.SecondaryKeys != "" {
				secondaryKeys := strings.Split(entry.SecondaryKeys, ",")
				for _, sk := range secondaryKeys {
					sk = strings.TrimSpace(strings.ToLower(sk))
					if sk != "" && strings.Contains(searchText, sk) {
						return true
					}
				}
				return false
			}
			return true
		}
	}
	return false
}

// applyRegexRules 应用正则规则到文本
// RegexTextType 正则规则作用的文本类型
type RegexTextType int

const (
	RegexTextTypeUserInput    RegexTextType = iota // 用户输入
	RegexTextTypeAIOutput                          // AI输出
	RegexTextTypeSlashCommand                      // 斜杠命令
	RegexTextTypeWorldInfo                         // 世界书
	RegexTextTypePrompt                            // 提示词
)

// applyRegexRules 应用正则规则到文本
// text: 待处理的文本
// rules: 正则规则列表
// textType: 文本类型，用于匹配正则规则的作用范围
func (c *chatImpl) applyRegexRules(text string, rules []*entity.RegexRule, textType RegexTextType) string {
	// 按排序顺序处理正则规则
	sortedRules := make([]*entity.RegexRule, len(rules))
	copy(sortedRules, rules)
	sort.Slice(sortedRules, func(i, j int) bool {
		return sortedRules[i].SortOrder < sortedRules[j].SortOrder
	})

	for _, rule := range sortedRules {
		// 根据文本类型检查规则是否适用
		if !c.isRuleApplicable(rule, textType) {
			continue
		}

		// 编译正则表达式
		re, err := regexp.Compile(rule.FindPattern)
		if err != nil {
			continue
		}

		// 执行替换
		if rule.SubstituteRegex {
			text = re.ReplaceAllString(text, rule.ReplacePattern)
		} else {
			text = re.ReplaceAllLiteralString(text, rule.ReplacePattern)
		}
	}

	return text
}

// isRuleApplicable 检查正则规则是否适用于指定的文本类型
func (c *chatImpl) isRuleApplicable(rule *entity.RegexRule, textType RegexTextType) bool {
	switch textType {
	case RegexTextTypeUserInput:
		return rule.AffectFlagsUserInput
	case RegexTextTypeAIOutput:
		return rule.AffectFlagsAIOutput
	case RegexTextTypeSlashCommand:
		return rule.AffectFlagsSlashCommand
	case RegexTextTypeWorldInfo:
		return rule.AffectFlagsWorldInfo
	case RegexTextTypePrompt:
		return rule.AffectFlagsPrompt
	default:
		return false
	}
}

// getLLM 根据提供商获取对应的 LLM 实例
func (c *chatImpl) getLLM(provider pb.APIProvider) model.LLMModel {
	switch provider {
	case pb.APIProvider_Gemini:
		return model.NewGemini()
	case pb.APIProvider_OpenAI:
		// TODO: 实现 OpenAI 模型
		return model.NewGemini() // 暂时返回 Gemini
	case pb.APIProvider_Claude:
		// TODO: 实现 Claude 模型
		return model.NewGemini() // 暂时返回 Gemini
	default:
		return model.NewGemini()
	}
}

func (c *chatImpl) RegenerateMessage(ctx context.Context, req *pb.RegenerateMessageRequest, stream RegenerateMessageStream) error {
	//TODO implement me
	panic("implement me")
}

func (c *chatImpl) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *chatImpl) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *chatImpl) SwitchSwipe(ctx context.Context, req *pb.SwitchSwipeRequest) (*pb.SwitchSwipeResponse, error) {
	//TODO implement me
	panic("implement me")
}
