package dfuse

import (
	"bytes"
	"context"
	"log"
	"sync"
	"time"

	"github.com/chenyihui555/dfuse-go/entity"
	jsoniter "github.com/json-iterator/go"
)

type callbackFunc func(string, *Callback)

type subscribe struct {
	cli *wssClient
	hub *Hub

	closeOnce sync.Once
	ctx       context.Context
	cancel    context.CancelFunc

	reqId      string
	actionType string
	handle     callbackFunc
	reqCache   []byte
	progress   *progress
	callback   *Callback
}

func newSubscribe(reqId, actionType string, intervalBlock int, param interface{}, handle callbackFunc, cli *wssClient, hub *Hub) (*subscribe, error) {
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
		progress:   newProgress(reqId, time.Second*time.Duration(intervalBlock)/2),
		callback:   newCallback(),
	}

	subscriber.ctx, subscriber.cancel = context.WithCancel(context.Background())

	subscriber.cli.sendChan <- msgBytes
	subscriber.hub.register <- &subscriber

	go subscriber.monitorProgress()
	return &subscriber, nil
}

func (s *subscribe) distribute(sendBytes []byte) {
	var resp entity.CommonResp
	if err := jsoniter.Unmarshal(sendBytes, &resp); err != nil {
		log.Println("unmarshal err :", err)
		return
	}

	if resp.Type == Ping {
		pong := bytes.Replace(sendBytes, []byte(`"ping"`), []byte(`"pong"`), 1)
		s.cli.sendChan <- pong
		return
	}

	if resp.Type == Progress {
		s.progress.refreshTime()
	}

	s.callback.msgBytes = sendBytes
	s.handle(resp.Type, s.callback)
}

func (s *subscribe) Close() {
	s.closeOnce.Do(func() {
		s.cancel()
	})
}

func (s *subscribe) monitorProgress() {
	ticker := time.NewTicker(s.progress.interval * 9 / 10)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case nowTime := <-ticker.C:
			if s.progress.isTimeout(nowTime) {
				s.cli.sendChan <- s.reqCache
			}
		}
	}
}
