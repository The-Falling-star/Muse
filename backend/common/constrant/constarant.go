package constrant

const (
	DefaultPageNum  = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
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
