package main

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func (repo *Repo) UpdateBalances(req AccountEvent) error {
	tx := repo.DB.Begin()

	if err := tx.Model(&AccountInfo{}).
		Where("account_no = ?", req.AccountNo).
		Update("account_balance", gorm.Expr("account_balance + ?", req.Amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func GenerateString() string {
	// Create a new Node with a Node number of 1
	node, _ := snowflake.NewNode(1)

	return node.Generate().String()
}
