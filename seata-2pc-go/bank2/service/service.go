package service

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func (repo *Repo) AddBalances(ctx context.Context, req AccountEvent) error {
	tx := repo.DB.WithContext(ctx).Begin(&sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})

	if err := tx.Model(&AccountInfo{}).
		Where("account_no = ?", req.AccountNo). // 这里会拿到李四的no
		Update("account_balance", gorm.Expr("account_balance + ?", req.Amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if true {
		tx.Rollback()
		return errors.New("ahhaha")
	}
	err := tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}
