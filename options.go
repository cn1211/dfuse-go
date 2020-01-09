package dfuse

type Options struct {
	Network    Network
	ApiKey     string
	Proxy      string
	TokenStore string
}

func (opt *Options) init() {

}

type Network struct {
	Name     string
	Endpoint string
}
