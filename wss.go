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
	cli  *Client
	conn *websocket.Conn

	closeOnce sync.Once
	msgChan   chan []byte
	err       error
	*action
}

func newWssClient(endpoint, token string) *wssClient {
	wssCli := &wssClient{
		msgChan: make(chan []byte),
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
	wssCli.action = &action{
		readBytes: make([]byte, 0),
		conn:      conn,
		handle:    make(map[string]HandleCallback),
	}

	go wssCli.read()
	go wssCli.publish()
	return wssCli
}

func (c *wssClient) Error() error {
	return c.err
}

func (c *wssClient) Close() {
	c.closeOnce.Do(func() {
		close(c.msgChan)
		_ = c.conn.Close()
	})
}

func (c *wssClient) read() {
	//c.conn.SetReadLimit(maxBufferSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pingWait))
	c.conn.SetPingHandler(func(appData string) error {
		fmt.Println(">>>执行到此")
		pong := bytes.Replace([]byte(appData), []byte(`"ping"`), []byte(`"pong"`), 1)
		_ = c.conn.WriteMessage(websocket.TextMessage, pong)
		_ = c.conn.SetReadDeadline(time.Now().Add(pingWait))
		return nil
	})

	for {
		msgType, context, err := c.conn.ReadMessage()
		if err != nil {
			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() {
					c.err = fmt.Errorf("read message timeout remote :%v", c.conn.RemoteAddr())
					continue
				}
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Printf("read message fail expect err:%v \n", err)
				continue
			}
			c.err = err
			continue
		}

		if msgType != websocket.TextMessage {
			c.err = fmt.Errorf("invalid msg type:%d", msgType)
			continue
		}

		c.msgChan <- context
	}
}

func (c *wssClient) publish() {
	for msg := range c.msgChan {
		var resp entity.CommonResp
		if err := jsoniter.Unmarshal(msg, &resp); err != nil {
			c.err = fmt.Errorf("unmarshal err:%v", err)
			continue
		}

		switch resp.Type {
		case Progress:

		case UnListened:
			c.action.UnListened()
		case Error:
			c.action.Error()
		default:
			//fmt.Printf("receive msg:%s \n", string(msg))
			c.action.callback(resp.Type, resp.ReqId, string(msg))
		}
	}
}
