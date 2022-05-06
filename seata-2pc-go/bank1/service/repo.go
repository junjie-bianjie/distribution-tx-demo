package service

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func (repo *Repo) Transfer(ctx context.Context, req *AccountEvent) (interface{}, error) {
	tx := repo.DB.WithContext(ctx).Begin(&sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})

	err := tx.Model(&AccountInfo{}).
		Where("account_no = ?", 1). // 1是张三的no
		Update("account_balance", gorm.Expr("account_balance - ?", req.Amount)).
		Commit().
		Error

	return nil, err
}
