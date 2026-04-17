package cache

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/groupcache/lru"
	"github.com/ling/muse/common/crypto"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	log "github.com/sirupsen/logrus"
)

// APIKeyCache APIKey的缓存
type APIKeyCache struct {
	mu     sync.RWMutex
	cache  *lru.Cache
	aesGcm *crypto.AESGCM
}

// APIKeyCacheEntry APIKey缓存项
type APIKeyCacheEntry struct {
	APIKeys  []*entity.APIConfig
	index    atomic.Int32
	createAt int64
}

// NewAPIKeyCache 新建一个APIKey LRU缓存
func NewAPIKeyCache() APIKeyCache {
	aesGcm, err := crypto.NewAESGCMFromKey(config.Get().APIEncrypt.EncryptionKey)
	if err != nil {
		log.Fatalf("初始化API Key GCM失败: %v", err)
	}
	return APIKeyCache{
		cache:  lru.New(config.Get().Cache.LRUMaxSize),
		aesGcm: aesGcm,
	}
}

// Get 获取缓存元素
func (a *APIKeyCache) Get(userID int, provider pb.APIProvider) *entity.APIConfig {
	a.mu.RLock()

	key := fmt.Sprintf("%d:%d", userID, provider)
	val, exist := a.cache.Get(key)
	if !exist {
		a.mu.RUnlock()
		return nil
	}

	entry, _ := val.(*APIKeyCacheEntry)
	// 检查过期
	if time.Now().Unix()-entry.createAt > int64(config.Get().Cache.LRUTTL) {
		a.mu.RUnlock()
		a.mu.Lock()
		defer a.mu.Unlock()
		a.cache.Remove(key)
		return nil
	}
	// 保证读锁解锁
	defer a.mu.RUnlock()

	if len(entry.APIKeys) == 0 {
		return nil
	}

	// 轮询返回一个 APIKey
	idx := entry.index.Add(1)
	return entry.APIKeys[int(idx)%len(entry.APIKeys)]
}

// Set 设置缓存数据
func (a *APIKeyCache) Set(userID int, provider pb.APIProvider, APIKeys []*entity.APIConfig) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cacheEntry := &APIKeyCacheEntry{
		APIKeys:  APIKeys,
		createAt: time.Now().Unix(),
	}

	// 解码API Key
	for _, apiKey := range APIKeys {
		decryptKey, err := a.aesGcm.Decrypt(apiKey.APIKey)
		if err != nil {
			log.Warnf("解码API Key: %s 失败: %v", apiKey.APIKey, err)
			continue
		}
		apiKey.APIKey = decryptKey
	}
	a.cache.Add(fmt.Sprintf("%d:%d", userID, provider), cacheEntry)
}
