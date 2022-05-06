# seata-go-samples

## Step0: setup TC server
```bash
git clone git@github.com:opentrx/seata-golang.git
cd seata-golang

vim ./cmd/profiles/dev/config.yml
# update storage.mysql.dsn
# update log.logPath

# create database `seata` on mysql server
# mysql> CREATE database if NOT EXISTS `seata` default character set utf8mb4 collate utf8mb4_unicode_ci;

cd cmd/tc
go run main.go start -config ../profiles/dev/config.yml
```

- ## AT mode example (gorm)
### Step1: setup bank1 client
```bash
vim ./bank1/conf/client.yml
# update log.logPath

export ConfigPath="./bank1/conf/client.yml"
go run bank1/main.go
```

### Step2: setup bank2 client
```bash
vim ./bank2/conf/client.yml
# update at.dsn
# update log.logPath

export ConfigPath="./bank2/conf/client.yml"
go run bank2/main.go
```

### Step3: access
- http://localhost:8000/transfer
