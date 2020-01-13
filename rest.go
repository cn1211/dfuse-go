package dfuse

type restClient struct {
	cli *Client
}

func newRestClient() *restClient {
	return &restClient{}
}
