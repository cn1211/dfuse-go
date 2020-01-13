package entity

type CommonResp struct {
	Type  string `json:"type"`   // 返回类型
	ReqId string `json:"req_id"` // 跟踪请求uuid
}

type TableSnapshotResp struct {
	CommonResp

	Data struct {
		Rows []map[string]interface{} `json:"rows"`
	} `json:"data"`
}

type TableDeltaResp struct {
	CommonResp

	Data struct {
		BlockNum int    `json:"block_num"`
		Step     string `json:"step"`
		DBOP     DBOp   `json:"dbop"`
	} `json:"data"`
}

type TransactionLifecycleResp struct {
	Type string `json:"type"`
	Data struct {
		Lifecycle TransactionLifecycle `json:"lifecycle"`
	} `json:"data"`
}

type PingResp struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type ProgressResp struct {
	CommonResp

	Data struct {
		BlockNum int    `json:"block_num"`
		BlockId  string `json:"block_id"`
	} `json:"data"`
}

type UnListenResp struct {
	Type string `json:"type"`
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type ListeningResp struct {
	CommonResp

	Data struct {
		NextBlock int `json:"next_block"`
	} `json:"data"`
}

type ErrorResp struct {
	Type string `json:"type"`
	Data struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details struct {
			TxId string `json:"tx_id"`
		} `json:"details"`
	} `json:"data"`
}
