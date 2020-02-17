package dfuse

import (
	"bytes"
	"context"

	"sync"
	"time"

	"github.com/chenyihui555/dfuse-go/entity"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
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
		logrus.Error("unmarshal err", err)
		return
	}

	if resp.ReqId != s.reqId {
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
	var interval time.Duration
	switch {
	case s.progress.interval <= 10:
		interval = time.Minute
	case 10 <= s.progress.interval && s.progress.interval <= 60:
		interval = s.progress.interval * 3 / 5
	case 60 < s.progress.interval:
		interval = s.progress.interval * 5
	}

	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case nowTime := <-ticker.C:
			if s.progress.isTimeout(nowTime) {
				logrus.Errorf("progress 接收超时 nowTime:%s, progressTime:%s", nowTime.String(), s.progress.nextTime.String())
				s.cli.sendChan <- s.reqCache
			}
		}
	}
}
