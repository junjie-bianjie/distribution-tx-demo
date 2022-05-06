package main

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) CreateUserAndLog(user User, msg MsgLog) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		msg.UserId = user.Id
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}

type MsgLogRepo struct {
	db *gorm.DB
}

func NewMsgLogRepo(db *gorm.DB) *MsgLogRepo {
	return &MsgLogRepo{db: db}
}

func (repo *MsgLogRepo) FindNeedExec() ([]MsgLog, error) {
	var res []MsgLog
	err := repo.db.Where(MsgLog{}).Where("status = ?", Undo).
		Find(&res).Error
	return res, err
}

func (repo *MsgLogRepo) BatchUpdateResolve(uuids []string) error {
	err := repo.db.Model(MsgLog{}).Where("uuid in ?", uuids).Update("status", Do).Error
	return err
}

func (repo *MsgLogRepo) Exists(uuids string) (bool, error) {
	var count int64
	err := repo.db.Model(MsgLog{}).Where("uuid = ?", uuids).Count(&count).Error
	return count > 0, err
}

func GenerateString() string {
	// Create a new Node with a Node number of 1
	node, _ := snowflake.NewNode(1)

	return node.Generate().String()
}

type IntegralRepo struct {
	db *gorm.DB
}

func NewIntegralRepo(db *gorm.DB) *IntegralRepo {
	return &IntegralRepo{db: db}
}

func (repo *IntegralRepo) CreateIntegralAndLog(integral Integral, msg MsgLog) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&integral).Error; err != nil {
			return err
		}

		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}
