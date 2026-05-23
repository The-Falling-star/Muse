// Package cache 提供会话缓存管理功能
// 当前使用内存 LRU 缓存实现，后续可迁移到 Redis
package cache

import (
	"container/list"
	"sync"
	"time"

	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
)

// 全局缓存管理器实例
var globalCacheManager *SessionCacheManager
var cacheOnce sync.Once

// GetSessionCacheManager 获取全局会话缓存管理器单例
func GetSessionCacheManager() *SessionCacheManager {
	cacheOnce.Do(func() {
		globalCacheManager = NewSessionCacheManager(100, 30*time.Minute)
	})
	return globalCacheManager
}

// InvalidateCacheByCharacter 使指定角色相关的所有缓存失效
func InvalidateCacheByCharacter(characterID int64) {
	GetSessionCacheManager().RemoveByCharacter(characterID)
}

// InvalidateCacheByPreset 使使用指定预设的所有缓存失效
func InvalidateCacheByPreset(presetID int64) {
	GetSessionCacheManager().RemoveByPreset(presetID)
}

// InvalidateCacheByWorldInfo 使使用指定世界书的所有缓存失效
func InvalidateCacheByWorldInfo(worldInfoID int64) {
	GetSessionCacheManager().RemoveByWorldInfo(worldInfoID)
}

// InvalidateCacheByRegexRule 使使用指定正则规则的所有缓存失效
func InvalidateCacheByRegexRule(regexID int64) {
	GetSessionCacheManager().RemoveByRegexRule(regexID)
}

// InvalidateCacheByUser 使指定用户的所有缓存失效
func InvalidateCacheByUser(userID int64) {
	GetSessionCacheManager().RemoveByUser(userID)
}

// SessionCache 会话缓存数据
// 缓存会话相关的所有数据，避免每次请求都查询数据库
type SessionCache struct {
	SessionID   int64             // 会话ID
	UserID      int64             // 用户ID
	CharacterID int64             // 角色ID
	Character   *entity.Character // 角色卡数据
	Persona     *entity.Persona   // 人设数据

	// 世界书相关
	WorldInfos       []*entity.WorldInfo      // 世界书列表
	WorldInfoEntries []*entity.WorldInfoEntry // 启用的世界书条目
	WorldInfoByDepth map[int][]string         // 按深度分组的世界书内容（仅包含 Constant 条目）

	// 预设相关
	Preset      *entity.Preset       // 预设数据
	PromptItems []*entity.PromptItem // 预设提示项

	// 正则规则
	RegexRules []*entity.RegexRule // 正则规则列表

	// 历史消息（增量更新）
	Messages []*entity.Message // 历史消息列表

	// 版本信息（用于校验缓存是否过期）
	CharacterVersion  int           // 角色卡版本
	PresetVersion     int           // 预设版本
	WorldInfoVersions map[int64]int // 世界书版本 world_info_id -> version
	RegexRuleVersions map[int64]int // 正则规则版本 regex_id -> version

	// 元数据
	LastAccessTime time.Time // 最后访问时间
	CreatedAt      time.Time // 创建时间
}

// VersionInfo 版本信息，用于校验缓存是否过期
type VersionInfo struct {
	CharacterVersion  int           // 角色卡版本
	PresetVersion     int           // 预设版本
	WorldInfoVersions map[int64]int // 世界书版本
	RegexRuleVersions map[int64]int // 正则规则版本
}

// cacheKey 缓存键
type cacheKey struct {
	UserID    int64
	SessionID int64
}

// lruEntry LRU 链表节点数据
type lruEntry struct {
	key   cacheKey
	cache *SessionCache
}

// SessionCacheManager 会话缓存管理器
// 使用 LRU 策略管理缓存，限制最大缓存数量
type SessionCacheManager struct {
	mu      sync.RWMutex
	cache   map[cacheKey]*list.Element // 缓存映射
	lruList *list.List                 // LRU 链表，最近使用的在前
	maxSize int                        // 最大缓存数量
	ttl     time.Duration              // 缓存过期时间
}

// NewSessionCacheManager 创建会话缓存管理器
func NewSessionCacheManager(maxSize int, ttl time.Duration) *SessionCacheManager {
	return &SessionCacheManager{
		cache:   make(map[cacheKey]*list.Element),
		lruList: list.New(),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get 获取缓存
// 返回缓存数据和是否命中
func (m *SessionCacheManager) Get(userID, sessionID int64) (*SessionCache, bool) {
	if !config.Get().Chat.EnableCache {
		return nil, false
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	key := cacheKey{UserID: userID, SessionID: sessionID}
	if elem, ok := m.cache[key]; ok {
		entry := elem.Value.(*lruEntry)

		// 检查是否过期
		if time.Since(entry.cache.LastAccessTime) > m.ttl {
			m.removeElement(elem)
			return nil, false
		}

		// 更新访问时间并移到链表头部
		entry.cache.LastAccessTime = time.Now()
		m.lruList.MoveToFront(elem)
		return entry.cache, true
	}
	return nil, false
}

// Set 设置缓存
func (m *SessionCacheManager) Set(userID, sessionID int64, cache *SessionCache) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := cacheKey{UserID: userID, SessionID: sessionID}
	now := time.Now()
	cache.LastAccessTime = now
	if cache.CreatedAt.IsZero() {
		cache.CreatedAt = now
	}

	// 如果已存在，更新并移到头部
	if elem, ok := m.cache[key]; ok {
		entry := elem.Value.(*lruEntry)
		entry.cache = cache
		m.lruList.MoveToFront(elem)
		return
	}

	// 添加新缓存
	entry := &lruEntry{key: key, cache: cache}
	elem := m.lruList.PushFront(entry)
	m.cache[key] = elem

	// 超出容量时淘汰最久未使用的
	for m.lruList.Len() > m.maxSize {
		m.removeOldest()
	}
}

// Remove 删除指定缓存
func (m *SessionCacheManager) Remove(userID, sessionID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := cacheKey{UserID: userID, SessionID: sessionID}
	if elem, ok := m.cache[key]; ok {
		m.removeElement(elem)
	}
}

// RemoveByCharacter 删除指定角色相关的所有缓存
func (m *SessionCacheManager) RemoveByCharacter(characterID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var toRemove []*list.Element
	for elem := m.lruList.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*lruEntry)
		if entry.cache.CharacterID == characterID {
			toRemove = append(toRemove, elem)
		}
	}

	for _, elem := range toRemove {
		m.removeElement(elem)
	}
}

// RemoveByPreset 删除使用指定预设的所有缓存
func (m *SessionCacheManager) RemoveByPreset(presetID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var toRemove []*list.Element
	for elem := m.lruList.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*lruEntry)
		if entry.cache.Preset != nil && int64(entry.cache.Preset.ID) == presetID {
			toRemove = append(toRemove, elem)
		}
	}

	for _, elem := range toRemove {
		m.removeElement(elem)
	}
}

// RemoveByWorldInfo 删除使用指定世界书的所有缓存
func (m *SessionCacheManager) RemoveByWorldInfo(worldInfoID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var toRemove []*list.Element
	for elem := m.lruList.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*lruEntry)
		if _, ok := entry.cache.WorldInfoVersions[worldInfoID]; ok {
			toRemove = append(toRemove, elem)
		}
	}

	for _, elem := range toRemove {
		m.removeElement(elem)
	}
}

// RemoveByRegexRule 删除使用指定正则规则的所有缓存
func (m *SessionCacheManager) RemoveByRegexRule(regexID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var toRemove []*list.Element
	for elem := m.lruList.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*lruEntry)
		if _, ok := entry.cache.RegexRuleVersions[regexID]; ok {
			toRemove = append(toRemove, elem)
		}
	}

	for _, elem := range toRemove {
		m.removeElement(elem)
	}
}

// RemoveByUser 删除指定用户的所有缓存
func (m *SessionCacheManager) RemoveByUser(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var toRemove []*list.Element
	for elem := m.lruList.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*lruEntry)
		if entry.cache.UserID == userID {
			toRemove = append(toRemove, elem)
		}
	}

	for _, elem := range toRemove {
		m.removeElement(elem)
	}
}

// Clear 清空所有缓存
func (m *SessionCacheManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache = make(map[cacheKey]*list.Element)
	m.lruList.Init()
}

// Size 返回当前缓存数量
func (m *SessionCacheManager) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lruList.Len()
}

// removeElement 删除链表节点
func (m *SessionCacheManager) removeElement(elem *list.Element) {
	entry := elem.Value.(*lruEntry)
	delete(m.cache, entry.key)
	m.lruList.Remove(elem)
}

// removeOldest 删除最久未使用的缓存
func (m *SessionCacheManager) removeOldest() {
	elem := m.lruList.Back()
	if elem != nil {
		m.removeElement(elem)
	}
}

// CleanExpired 清理过期缓存
// 可以定期调用此方法清理过期的缓存
func (m *SessionCacheManager) CleanExpired() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	var toRemove []*list.Element
	now := time.Now()

	for elem := m.lruList.Back(); elem != nil; elem = elem.Prev() {
		entry := elem.Value.(*lruEntry)
		if now.Sub(entry.cache.LastAccessTime) > m.ttl {
			toRemove = append(toRemove, elem)
		} else {
			// 因为是 LRU 排序，后面的都是更新的，可以提前结束
			break
		}
	}

	for _, elem := range toRemove {
		m.removeElement(elem)
	}

	return len(toRemove)
}

// IsVersionValid 检查缓存版本是否有效
func (c *SessionCache) IsVersionValid(versions *VersionInfo) bool {
	// 检查角色卡版本
	if c.CharacterVersion != versions.CharacterVersion {
		return false
	}

	// 检查预设版本
	if c.PresetVersion != versions.PresetVersion {
		return false
	}

	// 检查世界书版本
	if len(c.WorldInfoVersions) != len(versions.WorldInfoVersions) {
		return false
	}
	for id, v := range c.WorldInfoVersions {
		if versions.WorldInfoVersions[id] != v {
			return false
		}
	}

	// 检查正则规则版本
	if len(c.RegexRuleVersions) != len(versions.RegexRuleVersions) {
		return false
	}
	for id, v := range c.RegexRuleVersions {
		if versions.RegexRuleVersions[id] != v {
			return false
		}
	}

	return true
}

// AddMessage 向缓存中添加消息
func (c *SessionCache) AddMessage(msg *entity.Message) {
	c.Messages = append(c.Messages, msg)
}

// GetLastMessage 获取最后一条消息
func (c *SessionCache) GetLastMessage() *entity.Message {
	if len(c.Messages) == 0 {
		return nil
	}
	return c.Messages[len(c.Messages)-1]
}
