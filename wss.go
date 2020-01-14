package dfuse

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/chenyihui555/dfuse-go/entity"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	maxBufferSize = 2024
	pingWait      = time.Second * 60
)

type wssClient struct {
	*Client
	*memoryCache
	*action

	conn      *websocket.Conn
	dail      websocket.Dialer
	mux       sync.Mutex
	closeOnce sync.Once
	sendChan  chan []byte
	handle    map[string]handleCallback
}

func newWssClient(endpoint, token string, cli *Client) *wssClient {
	wssCli := &wssClient{
		Client: cli,

		handle:      make(map[string]handleCallback, 0),
		memoryCache: NewMemoryCache(time.Second*30, time.Minute*5),
		sendChan:    make(chan []byte, 0),
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
	return wssCli
}

func (w *wssClient) Close() {
	w.closeOnce.Do(func() {
		close(w.action.sendChan)
		_ = w.conn.Close()
	})
}

func (w *wssClient) read() {

	w.conn.SetReadLimit(maxBufferSize)
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
					fmt.Println("time out")
					w.reconnect()
					continue
				}
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				panic(fmt.Sprintf("read message fail expect err:%v", err))
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

	_ = w.conn.Close()
	w.conn = conn
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

		case Ping:
			pong := bytes.Replace([]byte(msg), []byte(`"ping"`), []byte(`"pong"`), 1)
			_ = w.conn.WriteMessage(websocket.TextMessage, pong)
			fmt.Println("ping:", string(msg))

		case UnListened:
			fmt.Println("unlistend:", string(msg))

		case Listening:
			fmt.Println("listening:", string(msg))

		case Error:
			fmt.Println("error:", string(msg))

		default:
			handle, has := w.handle[resp.ReqId]
			if has {
				handle(resp.Type, string(msg))
			}
		}
	}
}

func (w *wssClient) write(param interface{}) error {
	_ = w.conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
	writeBytes, err := jsoniter.Marshal(param)
	if err != nil {
		return err
	}

	w.mux.Lock()
	defer w.mux.Unlock()
	return w.conn.WriteMessage(websocket.TextMessage, writeBytes)
}
