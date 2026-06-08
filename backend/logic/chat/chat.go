package chat

import (
	"context"
	"errors"
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
	apiKeyCache  cache.APIKeyCache
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
		apiKeyCache:   cache.NewAPIKeyCache(),
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
			Content: content,
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
	var persona *entity.Persona

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
		persona = sessionCache.Persona

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
		session, preset, promptItems, globalWorldInfos,
			charWorldInfo, regexRules, persona, err = c.loadSessionDataFromDB(ctx, sessionID, userID, user)
		if err != nil {
			return errs.NewStandardf(errs.Code(err), "从数据库加载会话数据失败: %v", err)
		}
		// 预计算世界书按深度分组
		worldInfoByDepth := c.buildWorldInfoByDepth(globalWorldInfos, charWorldInfo)

		// 构建缓存
		sessionCache = c.buildSessionCache(session, preset, promptItems, globalWorldInfos,
			charWorldInfo, regexRules, worldInfoByDepth, persona)
		c.cacheManager.Set(int64(userID), int64(sessionID), sessionCache)
	}

	// 构建消息列表
	worldBook := append(globalWorldInfos, charWorldInfo)
	session.Messages = append(session.Messages, entity.Message{
		Role:             pb.Role_User,
		ActiveSwipeIndex: 0,
		Swipes: []entity.MessageSwipe{
			{
				Content: content,
			},
		},
	})
	messages, err := c.buildMessages(session, promptItems, worldBook, regexRules, persona)
	log.Debugf("构建的消息列表: %+v", messages)
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
		}
		if err = c.chatRepo.CreateMessageSwipe(ctx, aiSwipe); err != nil {
			return err
		}
		aiSwipes[i] = aiSwipe
	}

	// 选择 LLM 模型并调用
	// 获取活跃的 API 配置
	apiConfig := c.apiKeyCache.Get(userID, user.Provider)
	if apiConfig == nil {
		apiConfig = &entity.APIConfig{
			UserID:   userID,
			Provider: user.Provider,
			APIKey:   "",
			IsActive: true,
		}
		apiConfigs, err := c.userRepo.GetActiveAPIConfig(ctx, userID, user.Provider)
		if err != nil {
			return errs.NewStandardf(errs.Code(err), "获取API配置失败: %v", err)
		}
		if len(apiConfigs) != 0 {
			c.apiKeyCache.Set(userID, user.Provider, apiConfigs)
			apiConfig = apiConfigs[0]
		}
	}
	llm := c.getLLM(apiConfig.Provider)

	// 流式调用大模型
	resultChan := llm.StreamGenerateContent(ctx, apiConfig.APIKey, user.Model, user.ProxyUrl, *preset, messages)

	// 收集每个候选的内容
	contents := make([]string, preset.CandidateCount)

	// 处理流式响应
	for result := range resultChan {
		errCode := pb.ErrCode_Success
		errMsg := ""
		if result.Error != nil {
			log.Errorf("大模型输出失败: %v", result.Error)
			errCode = pb.ErrCode(connect.CodeInternal)
			errMsg = "大模型输出异常"

			var resultErr *connect.Error
			ok := errors.As(result.Error, &resultErr)
			if ok {
				errCode = pb.ErrCode(resultErr.Code())
				errMsg = resultErr.Message()
			}
		}

		if result.Index < len(contents) {
			contents[result.Index] += result.Content
		}
		//log.Debugf("大模型输出下标: %d, 内容: %s", result.Index, result.Content)

		// 发送流式响应给前端
		if sendErr := stream.Send(&pb.SendMessageResponse{
			Index:      int32(result.Index),
			Content:    result.Content,
			Done:       result.Done,
			ErrCode:    errCode,
			ErrMessage: errMsg,
		}); sendErr != nil {
			return errs.NewStandardf(connect.CodeInternal, "发送响应失败: %v", sendErr)
		}

		if result.Done {
			log.Info("大模型输出完成")
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

// loadSessionDataFromDB 从数据库加载会话数据（缓存未命中时调用）
func (c *chatImpl) loadSessionDataFromDB(ctx context.Context, sessionID, userID int, user *entity.User) (
	session *entity.ChatSession,
	preset *entity.Preset,
	promptItems []*entity.PromptItem,
	globalWorldInfos []*entity.WorldInfo,
	charWorldInfo *entity.WorldInfo,
	regexRules []*entity.RegexRule,
	persona *entity.Persona,
	err error,
) {
	// 从数据库获取会话数据
	session, err = c.chatRepo.GetSessionWithMessages(ctx, sessionID, userID)
	if err != nil {
		err = errs.NewStandardf(errs.Code(err), "获取会话数据失败: %v", err)
		return
	}
	if session == nil {
		err = errs.NewStandard(connect.CodeNotFound, errs.SessionNotFound)
		return
	}

	// 获取预设（包含提示项）
	if user.ActivePresetID > 0 {
		preset, err = c.presetRepo.GetByID(ctx, user.ActivePresetID, userID)
		if err != nil {
			err = errs.NewStandardf(errs.Code(err), "用户: %d 获取预设: %d 数据失败: %v",
				userID, user.ActivePresetID, err)
			return
		}
	}
	// TODO 有问题，没提示词项了
	if preset == nil {
		// 使用默认预设
		preset = &entity.Preset{
			Temperature:    constant.DefaultTemperature,
			TopP:           constant.DefaultTopP,
			MaxTokens:      constant.DefaultMaxTokens,
			CandidateCount: constant.DefaultCandidateCount,
		}
	}

	// 加载预设的提示项
	promptItems = make([]*entity.PromptItem, 0, len(preset.PromptItems))
	for _, item := range preset.PromptItems {
		if item.IsEnabled {
			promptItems = append(promptItems, &item)
		}
	}

	// 加载全局世界书
	globalWorldInfos, err = c.worldInfoRepo.ListGlobalWorldInfosWithEntries(ctx, userID)
	if err != nil {
		err = errs.NewStandardf(errs.Code(err), "用户: %d 获取全局世界书数据失败: %v", userID, err)
		return
	}

	// 加载角色卡关联的世界书
	if session.Character != nil && session.Character.WorldInfoID > 0 {
		charWorldInfo, err = c.worldInfoRepo.GetByIDWithEntries(ctx, session.Character.WorldInfoID, userID)
		if err != nil {
			err = errs.NewStandardf(errs.Code(err), "用户: %d 获取角色世界书数据失败: %v", userID, err)
			return
		}
	}

	// 加载正则规则
	regexRules, err = c.regexRepo.ListEnabledRules(ctx, preset.ID, session.CharacterID)
	if err != nil {
		err = errs.NewStandardf(errs.Code(err), "用户: %d 获取正则规则数据失败: %v", userID, err)
		return
	}
	persona, err = c.userRepo.GetPersonaByID(ctx, user.ActivePersonaID, userID)
	if err != nil {
		err = errs.NewStandardf(errs.Code(err), "用户: %d 获取人设数据失败: %v", userID, err)
	}
	if persona == nil {
		persona = &entity.Persona{
			UserID: userID,
			Name:   "凌溸",
		}
	}
	return
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
	persona *entity.Persona,
) *cache.SessionCache {
	sessionCache := &cache.SessionCache{
		SessionID:         int64(session.ID),
		UserID:            int64(session.UserID),
		CharacterID:       int64(session.CharacterID),
		Character:         session.Character,
		Persona:           persona,
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
		return model.NewMockModel() // 暂时返回 Gemini
	case pb.APIProvider_Claude:
		// TODO: 实现 Claude 模型
		return model.NewMockModel() // 暂时返回 Gemini
	default:
		return model.NewMockModel()
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
	persona *entity.Persona,
) ([]model.Message, error) {
	var messages []model.Message
	regexMap := make(map[string]*regexp2.Regexp, len(regexRules))
	// 根据预设填充对应的信息
	for _, promptItem := range promptItems {
		switch promptItem.Identifier {
		case pb.PromptItemIdentifier_Main, pb.PromptItemIdentifier_Jailbreak, pb.PromptItemIdentifier_Nsfw:
			message, err := createModelMsg(promptItem.Role, promptItem.Content, constant.Preset,
				regexMap, regexRules, 0)
			if err != nil {
				return nil, err
			}
			messages = append(messages, message)
		case pb.PromptItemIdentifier_ChatHistory:
			var msgs []model.Message
			// 收集需要按照深度插入的消息
			insertMsg := collectMsg(worldBook, promptItems, session)

			// 添加历史消息
			for i := len(session.Messages) - 1; i >= 0; i-- {
				history := session.Messages[i]
				if len(history.Swipes) == 0 {
					continue
				}
				depth := len(session.Messages) - i

				// 查看世界书和某些预设是否符合插入要求, 符合便插入
				if needWorldInfos, ok := insertMsg[depth]; ok {
					for _, needWorldInfo := range needWorldInfos {
						message, err := createModelMsg(
							needWorldInfo.Role, needWorldInfo.Content, needWorldInfo.Module,
							regexMap, regexRules, 0,
						)
						if err != nil {
							return nil, err
						}
						msgs = append(msgs, message)
					}
				}
				index := history.ActiveSwipeIndex
				if index >= len(history.Swipes) {
					index = len(history.Swipes) - 1
				}

				// 插入普通聊天记录
				var messageModule constant.ModuleType
				switch history.Role {
				case pb.Role_User:
					messageModule = constant.UserInput
				case pb.Role_Assistant:
					messageModule = constant.AIOutput
				}
				message, err := createModelMsg(history.Role, history.Swipes[index].Content, messageModule,
					regexMap, regexRules, depth)
				if err != nil {
					return nil, err
				}
				msgs = append(msgs, message)
			}

			// 插入那些大于当前消息深度的预设或世界书
			moreDeepMsgs := insertMaxDepthMsg(len(session.Messages), insertMsg)
			for _, msg := range moreDeepMsgs {
				message, err := createModelMsg(msg.Role, msg.Content, msg.Module, regexMap, regexRules, 0)
				if err != nil {
					return nil, err
				}
				msgs = append(msgs, message)
			}
			slices.Reverse(msgs)
			messages = append(messages, msgs...)
		case pb.PromptItemIdentifier_CharDescription:
			msg, err := createModelMsg(promptItem.Role, session.Character.Description, constant.Preset,
				regexMap, regexRules, 0)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		case pb.PromptItemIdentifier_DialogueExamples:
			for _, example := range session.Character.ExampleDialogue {
				msg, err := createModelMsg(promptItem.Role, example, constant.Preset,
					regexMap, regexRules, 0)
				if err != nil {
					return nil, err
				}
				messages = append(messages, msg)
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
					msg, err := createModelMsg(promptItem.Role, entry.Content, constant.WorldInfo,
						regexMap, regexRules, 0)
					if err != nil {
						return nil, err
					}
					messages = append(messages, msg)
				}
			}
		default:
			if promptItem.InjectionPosition == pb.InjectionPosition_Absolute {
				continue
			}
			msg, err := createModelMsg(promptItem.Role, promptItem.Content, constant.Preset,
				regexMap, regexRules, 0)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
	}
	messages = applyMacro(messages, session, persona)

	// 去除空内容
	n := 0
	for _, m := range messages {
		if m.Content != "" {
			messages[n] = m
			n++
		}
	}
	return messages[:n], nil
}

// PriorityMessage 继承model.Message并添加优先级字段
type PriorityMessage struct {
	model.Message
	Priority int
}

// collectMsg 收集需要按照深度插入的消息
func collectMsg(worldBooks []*entity.WorldInfo, promptItems []*entity.PromptItem,
	session *entity.ChatSession) map[int][]PriorityMessage {
	insertMsg := make(map[int][]PriorityMessage)
	// 处理世界书的深度插入信息
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
			insertMsg[entry.Depth] = append(insertMsg[entry.Depth], message)
		}
	}
	// 处理预设的深度插入信息
	for _, item := range promptItems {
		if item.InjectionPosition != pb.InjectionPosition_Absolute {
			continue
		}
		message := PriorityMessage{
			Message: model.Message{
				Role:    item.Role,
				Content: item.Content,
				Module:  constant.Preset,
			},
			// TODO Priority: item.SortOrder,
		}
		insertMsg[item.InjectionDepth] = append(insertMsg[item.InjectionDepth], message)
	}

	// 为每个深度的数组按优先级进行排序，优先级越大越靠前
	for depth := range insertMsg {
		sort.Slice(insertMsg[depth], func(i, j int) bool {
			return insertMsg[depth][i].Priority > insertMsg[depth][j].Priority
		})
	}

	return insertMsg
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

func insertMaxDepthMsg(maxDepth int, msgMap map[int][]PriorityMessage) []PriorityMessage {
	var depths []int
	for depth := range msgMap {
		depths = append(depths, depth)
	}
	msgs := make([]PriorityMessage, 0)
	sort.Ints(depths)
	for _, depth := range depths {
		if depth < maxDepth {
			continue
		}
		msgs = append(msgs, msgMap[depth]...)
	}
	return msgs
}

// createModelMsg 创建一个消息, 并进行正则匹配
func createModelMsg(role pb.Role, content string, module constant.ModuleType,
	regexsMap map[string]*regexp2.Regexp, regexs []*entity.RegexRule, depth int) (model.Message, error) {
	msg := model.Message{
		Role:    role,
		Content: content,
		Module:  module,
	}
	// 进行正则匹配
	replaceContent, err := applyRegex(regexsMap, regexs, msg, depth)
	if err != nil {
		log.Errorf("applyRegex error: %+v", err)
		return msg, err
	}
	msg.Content = replaceContent
	return msg, nil
}

func applyRegex(regexsMap map[string]*regexp2.Regexp, regexs []*entity.RegexRule,
	message model.Message, depth int) (string, error) {
	content := message.Content
	for _, regex := range regexs {
		if !regex.IsEnabled {
			log.Debugf("正则: %s规则未启用", regex.FindPattern)
			continue
		}
		if message.Module != constant.Preset && !regex.AffectFlagsPrompt ||
			message.Module != constant.UserInput && !regex.AffectFlagsUserInput ||
			message.Module != constant.AIOutput && !regex.AffectFlagsAIOutput ||
			message.Module != constant.WorldInfo && !regex.AffectFlagsWorldInfo {
			continue
		}
		// 检查消息深度是否在允许范围内, 且只有聊天记录才检查深度
		if message.Module != constant.Preset && message.Module != constant.WorldInfo &&
			depth < regex.MinDepth && depth > regex.MaxDepth {
			log.Debugf("正则: %s规则消息深度超出范围: %d, min: %d, max: %d",
				regex.FindPattern, depth, regex.MinDepth, regex.MaxDepth)
			continue
		}

		re, ok := regexsMap[regex.FindPattern]
		if regexsMap != nil && !ok {
			log.Debugf("正则: %s规则未编译，正在编译...", regex.FindPattern)
			compileRe, err := convert.ConvertToRegexp2(regex.FindPattern)
			if err != nil {
				return message.Content, errs.NewStandardf(connect.CodeInternal, "正则: %s转换失败: %v",
					regex.FindPattern, err)
			}
			regexsMap[regex.FindPattern] = compileRe.Engine
			re = compileRe.Engine
		}
		replaceContent, err := re.Replace(message.Content, regex.ReplacePattern, -1, -1)
		if err != nil {
			return message.Content, errs.NewStandardf(connect.CodeInternal, "正则: %s替换内容失败: %v",
				regex.FindPattern, err)
		}
		if content != replaceContent {
			log.Debugf("应用正则表达式: %s, 替换为文本: %s, 替换前: %s,\n\n\n 替换后: %s\n",
				regex.FindPattern, regex.ReplacePattern, content, replaceContent)
		}
		content = replaceContent
	}
	return content, nil
}

// applyMacro 对消息列表应用宏替换
// 支持的宏:
//   - {{char}}, <char>, <bot>, <BOT>: 角色名称
//   - {{newline}}: 换行符
//   - {{noop}}: 空字符串
//   - {{time}}: 当前时间 (HH:mm:ss)
//   - {{date}}: 当前日期 (YYYY-MM-DD)
//   - {{weekday}}: 星期几
//   - {{isotime}}: ISO 时间 (HH:mm:ss)
//   - {{isodate}}: ISO 日期 (YYYY-MM-DD)
//   - {{description}}: 角色描述
//   - {{mesExamples}}: 示例对话
//   - {{lastMessage}}: 最后一条消息
//   - {{lastUserMessage}}: 最后一条用户消息
//   - {{lastCharMessage}}: 最后一条角色消息
//   - {{lastMessageId}}: 最后一条消息ID
func applyMacro(messages []model.Message, session *entity.ChatSession, persona *entity.Persona) []model.Message {
	if session == nil || session.Character == nil {
		return messages
	}

	// 获取时间相关数据
	now := time.Now()
	charName := session.Character.Name
	userName := persona.Name
	description := session.Character.Description

	// 处理示例对话
	mesExamples := ""
	if len(session.Character.ExampleDialogue) > 0 {
		mesExamples = strings.Join(session.Character.ExampleDialogue, "\n")
	}

	// 获取消息相关数据
	var lastMessage, lastUserMessage, lastCharMessage string
	var lastMessageId int32

	if len(session.Messages) > 0 {
		lastMsg := session.Messages[len(session.Messages)-1]
		lastMessageId = int32(lastMsg.ID)
		if len(lastMsg.Swipes) > 0 && int(lastMsg.ActiveSwipeIndex) < len(lastMsg.Swipes) {
			lastMessage = lastMsg.Swipes[lastMsg.ActiveSwipeIndex].Content
		}

		// 反向遍历查找最后一条用户消息和角色消息
		for i := len(session.Messages) - 1; i >= 0; i-- {
			msg := session.Messages[i]
			if len(msg.Swipes) == 0 || int(msg.ActiveSwipeIndex) >= len(msg.Swipes) {
				continue
			}
			content := msg.Swipes[msg.ActiveSwipeIndex].Content

			if msg.Role == pb.Role_User && lastUserMessage == "" {
				lastUserMessage = content
			}
			if msg.Role == pb.Role_Assistant && lastCharMessage == "" {
				lastCharMessage = content
			}

			if lastUserMessage != "" && lastCharMessage != "" {
				break
			}
		}
	}

	// 使用 strings.NewReplacer 进行高效替换
	replacer := strings.NewReplacer(
		// 角色名称
		constant.Char1, charName,
		constant.Char2, charName,
		constant.Char3, charName,
		constant.Char4, charName,
		// 用户名称
		constant.User1, userName,
		constant.User2, userName,
		// 基础宏
		constant.Newline, "\n",
		constant.Noop, "",
		// 时间日期宏
		constant.Time, now.Format("15:04:05"),
		constant.Date, now.Format("2006-01-02"),
		constant.Weekday, now.Weekday().String(),
		constant.IsoTime, now.Format("15:04:05"),
		constant.IsoDate, now.Format("2006-01-02"),
		// 角色数据宏
		constant.Description, description,
		constant.MesExamples, mesExamples,
		// 消息相关宏
		constant.LastMessage1, lastMessage,
		constant.LastMessage2, lastMessage,
		constant.LastMessage3, lastMessage,
		constant.LastUserMessage1, lastUserMessage,
		constant.LastUserMessage2, lastUserMessage,
		constant.LastUserMessage3, lastUserMessage,
		constant.LastCharMessage1, lastCharMessage,
		constant.LastCharMessage2, lastCharMessage,
		constant.LastCharMessage3, lastCharMessage,
		constant.LastMessageId, fmt.Sprintf("%d", lastMessageId),
	)

	// 应用替换
	for i := range messages {
		messages[i].Content = replacer.Replace(messages[i].Content)
	}

	return messages
}

func (c *chatImpl) DeleteSwipe(ctx context.Context, req *pb.DeleteSwipeRequest) (*pb.DeleteSwipeResponse, error) {

	if err := c.chatRepo.DeleteSwipe(ctx, int(req.MessageId), int(req.SwipeId)); err != nil {
		return nil, err
	}
	return &pb.DeleteSwipeResponse{}, nil
}
