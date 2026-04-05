package database

import (
	"context"
	"time"

	"github.com/ling/muse/common/constant"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	"gorm.io/gorm"
)

// GetDB 获取数据库实例
func GetDB(ctx context.Context) *gorm.DB {
	if db := ctx.Value(constant.TransactionKey); db != nil {
		return db.(*gorm.DB)
	}
	return config.GetDB()
}

// CharacterRepository 角色卡数据仓库接口
type CharacterRepository interface {
	// Create 创建角色
	Create(ctx context.Context, character *entity.Character) error
	// GetByID 根据ID获取角色
	GetByID(ctx context.Context, id, userID int) (*entity.Character, error)
	// List 获取角色列表
	List(ctx context.Context, userID, page, pageSize int) ([]*entity.Character, int64, error)
	// Update 更新角色
	Update(ctx context.Context, character *entity.Character) error
	// Delete 删除角色
	Delete(ctx context.Context, id, userID int) error
	// GetVersion 获取角色的版本号
	GetVersion(ctx context.Context, id, userID int) (int, error)
}

// ChatRepository 聊天会话数据仓库接口
type ChatRepository interface {
	// CreateSession 创建聊天会话
	CreateSession(ctx context.Context, session *entity.ChatSession) error
	// GetSessionByID 根据ID获取聊天会话
	GetSessionByID(ctx context.Context, id int, userID int) (*entity.ChatSession, error)
	// ListSessions 获取会话列表
	ListSessions(ctx context.Context, userID, characterID, page, pageSize int) ([]*entity.ChatSession, int64, error)
	// UpdateSession 更新聊天会话
	UpdateSession(ctx context.Context, sessionID, userID, version int, sessionName string) (int, error)
	// DeleteSession 删除聊天会话
	DeleteSession(ctx context.Context, id, userID int) error
	// GetMessageByID 根据ID获取消息
	GetMessageByID(ctx context.Context, id int) (*entity.Message, error)
	// DeleteMessage 删除消息
	DeleteMessage(ctx context.Context, id int) error
	// UpdateMessage 更新消息的activeSwipeIndex
	UpdateMessage(ctx context.Context, message *entity.Message) error
	// GetSessionWithMessages 获取会话及其所有消息和角色卡
	GetSessionWithMessages(ctx context.Context, id int, userID int) (*entity.ChatSession, error)
	// CreateMessage 创建消息
	CreateMessage(ctx context.Context, message *entity.Message) error
	// CreateMessageSwipe 创建消息swipe
	CreateMessageSwipe(ctx context.Context, swipe *entity.MessageSwipe) error
	// GetMaxMessageSortOrder 获取会话中消息的最大排序号
	GetMaxMessageSortOrder(ctx context.Context, sessionID int) (int, error)
	// UpdateMessageSwipe 更新消息swipe内容
	UpdateMessageSwipe(ctx context.Context, swipeID int, content string) error
	// GetCharLatestSessionWithMsg 获取角色最新的会话
	GetCharLatestSessionWithMsg(ctx context.Context, charID, userID int) (*entity.ChatSession, error)
	// UpdateSessionTime 更新会话的更新时间
	UpdateSessionTime(ctx context.Context, sessionID int, updateTime time.Time) error
	// SwitchSwipe 切换消息的swipe
	SwitchSwipe(ctx context.Context, messageID, index int) error
}

// PresetRepository 预设数据仓库接口
type PresetRepository interface {
	// Create 创建预设
	Create(ctx context.Context, preset *entity.Preset) error
	// GetByID 根据ID获取预设（包含关联的PromptItems）
	GetByID(ctx context.Context, id int, userID int) (*entity.Preset, error)
	// List 获取预设列表
	List(ctx context.Context, userID int, page int, pageSize int) ([]*entity.Preset, []int, int64, error)
	// Update 更新预设
	Update(ctx context.Context, preset *entity.Preset) error
	// Delete 删除预设
	Delete(ctx context.Context, id int, userID int) error
	// CreatePromptItem 创建提示项
	CreatePromptItem(ctx context.Context, item *entity.PromptItem) error
	// GetPromptItemByID 根据ID获取提示项
	GetPromptItemByID(ctx context.Context, id int) (*entity.PromptItem, error)
	// ListPromptItems 获取预设的提示项列表
	ListPromptItems(ctx context.Context, presetID int) ([]*entity.PromptItem, error)
	// UpdatePromptItem 更新提示项
	UpdatePromptItem(ctx context.Context, item *entity.PromptItem) error
	// DeletePromptItem 删除提示项
	DeletePromptItem(ctx context.Context, id int) error
	// UpdatePromptItemsOrder 更新提示项排序
	UpdatePromptItemsOrder(ctx context.Context, presetID int, userID int, sourceID int, desID int, operation pb.SortOperation) error
	// GetVersion 获取预设的版本号
	GetVersion(ctx context.Context, id int, userID int) (int, error)
	// BatchCreatePromptItem 批量创建提示项目
	BatchCreatePromptItem(ctx context.Context, items []*entity.PromptItem) error
	// BatchUpdatePromptItemOrder 批量更新提示项排序
	BatchUpdatePromptItemOrder(ctx context.Context, items []*entity.PromptItem) error
	// SetActivePreset 设置用户的启用预设
	SetActivePreset(ctx context.Context, userID, presetID int) error
}

// RegexRuleRepository 正则规则数据仓库接口
type RegexRuleRepository interface {
	// Create 创建正则规则
	Create(ctx context.Context, rule *entity.RegexRule) error
	// GetByID 根据ID获取正则规则
	GetByID(ctx context.Context, id int) (*entity.RegexRule, error)
	// List 获取预设的正则规则列表
	List(ctx context.Context, presetID int) ([]*entity.RegexRule, error)
	// Update 更新正则规则
	Update(ctx context.Context, rule *entity.RegexRule) error
	// Delete 删除正则规则
	Delete(ctx context.Context, id int) error
	// UpdateRulesOrder 更新正则规则排序
	UpdateRulesOrder(ctx context.Context, presetID int, ruleOrders map[int]int) error
	// BatchCreate 批量创建正则规则
	BatchCreate(ctx context.Context, rules []*entity.RegexRule) error
	// GetMaxSortOrder 获取预设下正则规则的最大排序号
	GetMaxSortOrder(ctx context.Context, presetID int) (int, error)
	// ListEnabledRules 获取启用的正则规则列表
	ListEnabledRules(ctx context.Context, presetID int, characterID int) ([]*entity.RegexRule, error)
	// GetEnabledRuleVersions 批量获取启用的正则规则的版本号
	GetEnabledRuleVersions(ctx context.Context, presetID int, characterID int) (map[int64]int, error)
}

// WorldInfoRepository 世界书数据仓库接口
type WorldInfoRepository interface {
	// Create 创建世界书
	Create(ctx context.Context, worldInfo *entity.WorldInfo) error
	// GetByID 根据ID获取世界书
	GetByID(ctx context.Context, id int, userID int) (*entity.WorldInfo, error)
	// List 获取世界书列表
	List(ctx context.Context, userID, page, pageSize int) ([]*entity.WorldInfo, int64, error)
	// Update 更新世界书
	Update(ctx context.Context, worldInfo *entity.WorldInfo) error
	// Delete 删除世界书
	Delete(ctx context.Context, id, userID int) error
	// CreateEntry 创建世界书条目
	CreateEntry(ctx context.Context, entry *entity.WorldInfoEntry) error
	// BatchCreateEntries 批量创建世界书条目
	BatchCreateEntries(ctx context.Context, entries []*entity.WorldInfoEntry) error
	// GetEntryByID 根据ID获取世界书条目
	GetEntryByID(ctx context.Context, id int) (*entity.WorldInfoEntry, error)
	// ListEntries 获取世界书的条目列表
	ListEntries(ctx context.Context, worldInfoID int) ([]*entity.WorldInfoEntry, error)
	// UpdateEntry 更新世界书条目
	UpdateEntry(ctx context.Context, entry *entity.WorldInfoEntry) error
	// DeleteEntry 删除世界书条目
	DeleteEntry(ctx context.Context, id int) error
	// UpdateEntriesOrder 更新世界书条目排序
	UpdateEntriesOrder(ctx context.Context, worldInfoID int, entryOrders map[int]int) error
	// ListGlobalWorldInfosWithEntries 获取用户所有全局世界书及其启用的条目
	ListGlobalWorldInfosWithEntries(ctx context.Context, userID int) ([]*entity.WorldInfo, error)
	// GetByIDWithEntries 根据ID获取世界书及其所有启用的条目
	GetByIDWithEntries(ctx context.Context, id int, userID int) (*entity.WorldInfo, error)
	// GetVersions 批量获取世界书的版本号
	GetVersions(ctx context.Context, ids []int, userID int) (map[int64]int, error)
	// GetGlobalWorldInfoVersions 获取用户所有全局世界书的版本号
	GetGlobalWorldInfoVersions(ctx context.Context, userID int) (map[int64]int, error)
}

// UserRepository 用户数据仓库接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error
	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id int) (*entity.User, error)
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	// Update 更新用户
	Update(ctx context.Context, user *entity.User) error
	// UpdateActivePersonaID 更新用户的活跃人设ID
	UpdateActivePersonaID(ctx context.Context, userID int, personaID int) error
	// UpdateActivePresetID 更新用户的活跃预设ID
	UpdateActivePresetID(ctx context.Context, userID int, presetID int) error

	// CreatePersona 创建人设
	CreatePersona(ctx context.Context, persona *entity.Persona) error
	// GetPersonaByID 根据ID获取人设
	GetPersonaByID(ctx context.Context, id int, userID int) (*entity.Persona, error)
	// ListPersonas 获取用户的人设列表
	ListPersonas(ctx context.Context, userID int) ([]*entity.Persona, error)
	// UpdatePersona 更新人设
	UpdatePersona(ctx context.Context, persona *entity.Persona) error
	// DeletePersona 删除人设
	DeletePersona(ctx context.Context, id int, userID int) error

	// CreateUserSetting 创建用户设置
	CreateUserSetting(ctx context.Context, setting *entity.UserSetting) error
	// GetUserSettingByUserID 根据用户ID获取用户设置
	GetUserSettingByUserID(ctx context.Context, userID int) (*entity.UserSetting, error)
	// UpdateUserSetting 更新用户设置
	UpdateUserSetting(ctx context.Context, setting *entity.UserSetting) error

	// CreateAPIConfig 创建API配置
	CreateAPIConfig(ctx context.Context, apiConfig *entity.APIConfig) error
	// GetAPIConfigByID 根据ID获取API配置
	GetAPIConfigByID(ctx context.Context, id int, userID int) (*entity.APIConfig, error)
	// ListAPIConfigs 获取用户的API配置列表
	ListAPIConfigs(ctx context.Context, userID int) ([]*entity.APIConfig, error)
	// UpdateAPIConfig 更新API配置
	UpdateAPIConfig(ctx context.Context, apiConfig *entity.APIConfig) error
	// DeleteAPIConfig 删除API配置
	DeleteAPIConfig(ctx context.Context, id int, userID int) error
	// DeactivateAllAPIConfigs 停用用户的所有API配置
	DeactivateAllAPIConfigs(ctx context.Context, userID int) error
	// ActivateAPIConfig 激活API配置
	ActivateAPIConfig(ctx context.Context, id int, userID int) error
	// GetActiveAPIConfig 获取用户当前活跃的API配置
	GetActiveAPIConfig(ctx context.Context, userID int) (*entity.APIConfig, error)
}
