package dfuse

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chenyihui555/dfuse-go/entity"
	jsoniter "github.com/json-iterator/go"
)

type callback func(string, string)

type subscribe struct {
	cli *wssClient
	hub *Hub

	closeOnce sync.Once
	ctx       context.Context
	cancel    context.CancelFunc

	reqId      string
	actionType string
	handle     callback
	reqCache   []byte
	progress   *progress
}

func newSubscribe(reqId, actionType string, param interface{}, handle callback, cli *wssClient, hub *Hub) (*subscribe, error) {
	msgBytes, err := jsoniter.Marshal(param)
	if err != nil {
		return nil, err
	}

	subscriber := subscribe{
		cli:        cli,
		hub:        hub,
		reqCache:   msgBytes,
		reqId:      reqId,
		actionType: actionType,
		handle:     handle,
		progress:   newProgress(reqId, time.Second*10),
	}

	subscriber.ctx, subscriber.cancel = context.WithCancel(context.Background())

	subscriber.cli.sendChan <- msgBytes
	subscriber.hub.register <- &subscriber

	go subscriber.monitorProgress()
	return &subscriber, nil
}

// TODO 回调处理
func (s *subscribe) callback(sendBytes []byte) {
	var resp entity.CommonResp
	if err := jsoniter.Unmarshal(sendBytes, &resp); err != nil {
		fmt.Println("unmarshal err :", err)
		return
	}

	switch resp.Type {
	case Progress:
		s.progress.refreshTime()

	case Ping:
		pong := bytes.Replace(sendBytes, []byte(`"ping"`), []byte(`"pong"`), 1)
		s.cli.sendChan <- pong

	case UnListened:
		fmt.Printf("unlisten success msg:%s \n", string(sendBytes))

	case Listening:
		fmt.Printf("listening msg:%s \n", string(sendBytes))

	case Error:
		fmt.Printf("err:%s \n", string(sendBytes))

	default:
		s.handle(resp.Type, string(sendBytes))
	}
}

func (s *subscribe) Close() {
	s.closeOnce.Do(func() {
		s.cancel()
	})
}

func (s *subscribe) monitorProgress() {
	ticker := time.NewTicker(time.Second * 10)
	defer func() {
		ticker.Stop()
	}()

	select {
	case <-s.ctx.Done():
		return

	case nowTime := <-ticker.C:
		if s.progress.isTimeout(nowTime) {
			s.cli.sendChan <- s.reqCache
		}
	}
}
