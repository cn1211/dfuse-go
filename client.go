package dfuse

import (
	"sync"
)

type Client struct {
	*Options

	wssOnce  sync.Once
	restOnce sync.Once

	wssCli  *wssClient
	restCli *restClient
}

func NewClient(opt *Options) *Client {
	if opt == nil {
		opt = defaultOpt
	}

	opt.init()
	return &Client{Options: opt}
}

func (c *Client) Wss() *wssClient {
	c.wssOnce.Do(func() {
		c.wssCli = newWssClient(c.Network.WssEndPoint(), c.tokenStore.GetAuth().Token, c)
	})

	return c.wssCli
}

// TODO 待定
func (c *Client) Rest() *restClient {
	c.restOnce.Do(func() {
		c.restCli = newRestClient()
	})
	return c.restCli
}
