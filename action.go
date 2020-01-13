package dfuse

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/chenyihui555/dfuse-go/entity"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
)

type WssActionInter interface {
	wssRequestAction
	wssResponseAction
}

type wssRequestAction interface {
	GetTableRows(*entity.GetTableRows, HandleCallback, *entity.OptionReq) error
	GetTransactionLifecycle(string, HandleCallback, *entity.OptionReq) error
}

type wssResponseAction interface {
	UnListen(string) error
	TableSnapshot() (*entity.TableSnapshotResp, error)
	TableDelta() (*entity.TableDeltaResp, error)
	TransactionLifecycle() (*entity.TransactionLifecycleResp, error)
	Progress() (*entity.ProgressResp, error)
	Listening() (*entity.ListeningResp, error)
	UnListened()
	Error() error
}

type action struct {
	readBytes []byte
	mux       sync.Mutex
	conn      *websocket.Conn
	handle    map[string]HandleCallback // map[请求id]回调
}

type HandleCallback func(msgType string, payload Payload)

func (a *action) GetTableRows(request *entity.GetTableRows, handle HandleCallback, opt *entity.OptionReq) error {
	param := entity.TableRowsReq{
		CommonReq: entity.CommonReq{
			Type:      GetTableRows,
			ReqId:     uuid.NewV4().String(),
			OptionReq: opt,
		},
		Data: *request,
	}

	if err := a.write(param); err != nil {
		return err
	}

	a.handle[param.ReqId] = handle
	return nil
}

func (a *action) GetTransactionLifecycle(txHash string, handle HandleCallback, opt *entity.OptionReq) error {
	param := entity.TransactionLifecycleReq{
		CommonReq: entity.CommonReq{
			Type:      TransactionLifecycle,
			ReqId:     uuid.NewV4().String(),
			OptionReq: opt,
		},
		Data: struct {
			Id string `json:"id"`
		}{
			Id: txHash,
		},
	}

	if err := a.write(param); err != nil {
		return err
	}

	a.handle[param.ReqId] = handle
	return nil
}

func (a *action) UnListen(reqId string) error {
	param := entity.UnListenReq{
		Type: UnListen,
		Data: struct {
			ReqId string `json:"req_id"`
		}{
			ReqId: reqId,
		},
	}

	return a.write(param)
}

func (a *action) TableSnapshot() (*entity.TableSnapshotResp, error) {
	panic("implement me")
}

func (a *action) TableDelta() (*entity.TableDeltaResp, error) {
	panic("implement me")
}

func (a *action) TransactionLifecycle() (*entity.TransactionLifecycleResp, error) {
	panic("implement me")
}

func (a *action) Progress() (*entity.ProgressResp, error) {
	panic("implement me")
}

func (a *action) Listening() (*entity.ListeningResp, error) {
	panic("implement me")
}

func (a *action) UnListened() {
	panic("implement me")
}

func (a *action) Error() error {
	return errors.New(string(a.readBytes))
}

func (a *action) write(param interface{}) error {
	_ = a.conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
	writeBytes, err := jsoniter.Marshal(param)
	if err != nil {
		return err
	}

	err = a.conn.WriteMessage(websocket.TextMessage, writeBytes)
	if err != nil {
		return err
	}

	return nil
}

func (a *action) callback(respType, reqId, context string) {
	if handle, has := a.handle[reqId]; has {
		handle(respType, Payload{context: context})
	}
}

type Payload struct {
	context string
}

func (p *Payload) TableSnapshot() (*entity.TableSnapshotResp, error) {
	snapshot := entity.TableSnapshotResp{}
	if err := jsoniter.UnmarshalFromString(p.context, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (p *Payload) TableDelta() (*entity.TableDeltaResp, error) {
	delta := entity.TableDeltaResp{}
	if err := jsoniter.UnmarshalFromString(p.context, &delta); err != nil {
		return nil, err
	}
	return &delta, nil
}

func (p *Payload) TransactionLifecycle() (*entity.TransactionLifecycleResp, error) {
	lifecycleResp := entity.TransactionLifecycleResp{}
	if err := jsoniter.UnmarshalFromString(p.context, &lifecycleResp); err != nil {
		return nil, err
	}
	return &lifecycleResp, nil
}

func (p *Payload) Progress() (*entity.ProgressResp, error) {
	progressResp := entity.ProgressResp{}
	if err := jsoniter.UnmarshalFromString(p.context, &progressResp); err != nil {
		return nil, err
	}
	return &progressResp, nil
}

func (p *Payload) Listening() (*entity.ListeningResp, error) {
	listenResp := entity.ListeningResp{}
	if err := jsoniter.UnmarshalFromString(p.context, &listenResp); err != nil {
		return nil, err
	}
	return &listenResp, nil
}

func (p *Payload) Error() error {
	return errors.New(p.context)
}
