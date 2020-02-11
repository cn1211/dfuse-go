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
