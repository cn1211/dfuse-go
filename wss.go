package dfuse

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait = time.Second * 60
	//pingPeriod = (pongWait * 9) / 10
)

type wssClient struct {
	*Client
	*action

	conn      *websocket.Conn
	dail      websocket.Dialer
	closeOnce sync.Once
	sendChan  chan []byte
	cacheReq  map[string][]byte
	notify    chan os.Signal

	hub          *Hub
	subscribeMap map[string]*subscribe
}

func newWssClient(endpoint, token string, cli *Client) *wssClient {
	wssCli := &wssClient{
		Client: cli,

		sendChan:     make(chan []byte),
		cacheReq:     make(map[string][]byte),
		notify:       make(chan os.Signal, 1),
		subscribeMap: make(map[string]*subscribe),
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
					log.Println("time out ", time.Now().String())
					w.reconnect()
					continue
				}
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// goingaway:页面跳转  abnormal:断网.没有收到帧
				log.Printf("read message fail expect err:%v \n", err)
				w.reconnect()
				continue
			}

			log.Println("read message fail ", err)
			continue
		}

		if msgType != websocket.TextMessage {
			log.Printf("invalid msg type:%d \n", msgType)
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
		log.Println("websocket dail fail :", err)
		time.Sleep(time.Second * 3)
		w.reconnect()
	}

	_ = w.conn.Close()
	w.conn = conn
	w.resubscribe()

	log.Println("reconnect success")
}

func (w *wssClient) write() {
	defer func() {
		w.Close()
	}()

	for {
		select {
		case <-w.notify:
			log.Println("interrupt")
			err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close message err", err)
				return
			}

		case msg := <-w.sendChan:
			//log.Println("send chan msg :", string(msg))
			err := w.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write text message err", err)
				return
			}
		}
	}

}

// resubscribe
func (w *wssClient) resubscribe() {
	for reqId, subscriber := range w.subscribeMap {
		w.sendChan <- subscriber.reqCache
		subscriber.progress = newProgress(reqId, subscriber.progress.interval)
	}
}

// subscribe
func (w *wssClient) subscribe(reqId, actionType string, intervalBlock int, param interface{}, handle callbackFunc) error {
	subscriber, err := newSubscribe(reqId, actionType, intervalBlock, param, handle, w, w.hub)
	if err != nil {
		return err
	}
	w.wssCli.subscribeMap[reqId] = subscriber
	return nil
}

// unsubscribe
func (w *wssClient) unsubscribe(reqId string) {
	if subscriber, has := w.subscribeMap[reqId]; has {
		w.hub.unregister <- subscriber
		delete(w.subscribeMap, reqId)
		subscriber.Close()
	}
}
