package dfuse

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/chenyihui555/dfuse-go/entity"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type wssClient struct {
	cli  *Client
	conn *websocket.Conn

	mutex     *sync.Mutex
	msgChan   chan interface{}
	closeOnce sync.Once
	handle    map[string]func()
}

var wssCli *wssClient

func newWssClient(endpoint, token string) *wssClient {
	wssCli = &wssClient{
		msgChan: make(chan interface{}, 0),
		handle:  make(map[string]func(), 0),
	}

	dail := websocket.Dialer{
		HandshakeTimeout:  time.Second * 15,
		EnableCompression: true,
	}

	u := fmt.Sprintf("%s?token=%s", endpoint, token)
	conn, _, err := dail.Dial(u, nil)
	if err != nil {
		panic(err)
	}

	wssCli.conn = conn

	go wssCli.read()
	return wssCli
}

func (c *wssClient) Close() {
	c.closeOnce.Do(func() {
		_ = c.conn.Close()
	})
}

func (c *wssClient) read() {
	for {
		msgType, cnt, err := c.conn.ReadMessage()
		if err != nil {
			panic(fmt.Sprintf("read message err%v", err))
		}

		if msgType != websocket.TextMessage {
			panic(fmt.Sprintf("msgtype invaild type%d", msgType))
		}

		var resp entity.CommonResp
		if err := jsoniter.Unmarshal(cnt, &resp); err != nil {
			panic(fmt.Sprintf("unmarshal err%v", err))
		}

		if resp.Type == "ping" {
			pong := bytes.Replace(cnt, []byte(`"ping"`), []byte(`"pong"`), 1)
			_ = c.conn.WriteMessage(websocket.TextMessage, pong)
			continue
		}

		fmt.Printf("receive msg:%s \n", string(cnt))
	}
}

func (c *wssClient) write(data interface{}) error {
	reqBytes, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		return err
	}

	return nil
}
