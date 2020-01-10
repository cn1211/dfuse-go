package dfuse

import (
	"sync"
)

type Client struct {
	//*baseClient
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
		c.wssCli = newWssClient(c.Network.Endpoint, c.tokenStore.GetAuth().Token)
	})

	return c.wssCli
}

// TODO
func (c *Client) Rest() *restClient {
	c.restOnce.Do(func() {

	})
	return nil
}
