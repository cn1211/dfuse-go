package dfuse

import "fmt"

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

func (n Network) GQLEndpoint() string {
	return fmt.Sprintf("https://%s/graphql", n.Endpoint)
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
