package entity

type StateTableReq struct {
	Account      string `json:"account"`
	Scope        string `json:"scope"`
	Table        string `json:"table"`
	BlockNum     int    `json:"block_num,omitempty"`
	Json         bool   `json:"json,omitempty"`
	KeyType      string `json:"key_type,omitempty"`
	WithBlockNum int    `json:"with_block_num,omitempty"`
	WithAbi      bool   `json:"with_abi,omitempty"`
}

type StateTableResp struct {
	UpToBlockId              int64    `json:"up_to_block_id,omitempty"`
	UpToBlockNum             int64    `json:"up_to_block_num,omitempty"`
	LastIrreversibleBlockId  string   `json:"last_irreversible_block_id,omitempty"`
	LastIrreversibleBlockNum int64    `json:"last_irreversible_block_num,omitempty"`
	Abi                      struct{} `json:"abi,omitempty"`
	Rows                     *[]DBRow `json:"rows"`
}
