package constrant

const (
	// DefaultPageNum 默认页码
	DefaultPageNum = 1
	// DefaultPageSize 默认页大小
	DefaultPageSize = 20
	// MaxPageSize 最大页码
	MaxPageSize = 100
)

// TransactionKey 获取ctx中的事务键
const TransactionKey = "transaction"

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
