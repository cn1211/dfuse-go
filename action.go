package dfuse

import (
	"github.com/pkg/errors"

	"github.com/chenyihui555/dfuse-go/entity"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type action struct {
	*wssClient

	conn *websocket.Conn
}

type callback func(string, string)

var defaultOptionReq = &entity.OptionReq{
	Fetch:            true,
	Listen:           true,
	StartBlock:       0,
	IrreversibleOnly: true,
	WithProgress:     1,
}

func (a *action) GetTableRows(reqId string, request *entity.GetTableRows, handle callback, opt *entity.OptionReq) error {
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

	if err := a.write(reqId, param); err != nil {
		return err
	}

	a.registerCallback(reqId, handle)
	return nil
}

func (a *action) GetTransactionLifecycle(reqId, txHash string, handle callback, opt *entity.OptionReq) error {
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

	if err := a.write(reqId, param); err != nil {
		return err
	}

	a.registerCallback(reqId, handle)
	return nil
}

// interrupt stream
func (a *action) UnListen(reqId string) error {
	param := entity.UnListenReq{
		Type: UnListen,
		Data: struct {
			ReqId string `json:"req_id"`
		}{
			ReqId: reqId,
		},
	}

	writeBytes, err := jsoniter.Marshal(param)
	if err != nil {
		return err
	}

	return a.conn.WriteMessage(websocket.TextMessage, writeBytes)
}

func (a *action) TableSnapshot() (*entity.TableSnapshotResp, error) {
	snapshot := entity.TableSnapshotResp{}
	if err := jsoniter.UnmarshalFromString("", &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (a *action) TableDelta() (*entity.TableDeltaResp, error) {
	delta := entity.TableDeltaResp{}
	if err := jsoniter.UnmarshalFromString("", &delta); err != nil {
		return nil, err
	}
	return &delta, nil
}

func (a *action) TransactionLifecycle() (*entity.TransactionLifecycleResp, error) {
	lifecycleResp := entity.TransactionLifecycleResp{}
	if err := jsoniter.UnmarshalFromString("", &lifecycleResp); err != nil {
		return nil, err
	}
	return &lifecycleResp, nil
}

func (a *action) Progress() (*entity.ProgressResp, error) {
	progressResp := entity.ProgressResp{}
	if err := jsoniter.UnmarshalFromString("", &progressResp); err != nil {
		return nil, err
	}
	return &progressResp, nil
}

func (a *action) Listening() (*entity.ListeningResp, error) {
	listenResp := entity.ListeningResp{}
	if err := jsoniter.UnmarshalFromString("", &listenResp); err != nil {
		return nil, err
	}
	return &listenResp, nil
}

func (a *action) Error() error {
	return errors.New("")
}
