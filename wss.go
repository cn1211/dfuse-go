package dfuse

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/chenyihui555/dfuse-go/entity"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	pingWait = time.Second * 60
)

type wssClient struct {
	*Client
	*memoryCache
	*action

	conn        *websocket.Conn
	dail        websocket.Dialer
	mux         sync.Mutex
	closeOnce   sync.Once
	sendChan    chan []byte
	handle      map[string]callback
	cacheReq    map[string][]byte
	progressMap map[string]progress
}

func newWssClient(endpoint, token string, cli *Client) *wssClient {
	wssCli := &wssClient{
		Client: cli,

		handle:      make(map[string]callback),
		memoryCache: NewMemoryCache(time.Second*30, time.Minute*5),
		sendChan:    make(chan []byte),
		cacheReq:    make(map[string][]byte),
		progressMap: make(map[string]progress),
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
	wssCli.action = &action{
		wssClient: wssCli,
		conn:      conn,
	}

	go wssCli.read()
	go wssCli.publish()
	go wssCli.listeningProgress()
	return wssCli
}

func (w *wssClient) Close() {
	w.closeOnce.Do(func() {
		w.batchUnListen()
		close(w.action.sendChan)
		_ = w.conn.Close()
	})
}

func (w *wssClient) read() {
	_ = w.conn.SetReadDeadline(time.Now().Add(pingWait))
	w.conn.SetPingHandler(func(appData string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(pingWait))
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
				fmt.Printf("read message fail expect err:%v \n", err)
				w.reconnect()
				continue
			}

			fmt.Println("read message fail ", err)
			w.reconnect()
			continue
		}

		if msgType != websocket.TextMessage {
			fmt.Printf("invalid msg type:%d \n", msgType)
			continue
		}

		w.sendChan <- cnt
	}
}

func (w *wssClient) reconnect() {
	w.Options.refreshToken()
	u := fmt.Sprintf("%s?token=%s", w.Network.WssEndPoint(), w.tokenStore.GetAuth().Token)
	var err error
	conn, _, err := w.dail.Dial(u, nil)
	if err != nil {
		panic(fmt.Sprintf("dail fail err:%v", err))
	}

	w.batchUnListen()
	_ = w.conn.Close()
	w.conn = conn
	w.resubscribe()

	fmt.Println("reconnect success")
}

func (w *wssClient) publish() {
	for msg := range w.sendChan {
		var resp entity.CommonResp
		if err := jsoniter.Unmarshal(msg, &resp); err != nil {
			panic(fmt.Sprint("unmarshal err", err))
		}

		switch resp.Type {
		case Progress:
			fmt.Println("progress:", string(msg))
			if progress, has := w.progressMap[resp.ReqId]; has {
				progress.refreshTime()
			}

		case Ping:
			pong := bytes.Replace([]byte(msg), []byte(`"ping"`), []byte(`"pong"`), 1)
			fmt.Println("pong:", string(pong))
			fmt.Println("ping:", string(msg))
			_ = w.conn.WriteMessage(websocket.TextMessage, pong)

		case UnListened:
			fmt.Println("unlistend:", string(msg))

		case Listening:
			fmt.Println("listening:", string(msg))

		case Error:
			fmt.Println("error:", string(msg))

		default:
			fmt.Println(">>>消息msg", string(msg))
			handle, has := w.handle[resp.ReqId]
			if has {
				handle(resp.Type, string(msg))
			}
		}
	}
}

func (w *wssClient) write(reqId string, param interface{}) error {
	writeBytes, err := jsoniter.Marshal(param)
	if err != nil {
		return err
	}

	if w.isRepeatReq(reqId, writeBytes) {
		return nil
	}

	w.mux.Lock()
	defer w.mux.Unlock()

	err = w.conn.WriteMessage(websocket.TextMessage, writeBytes)
	if err != nil {
		return err
	}

	w.cacheReq[reqId] = writeBytes // 记录用户请求数据,用于重新订阅
	w.progressMap[reqId] = newProgress(reqId, time.Second*5)
	return nil
}

func (w *wssClient) batchUnListen() {
	for reqId := range w.handle {
		if err := w.UnListen(reqId); err != nil {
			fmt.Printf("unlisten err :%v \n", err)
			continue
		}
	}
}

func (w *wssClient) resubscribe() {
	if w.conn == nil {
		return
	}

	for reqId, writeBytes := range w.cacheReq {
		if err := w.conn.WriteMessage(websocket.TextMessage, writeBytes); err != nil {
			fmt.Println("write message fail", err)
			continue
		}

		w.progressMap[reqId] = newProgress(reqId, time.Second*5)
	}
}

func (w *wssClient) registerCallback(reqId string, handle callback) {
	w.handle[reqId] = handle
}

func (w *wssClient) listeningProgress() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for nowTime := range ticker.C {
		for reqId, progress := range w.progressMap {
			if !progress.isTimeout(nowTime) {
				continue
			}

			if err := w.UnListen(reqId); err != nil {
				fmt.Printf("unlisten fail err:%+v \n", err)
				continue
			}

			writeBytes, has := w.cacheReq[reqId]
			if !has {
				fmt.Println("reconnect fail not exist req cache ")
				continue
			}

			if err := w.conn.WriteMessage(websocket.TextMessage, writeBytes); err != nil {
				fmt.Println("write message fail", err)
				continue
			}

			progress.refreshTime()
		}
	}
}

// checks for repeat requests
func (w *wssClient) isRepeatReq(reqId string, writeBytes []byte) bool {
	return reflect.DeepEqual(w.cacheReq[reqId], writeBytes)
}
