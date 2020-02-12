package dfuse

import (
	"fmt"

	"github.com/chenyihui555/dfuse-go/entity"
	jsoniter "github.com/json-iterator/go"
)

type action struct {
	*wssClient
}

var defaultOptionReq = &entity.OptionReq{
	Fetch:            true,
	Listen:           true,
	StartBlock:       0,
	IrreversibleOnly: true,
	WithProgress:     1,
}

func (a *action) GetTableRows(reqId string, request *entity.GetTableRows, handle callbackFunc, opt *entity.OptionReq) error {
	if opt == nil {
		opt = defaultOptionReq
	}

	param := entity.TableRowsReq{
		CommonReq: entity.CommonReq{
			Type:      GetTableRows,
			ReqId:     reqId,
			OptionReq: opt,
		},
		Data: *request,
	}

	if _, has := a.wssCli.subscribeMap[reqId]; has {
		return fmt.Errorf("req id already exists :%s", reqId)
	}

	return a.subscribe(reqId, GetTableRows, opt.WithProgress, param, handle)
}

func (a *action) GetTransactionLifecycle(reqId, txHash string, handle callbackFunc, opt *entity.OptionReq) error {
	if opt == nil {
		opt = defaultOptionReq
	}

	param := entity.TransactionLifecycleReq{
		CommonReq: entity.CommonReq{
			Type:      TransactionLifecycle,
			ReqId:     reqId,
			OptionReq: opt,
		},
		Data: struct {
			Id string `json:"id"`
		}{
			Id: txHash,
		},
	}

	if _, has := a.wssCli.subscribeMap[reqId]; has {
		return fmt.Errorf("req id already exists :%s", reqId)
	}

	return a.subscribe(reqId, TransactionLifecycle, opt.WithProgress, param, handle)
}

// interrupt subscribe stream
func (a *action) UnListen(reqId string) error {
	param := entity.UnListenReq{
		Type: UnListen,
		Data: struct {
			ReqId string `json:"req_id"`
		}{
			ReqId: reqId,
		},
	}

	sendBytes, err := jsoniter.Marshal(param)
	if err != nil {
		return err
	}

	a.sendChan <- sendBytes

	a.unsubscribe(reqId)
	return nil
}
