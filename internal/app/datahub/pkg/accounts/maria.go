package accounts

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

	"github.com/elgris/sqrl"
	"go.uber.org/zap"
)

// SelectAccount 获取账号
func SelectAccount(id string, qq uint64) (data []*Tab, err error) {

	sqr := sqrl.Select("*").From("accounts")
	if id != "" {
		sqr = sqr.Where("id=?", id)
	} else {
		sqr = sqr.Where("qq=?", qq)
	}

	sql, args, err := sqr.Limit(1).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	err = maria.Select(&data, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}

	return
}
