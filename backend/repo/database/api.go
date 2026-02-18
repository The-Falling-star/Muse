package database

import (
	"context"

	"github.com/ling/muse/common/constant"
	"github.com/ling/muse/config"
	"gorm.io/gorm"
)

// GetDB 获取数据库实例
func GetDB(ctx context.Context) *gorm.DB {
	if db := ctx.Value(constant.TransactionKey); db != nil {
		return db.(*gorm.DB)
	}
	return config.GetDB()
}
