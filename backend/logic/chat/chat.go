package chat

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/dlclark/regexp2"
	"github.com/ling/muse/common/constant"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/common/jwt"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/cache"
	"github.com/ling/muse/repo/database"
	"github.com/ling/muse/repo/model"
	log "github.com/sirupsen/logrus"
)

type chatImpl struct {
	chatRepo      database.ChatRepository
	presetRepo    database.PresetRepository
	worldInfoRepo database.WorldInfoRepository
	regexRepo     database.RegexRuleRepository
	charRepo      database.CharacterRepository
	userRepo      database.UserRepository

	// 会话缓存管理器
	cacheManager *cache.SessionCacheManager
}

func newChat() *chatImpl {
	return &chatImpl{
		chatRepo:      database.NewChatRepo(),
		presetRepo:    database.NewPresetRepo(),
		worldInfoRepo: database.NewWorldInfoRepo(),
		regexRepo:     database.NewRegexRuleRepo(),
		charRepo:      database.NewCharacterRepo(),
		userRepo:      database.NewUserRepo(),
		cacheManager:  cache.GetSessionCacheManager(),
	}
}

func (c *chatImpl) ListChatSessions(ctx context.Context, req *pb.ListChatSessionsRequest) (*pb.ListChatSessionsResponse, error) {
	// 获取并规范化分页参数
	page, pageSize := constant.NormalizePagination(int(req.GetPage()), int(req.GetPageSize()))
	characterID := int(req.GetCharacterId())

	// 从数据库获取会话列表
	userId := jwt.GetUserId(ctx)
	sessions, total, err := c.chatRepo.ListSessions(ctx, userId, characterID, page, pageSize)
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
	session, err := c.chatRepo.GetSessionByID(ctx, id, userId)
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

	fullSession, err := c.createSession(ctx, int(req.GetCharacterId()), req.GetName())
	if err != nil {
		return nil, err
	}

	return &pb.CreateChatSessionResponse{
		Session: convert.SessionEntityToPb(fullSession),
	}, nil
}

func (c *chatImpl) createSession(ctx context.Context, characterID int, name string) (*entity.ChatSession, error) {
	userId := jwt.GetUserId(ctx)

	// 插入第一条消息
	character, err := c.charRepo.GetByID(ctx, characterID, userId)
	if err != nil {
		log.Error("failed to get character")
		return nil, err
	}
	if character == nil {
		log.Infof("未找到角色卡")
		return nil, errs.NewStandardf(connect.CodeNotFound, errs.CharacterNotFound)
	}
	if name == "" {
		name = fmt.Sprintf("%s_%s", character.Name, time.Now().Format(time.DateTime))
	}
	swipes := make([]entity.MessageSwipe, len(character.FirstMessage))
	for i, content := range character.FirstMessage {
		swipes[i] = entity.MessageSwipe{
			Content:   content,
			SortOrder: i * constant.SortOrderInterval,
		}
	}

	message := []entity.Message{{
		Role:             pb.Role_Assistant,
		ActiveSwipeIndex: 0,
		SortOrder:        0,
		Swipes:           swipes,
	}}

	// 构建会话实体
	session := &entity.ChatSession{
		UserID:      userId,
		CharacterID: characterID,
		Name:        name,
		Version:     1,
		Messages:    message,
	}

	// 保存到数据库
	if err = c.chatRepo.CreateSession(ctx, session); err != nil {
		log.Error("failed to create session")
		return nil, err
	}

	// 重新获取完整数据（包含关联）
	fullSession, err := c.chatRepo.GetSessionByID(ctx, session.ID, userId)
	if err != nil {
		log.Error("failed to get session")
		return nil, err
	}
	return fullSession, nil
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
	_, err := c.chatRepo.UpdateSession(ctx, int(req.GetId()), userId, int(req.GetVersion()), req.GetName())
	if err != nil {
		return nil, err
	}
	// 重新获取更新后的会话（包含关联数据）
	updatedSession, err := c.chatRepo.GetSessionByID(ctx, id, userId)
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
	if err := c.chatRepo.DeleteSession(ctx, id, userId); err != nil {
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
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errs.NewStandard(connect.CodeNotFound, "用户不存在")
	}

	// 尝试从缓存获取会话数据
	sessionCache, cacheHit := c.cacheManager.Get(int64(userID), int64(sessionID))
	// 如果缓存命中，校验版本是否有效
	if cacheHit {
		versions, verifyErr := c.getVersionInfo(ctx, userID, sessionCache)
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
		session, err = c.chatRepo.GetSessionWithMessages(ctx, sessionID, userID)
		if err != nil {
			return err
		}
		if session == nil {
			return errs.NewStandard(connect.CodeNotFound, errs.SessionNotFound)
		}

		// 获取预设（包含提示项）
		if user.ActivePresetID > 0 {
			preset, err = c.presetRepo.GetByID(ctx, user.ActivePresetID, userID)
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
		promptItems, err = c.presetRepo.ListPromptItems(ctx, preset.ID)
		if err != nil {
			return err
		}

		// 加载全局世界书
		globalWorldInfos, err = c.worldInfoRepo.ListGlobalWorldInfosWithEntries(ctx, userID)
		if err != nil {
			return err
		}

		// 加载角色卡关联的世界书
		if session.Character != nil && session.Character.WorldInfoID > 0 {
			charWorldInfo, err = c.worldInfoRepo.GetByIDWithEntries(ctx, session.Character.WorldInfoID, userID)
			if err != nil {
				return err
			}
		}

		// 加载正则规则
		regexRules, err = c.regexRepo.ListEnabledRules(ctx, preset.ID, session.CharacterID)
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
	worldBook := append(globalWorldInfos, charWorldInfo)
	log.Debugf("会话: %+v, 世界书: %+v, 预设: %+v, 正则: %+v",
		session, globalWorldInfos, preset, regexRules)
	messages, err := c.buildMessages(session, promptItems, worldBook, regexRules)
	log.Debugf("构建的消息列表: %v", messages)
	if err != nil {
		return err
	}
	// 保存用户消息到数据库
	maxOrder, err := c.chatRepo.GetMaxMessageSortOrder(ctx, sessionID)
	if err != nil {
		return err
	}
	userMessage := &entity.Message{
		SessionID:        sessionID,
		Role:             pb.Role_User,
		ActiveSwipeIndex: 0,
		SortOrder:        maxOrder + 1,
	}
	if err = c.chatRepo.CreateMessage(ctx, userMessage); err != nil {
		return err
	}
	userSwipe := &entity.MessageSwipe{
		MessageID: userMessage.ID,
		Content:   content,
		SortOrder: 0,
	}
	if err = c.chatRepo.CreateMessageSwipe(ctx, userSwipe); err != nil {
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
	if err = c.chatRepo.CreateMessage(ctx, aiMessage); err != nil {
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
		if err = c.chatRepo.CreateMessageSwipe(ctx, aiSwipe); err != nil {
			return err
		}
		aiSwipes[i] = aiSwipe
	}

	// 选择 LLM 模型并调用
	// 获取活跃的 API 配置
	apiConfig, err := c.userRepo.GetActiveAPIConfig(ctx, userID)
	if err != nil {
		return err
	}
	if apiConfig == nil {
		return errs.NewStandard(connect.CodeFailedPrecondition, "请先配置并激活 API")
	}
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
			if err = c.chatRepo.UpdateMessageSwipe(ctx, aiSwipes[i].ID, content); err != nil {
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
func (c *chatImpl) getVersionInfo(ctx context.Context, userID int, sessionCache *cache.SessionCache) (*cache.VersionInfo, error) {
	versions := &cache.VersionInfo{
		WorldInfoVersions: make(map[int64]int),
		RegexRuleVersions: make(map[int64]int),
	}

	// 获取角色卡版本
	if sessionCache.CharacterID > 0 {
		charVersion, err := c.charRepo.GetVersion(ctx, int(sessionCache.CharacterID), userID)
		if err != nil {
			return nil, err
		}
		versions.CharacterVersion = charVersion
	}

	// 获取预设版本
	if sessionCache.Preset != nil && sessionCache.Preset.ID > 0 {
		presetVersion, err := c.presetRepo.GetVersion(ctx, sessionCache.Preset.ID, userID)
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
		worldInfoVersions, err := c.worldInfoRepo.GetVersions(ctx, worldInfoIDs, userID)
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
	regexVersions, err := c.regexRepo.GetEnabledRuleVersions(ctx, presetID, int(sessionCache.CharacterID))
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

	for _, key := range entry.Keys {
		key = strings.TrimSpace(strings.ToLower(key))
		if key != "" && strings.Contains(searchText, key) {
			// 如果有次要关键词，需要同时匹配
			if entry.Selective && len(entry.SecondaryKeys) > 0 {
				for _, sk := range entry.SecondaryKeys {
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
	err := c.chatRepo.SwitchSwipe(ctx, int(req.MessageId), int(req.SwipeIndex))
	if err != nil {
		if errs.Code(err) == connect.CodeNotFound {
			return nil, errs.NewStandardf(connect.CodeInvalidArgument,
				"参数错误, 消息或swipe不存在, messageID: %d", req.MessageId)
		}
		return nil, err
	}
	return &pb.SwitchSwipeResponse{}, nil
}

func (c *chatImpl) buildMessages(
	session *entity.ChatSession,
	promptItems []*entity.PromptItem,
	worldBook []*entity.WorldInfo,
	regexRules []*entity.RegexRule,
) ([]model.Message, error) {
	var messages []model.Message
	for _, promptItem := range promptItems {
		switch promptItem.Identifier {
		case pb.PromptItemIdentifier_PromptItemIdentifierUnspecified, pb.PromptItemIdentifier_Main,
			pb.PromptItemIdentifier_Jailbreak, pb.PromptItemIdentifier_Nsfw:
			messages = append(messages, model.Message{
				Role:    promptItem.Role,
				Content: promptItem.Content,
				Module:  constant.Preset,
			})
		case pb.PromptItemIdentifier_ChatHistory:
			var msgs []model.Message
			var systemCount, assistantCount, userCount int
			systemMsg, assistantMsg, userMsg := collectMsg(worldBook, promptItems, session)
			// 添加历史消息
			for i := len(session.Messages) - 1; i >= 0; i-- {
				history := session.Messages[i]
				if len(history.Swipes) == 0 {
					continue
				}
				switch history.Role {
				case pb.Role_System:
					if insertContent, exist := systemMsg[systemCount]; exist {
						msgs = append(msgs, insertContent...)
					}
					systemCount++
				case pb.Role_User:
					if insertContent, exist := userMsg[userCount]; exist {
						msgs = append(msgs, insertContent...)
					}
					userCount++
				case pb.Role_Assistant:
					if insertContent, exist := assistantMsg[assistantCount]; exist {
						msgs = append(msgs, insertContent...)
					}
					assistantCount++
				}
				index := history.ActiveSwipeIndex
				if index >= len(history.Swipes) {
					index = len(history.Swipes) - 1
				}
				message := model.Message{
					Role:    history.Role,
					Content: history.Swipes[index].Content,
				}
				switch history.Role {
				case pb.Role_User:
					message.Module = constant.UserInput
				case pb.Role_Assistant:
					message.Module = constant.AIOutput
				}
				msgs = append(msgs, message)
			}
			slices.Reverse(msgs)
			messages = append(messages, msgs...)
		case pb.PromptItemIdentifier_CharDescription, pb.PromptItemIdentifier_PersonaDescription:
			messages = append(messages, model.Message{
				Role:    promptItem.Role,
				Content: session.Character.Description,
				Module:  constant.Preset,
			})
		case pb.PromptItemIdentifier_DialogueExamples:
			for _, example := range session.Character.ExampleDialogue {
				messages = append(messages, model.Message{
					Role:    promptItem.Role,
					Content: example,
					Module:  constant.Preset,
				})
			}
		case pb.PromptItemIdentifier_WorldInfoBefore, pb.PromptItemIdentifier_WorldInfoAfter:
			for _, worldInfo := range worldBook {
				for _, entry := range worldInfo.Entries {
					if !(promptItem.Identifier == pb.PromptItemIdentifier_WorldInfoBefore &&
						entry.Position == pb.EntryPosition_BeforeChar ||
						promptItem.Identifier == pb.PromptItemIdentifier_WorldInfoAfter &&
							entry.Position == pb.EntryPosition_AfterChar) {
						continue
					}
					if !isWorldInfoEntryValid(&entry, session) {
						continue
					}
					messages = append(messages, model.Message{
						Role:    promptItem.Role,
						Content: entry.Content,
						Module:  constant.WorldInfo,
					})
				}
			}
		}
	}
	for _, regex := range regexRules {
		if !regex.IsEnabled {
			continue
		}
		re, err := regexp2.Compile(regex.FindPattern, 0)
		if err != nil {
			log.Warnf("正则规则编译失败: %v", err)
			continue
		}
		for i, message := range messages {
			if message.Module != constant.Preset && regex.AffectFlagsPrompt ||
				message.Module != constant.UserInput && regex.AffectFlagsUserInput ||
				message.Module != constant.AIOutput && regex.AffectFlagsAIOutput ||
				message.Module != constant.WorldInfo && regex.AffectFlagsWorldInfo {
				continue
			}
			// 检查消息深度是否在允许范围内
			if regex.MinDepth < len(messages)-i && regex.MaxDepth > len(messages)-i {
				continue
			}
			content, err := re.Replace(message.Content, regex.ReplacePattern, -1, -1)
			if err != nil {
				log.Warnf("正则规则替换内容失败: %v", err)
				return nil, errs.NewStandardf(connect.CodeInternal, "正则规则替换内容失败: %v", err)
			}
			messages[i].Content = content
		}
	}
	return messages, nil
}

// PriorityMessage 继承model.Message并添加优先级字段
type PriorityMessage struct {
	model.Message
	Priority int
}

func collectMsg(worldBooks []*entity.WorldInfo, promptItems []*entity.PromptItem, session *entity.ChatSession) (
	systemMsg, assistantMsg, userMsg map[int][]model.Message) {
	systemMsgTemp := make(map[int][]PriorityMessage)
	assistantMsgTemp := make(map[int][]PriorityMessage)
	userMsgTemp := make(map[int][]PriorityMessage)
	for _, worldBook := range worldBooks {
		for _, entry := range worldBook.Entries {
			if entry.Position != pb.EntryPosition_AtDepth {
				continue
			}
			if !isWorldInfoEntryValid(&entry, session) {
				continue
			}
			message := PriorityMessage{
				Message: model.Message{
					Role:    entry.Role,
					Content: entry.Content,
					Module:  constant.WorldInfo,
				},
				Priority: entry.SortOrder,
			}
			switch entry.Role {
			case pb.Role_System:
				systemMsgTemp[entry.Depth] = append(systemMsgTemp[entry.Depth], message)
			case pb.Role_User:
				userMsgTemp[entry.Depth] = append(userMsgTemp[entry.Depth], message)
			case pb.Role_Assistant:
				assistantMsgTemp[entry.Depth] = append(assistantMsgTemp[entry.Depth], message)
			}
		}
	}
	for _, item := range promptItems {
		message := PriorityMessage{
			Message: model.Message{
				Role:    item.Role,
				Content: item.Content,
				Module:  constant.Preset,
			},
			// TODO Priority: item.SortOrder,
		}
		switch item.Role {
		case pb.Role_System:
			systemMsgTemp[item.InjectionDepth] = append(systemMsgTemp[item.InjectionDepth], message)
		case pb.Role_User:
			userMsgTemp[item.InjectionDepth] = append(userMsgTemp[item.InjectionDepth], message)
		case pb.Role_Assistant:
			assistantMsgTemp[item.InjectionDepth] = append(assistantMsgTemp[item.InjectionDepth], message)
		}
	}

	// 为每个深度的数组按优先级进行排序，优先级越大越靠前
	for depth := range systemMsgTemp {
		sort.Slice(systemMsgTemp[depth], func(i, j int) bool {
			return systemMsgTemp[depth][i].Priority > systemMsgTemp[depth][j].Priority
		})
	}
	for depth := range assistantMsgTemp {
		sort.Slice(assistantMsgTemp[depth], func(i, j int) bool {
			return assistantMsgTemp[depth][i].Priority > assistantMsgTemp[depth][j].Priority
		})
	}
	for depth := range userMsgTemp {
		sort.Slice(userMsgTemp[depth], func(i, j int) bool {
			return userMsgTemp[depth][i].Priority > userMsgTemp[depth][j].Priority
		})
	}

	// 将排序后的PriorityMessage转换为model.Message并写入返回值
	systemMsg = make(map[int][]model.Message)
	assistantMsg = make(map[int][]model.Message)
	userMsg = make(map[int][]model.Message)
	for depth, messages := range systemMsgTemp {
		systemMsg[depth] = convertPriorityMessages(messages)
	}
	for depth, messages := range assistantMsgTemp {
		assistantMsg[depth] = convertPriorityMessages(messages)
	}
	for depth, messages := range userMsgTemp {
		userMsg[depth] = convertPriorityMessages(messages)
	}
	return systemMsg, assistantMsg, userMsg
}

// convertPriorityMessages 将PriorityMessage切片转换为model.Message切片
func convertPriorityMessages(messages []PriorityMessage) []model.Message {
	result := make([]model.Message, len(messages))
	for i, msg := range messages {
		result[i] = msg.Message
	}
	return result
}

func isWorldInfoEntryValid(entry *entity.WorldInfoEntry, session *entity.ChatSession) bool {
	if !entry.IsEnabled {
		return false
	}
	if entry.Constant {
		return true
	}
	if len(session.Messages) == 0 {
		return false
	}
	for _, message := range session.Messages {
		if len(message.Swipes) == 0 {
			continue
		}
		index := message.ActiveSwipeIndex
		if index >= len(message.Swipes) {
			index = len(message.Swipes) - 1
		}

		for _, key := range entry.Keys {
			if strings.Contains(message.Swipes[index].Content, key) {
				return true
			}
		}
	}
	return false
}

func (c *chatImpl) GetCharLatestSession(ctx context.Context, req *pb.GetCharLatestSessionRequest) (
	*pb.GetCharLatestSessionResponse, error) {
	if req.GetCharacterId() <= 0 {
		return nil, errs.NewStandardf(connect.CodeInvalidArgument, "character_id不能为空")
	}

	userId := jwt.GetUserId(ctx)
	session, err := c.chatRepo.GetCharLatestSessionWithMsg(ctx, int(req.GetCharacterId()), userId)
	if err != nil {
		return nil, err
	}
	if session == nil {
		session, err = c.createSession(ctx, int(req.GetCharacterId()), "")
		if err != nil {
			return nil, err
		}
	}
	return &pb.GetCharLatestSessionResponse{
		Session: convert.SessionEntityToPb(session),
	}, nil
}

func (c *chatImpl) UpdateSessionTime(ctx context.Context, req *pb.UpdateSessionTimeRequest) (
	*pb.UpdateSessionTimeResponse, error) {
	if err := c.chatRepo.UpdateSessionTime(ctx, int(req.GetSessionId()), time.Now()); err != nil {
		return nil, err
	}
	return &pb.UpdateSessionTimeResponse{}, nil
}
