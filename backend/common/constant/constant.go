package constant

const (
	// DefaultPageNum 默认页码
	DefaultPageNum = 1
	// DefaultPageSize 默认页大小
	DefaultPageSize = 20
	// MaxPageSize 最大页码
	MaxPageSize = 100
)

// CtxKey 上下文键
type CtxKey string

const (
	// TransactionKey 获取ctx中的事务键
	TransactionKey CtxKey = "transaction"
	// UserIDKey 获取ctx中的用户ID键
	UserIDKey CtxKey = "userId"
)

const SortOrderInterval = 100

// ModuleType 模块类型
type ModuleType int

const (
	// UserInput 用户输入
	UserInput ModuleType = iota
	// AIOutput AI输出
	AIOutput
	// Preset 预设
	Preset
	// WorldInfo 世界信息
	WorldInfo
)

// NormalizePagination 处理分页参数，返回规范化后的 page 和 pageSize
func NormalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = DefaultPageNum
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return page, pageSize
}
