package entity

type CommonResp struct {
	Type  string `json:"type"`
	ReqId string `json:"req_id"`
}

// get_table_rows resp snapshot
type TableSnapshotResp struct {
	CommonResp

	Data struct {
		Rows []map[string]interface{} `json:"rows"`
	} `json:"data"`
}

// get_table_rows resp delta
type TableDeltaResp struct {
	CommonResp

	Data struct {
		BlockNum int    `json:"block_num"`
		Step     string `json:"step"`
		//DBOP     DBOp   `json:"dbop"`
	} `json:"data"`
}
