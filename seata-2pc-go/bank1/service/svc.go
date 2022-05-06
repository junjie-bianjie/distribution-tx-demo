package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	context2 "github.com/opentrx/seata-golang/v2/pkg/client/base/context"
	"github.com/opentrx/seata-golang/v2/pkg/client/base/model"
)

type SVC struct {
	repo *Repo
}

func NewSvc(repo *Repo) *SVC {
	return &SVC{repo: repo}
}

func (svc *SVC) Transfer(ctx context.Context, rollback bool) error {
	rootContext := ctx.(*context2.RootContext)
	event := AccountEvent{
		AccountNo: 2, // 给李四转账
		Amount:    1,
	}

	_, err := svc.repo.Transfer(ctx, &event)
	if err != nil {
		panic(err)
	}

	req, err := json.Marshal(event)
	fmt.Println(string(req))
	req1, err := http.NewRequest("POST", "http://localhost:8001/add-balances", bytes.NewBuffer(req))
	if err != nil {
		panic(err)
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("xid", rootContext.GetXID())

	client := &http.Client{}
	result1, err1 := client.Do(req1)
	if err1 != nil {
		return err1
	}

	if result1.StatusCode != 200 {
		return errors.New("err")
	}

	if rollback {
		return errors.New("there is a error")
	}
	return nil
}

type ProxyService struct {
	*SVC
	Transfer func(ctx context.Context, rollback bool) error
}

var methodTransactionInfo = make(map[string]*model.TransactionInfo)

func init() {
	methodTransactionInfo["Transfer"] = &model.TransactionInfo{
		TimeOut:     60000000,
		Name:        "Transfer",
		Propagation: model.Required,
	}
}

func (svc *ProxyService) GetProxyService() interface{} {
	return svc.SVC
}

func (svc *ProxyService) GetMethodTransactionInfo(methodName string) *model.TransactionInfo {
	return methodTransactionInfo[methodName]
}
