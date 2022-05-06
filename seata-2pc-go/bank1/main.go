package main

import (
	"github.com/opentrx/mysql/v2"
	"github.com/opentrx/seata-golang/v2/pkg/client/tm"
	"os"
	"time"

	svc "distribution-tx-demo/seata-2pc-go/bank1/service"
	dialector "distribution-tx-demo/seata-2pc-go/dialector/mysql"
	"github.com/gin-gonic/gin"
	"github.com/opentrx/seata-golang/v2/pkg/client"
	"github.com/opentrx/seata-golang/v2/pkg/client/config"
	"github.com/opentrx/seata-golang/v2/pkg/client/rm"
	"github.com/opentrx/seata-golang/v2/pkg/util/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	r := gin.Default()

	configPath := os.Getenv("ConfigPath")
	conf := config.InitConfiguration(configPath)

	log.Init(conf.Log.LogPath, conf.Log.LogLevel)
	client.Init(conf)

	rm.RegisterTransactionServiceServer(mysql.GetDataSourceManager())
	mysql.RegisterResource(config.GetATConfig().DSN)

	db, err := gorm.Open(
		dialector.Open(config.GetATConfig().DSN),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			}})
	if err != nil {
		panic(err)
	}
	DB, err := db.DB()
	if err != nil {
		panic(err)
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(20)
	DB.SetConnMaxLifetime(4 * time.Hour)
	if err := DB.Ping(); err != nil {
		panic(err)
	}
	repo := &svc.Repo{DB: db}
	proxyService := svc.ProxyService{SVC: svc.NewSvc(repo)}
	tm.Implement(&proxyService) // 重点，创建代理对象

	r.POST("/transfer", func(c *gin.Context) {
		if err := proxyService.Transfer(c, false); err == nil {
			c.JSON(200, gin.H{
				"success": true,
				"message": "success",
			})
		} else {
			c.JSON(500, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}

	})

	r.Run(":8000")
}
