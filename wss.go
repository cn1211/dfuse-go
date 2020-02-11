package dfuse

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = time.Second * 60
	pingPeriod = (pongWait * 9) / 10
)

type wssClient struct {
	*Client
	*action

	conn      *websocket.Conn
	dail      websocket.Dialer
	closeOnce sync.Once
	sendChan  chan []byte
	cacheReq  map[string][]byte
	errChan   chan error
	notify    chan os.Signal
	err       error

	hub           *Hub
	subscriberMap map[string]*subscribe
}

func newWssClient(endpoint, token string, cli *Client) *wssClient {
	wssCli := &wssClient{
		Client: cli,

		sendChan:      make(chan []byte),
		cacheReq:      make(map[string][]byte),
		errChan:       make(chan error),
		notify:        make(chan os.Signal, 1),
		subscriberMap: make(map[string]*subscribe),
	}

	wssCli.dail = websocket.Dialer{
		HandshakeTimeout:  time.Second * 30,
		EnableCompression: true,
	}

	u := fmt.Sprintf("%s?token=%s", endpoint, token)
	conn, _, err := wssCli.dail.Dial(u, nil)
	if err != nil {
		panic(err)
	}

	wssCli.conn = conn
	wssCli.action = &action{wssClient: wssCli}
	wssCli.hub = newHub(wssCli)

	signal.Notify(wssCli.notify, os.Interrupt)

	go wssCli.hub.run()
	go wssCli.read()
	go wssCli.write()

	return wssCli
}

func (w *wssClient) Close() {
	w.closeOnce.Do(func() {
		close(w.sendChan)
		close(w.errChan)
		_ = w.conn.Close()
	})
}

func (w *wssClient) read() {
	defer func() {
		w.Close()
	}()

	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPingHandler(func(appData string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		msgType, cnt, err := w.conn.ReadMessage()
		if err != nil {
			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() {
					fmt.Println("time out ", time.Now().String())
					w.reconnect()
					continue
				}
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// goingaway:页面跳转  abnormal:断网.没有收到帧
				fmt.Printf("read message fail expect err:%v \n", err)
				w.reconnect()
				continue
			}

			fmt.Println("read message fail ", err)
			continue
		}

		if msgType != websocket.TextMessage {
			fmt.Printf("invalid msg type:%d \n", msgType)
			continue
		}

		w.hub.broadcast <- cnt
	}
}

func (w *wssClient) reconnect() {
	w.Options.refreshToken()
	u := fmt.Sprintf("%s?token=%s", w.Network.WssEndPoint(), w.tokenStore.GetAuth().Token)

	conn, _, err := w.dail.Dial(u, nil)
	if err != nil {
		panic(fmt.Sprintf("dail fail err:%v", err))
	}

	_ = w.conn.Close()
	w.conn = conn
	w.resubscribe()

	fmt.Println("reconnect success")
}

func (w *wssClient) write() {
	defer func() {
		w.Close()
	}()

	for {
		select {
		case <-w.notify:
			fmt.Println("interrupt")
			err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close message err", err)
				w.errChan <- err
				w.err = err
				return
			}

		case msg := <-w.sendChan:
			err := w.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("write text message err", err)
				return
			}
		}
	}

}

// resubscribe
func (w *wssClient) resubscribe() {
	for reqId, subscriber := range w.subscriberMap {
		w.sendChan <- subscriber.reqCache
		subscriber.progress = newProgress(reqId, time.Second*10)
	}
}

// subscribe
func (w *wssClient) subscribe(reqId, actionType string, param interface{}, handle callback) error {
	subscriber, err := newSubscribe(reqId, actionType, param, handle, w, w.hub)
	if err != nil {
		return err
	}
	w.wssCli.subscriberMap[reqId] = subscriber
	w.hub.register <- subscriber
	return nil
}

// unsubscribe
func (w *wssClient) unsubscribe(reqId string) {
	if subscriber, has := w.subscriberMap[reqId]; has {
		w.hub.unregister <- subscriber
		delete(w.subscriberMap, reqId)
		subscriber.Close()
	}
}
