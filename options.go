package dfuse

var defaultOpt = &Options{}

type Options struct {
	Network    Network
	ApiKey     string
	Proxy      string
	tokenStore TokenStore
}

func (opt *Options) init() {
	auth := opt.tokenStore.GetAuth()
	if auth == nil || auth.IsExpired() {
		auth, err := fetchAuth(opt.ApiKey)
		if err != nil {
			return
		}
		opt.tokenStore.SetAuth(auth)
	}
}

type Network struct {
	Name     string
	Endpoint string
}
