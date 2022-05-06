package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	cron "github.com/robfig/cron/v3"
	"github.com/streadway/amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	userDb := NewMysqlDB("user")
	CreateUserDb(userDb)
	integralDb := NewMysqlDB("integral")
	CreateIntegral(integralDb)
	go startConsumer(integralDb)

	go startCronJobProducer(userDb)
	startServer(userDb)
	select {}
}

func startServer(userDb *gorm.DB) {
	r := gin.Default()
	r.Handle("POST", "/user", func(c *gin.Context) {
		var req UserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, errors.New("params"))
			return
		}

		userRepo := NewUserRepo(userDb)
		if err := userRepo.CreateUserAndLog(User{
			Username: req.Username,
		}, MsgLog{Integral: 100, UUID: GenerateString(), Status: -1}); err != nil {
			panic(err)
		}
	})
	r.Run(":8090")
}

func startCronJobProducer(userDb *gorm.DB) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, err := conn.Channel()

	q, err := ch.QueueDeclare(
		"my_queue", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	c := cron.New(cron.WithSeconds())

	msgRepo := NewMsgLogRepo(userDb)
	_, err = c.AddFunc(EveryMinute, func() {

		tasks, err := msgRepo.FindNeedExec()
		if err != nil {
			panic(err)
		}

		var resolveUUIDs []string
		for _, task := range tasks {
			body, _ := json.Marshal(task)

			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			// mark resolve
			resolveUUIDs = append(resolveUUIDs, task.UUID)
		}

		if err := msgRepo.BatchUpdateResolve(resolveUUIDs); err != nil {
			// do nothing
			fmt.Println(err.Error())
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}

func startConsumer(integralDb *gorm.DB) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"my_queue", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,     // queue
		"my_queue", // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)

	integralRepo := NewIntegralRepo(integralDb)
	msgLogRepo := NewMsgLogRepo(integralDb)
	go func() {
		for d := range msgs {
			body := string(d.Body)
			fmt.Printf("consumed: %v", body)
			var msg MsgLog
			exists, err2 := msgLogRepo.Exists(msg.UUID)
			if err2 != nil {
				panic(err)
			}
			if exists {
				continue // 幂等处理
			}

			if err := json.Unmarshal(d.Body, &msg); err != nil {
				panic(err)
			}

			if err := integralRepo.CreateIntegralAndLog(Integral{
				UserId:   msg.UserId,
				Integral: msg.Integral,
			}, msg); err != nil {
				panic(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func NewMysqlDB(database string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		"root",
		"123456",
		"localhost",
		3306,
		database,
		"utf8mb4",
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("start mysql client failed, db:%s, err:%s" + err.Error())
	}
	return db
}

func CreateUserDb(db *gorm.DB) {
	_ = db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(
		&MsgLog{},
		&User{},
	)
}

func CreateIntegral(db *gorm.DB) {
	_ = db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(
		&MsgLog{},
		&Integral{},
	)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
