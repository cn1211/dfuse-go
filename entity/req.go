package entity

type CommonReq struct {
	Type             string `json:"type"`
	ReqId            string `json:"req_id,omitempty"`
	Fetch            bool   `json:"fetch,omitempty"`
	Listen           bool   `json:"listen,omitempty"`
	StartBlock       int    `json:"start_block,1,omitempty"`
	IrreversibleOnly bool   `json:"irreversible_only,omitempty"`
	WithProgress     int    `json:"with_progress,omitempty"`
}

type TableRowsReq struct {
	CommonReq

	Data struct {
		Code  string `json:"code"`
		Scope string `json:"scope"`
		Table string `json:"table"`
		Json  bool   `json:"json,omitempty"`
	} `json:"data"`
}
