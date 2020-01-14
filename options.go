package dfuse

import "fmt"

var defaultOpt = &Options{}

type Options struct {
	Network    Network
	ApiKey     string
	Proxy      string
	tokenStore TokenStore
}

func (opt *Options) init() {
	opt.refreshToken()
}

// token refresh
func (opt *Options) refreshToken() {
	opt.tokenStore = &InMemoryTokenStore{}
	auth := opt.tokenStore.GetAuth()
	if auth == nil || auth.IsExpired() {
		auth, err := fetchAuth(opt.ApiKey)
		if err != nil {
			panic(fmt.Sprintf("auth fail err:%v", err))
		}
		opt.tokenStore.SetAuth(auth)
	}
}

type Network struct {
	Name     string
	Endpoint string
}

func (n *Network) WssEndPoint() string {
	return fmt.Sprintf("wss://%s/v1/stream", n.Endpoint)
}

func (n *Network) RestEndPoint() string {
	return fmt.Sprintf("https://%s", n.Endpoint)
}

var MainNet = Network{
	Name:     "mainnet",
	Endpoint: "mainnet.eos.dfuse.io",
}

var Jungle = Network{
	Name:     "jungle",
	Endpoint: "jungle.eos.alt.dfuse.io",
}

var Kylin = Network{
	Name:     "kylin",
	Endpoint: "kylin.eos.dfuse.io",
}
