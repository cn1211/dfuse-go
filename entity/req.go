package entity

// common req
type CommonReq struct {
	Type  string `json:"type"`   // 请求类型
	ReqId string `json:"req_id"` // 跟踪请求uuid
	*OptionReq
}

// request option
type OptionReq struct {
	Fetch            bool `json:"fetch"`             // 是否捕捉快照(仅支持get_table_rows,get_transaction)
	Listen           bool `json:"listen"`            // 是否监听请求流后续的改变
	StartBlock       int  `json:"start_block"`       // 开始监听的区块高度(0代表当前区块 负数:例-200表示当前区块往前200个区块 正数:例200代表从整条eos链的第200个区块开始输送(需要额外的key))
	IrreversibleOnly bool `json:"irreversible_only"` // 是否只输出不可逆区块
	WithProgress     int  `json:"with_progress"`     // 区块输出之间的间隔
}

type TableRowsReq struct {
	CommonReq

	Data GetTableRows `json:"data"`
}

type GetTableRows struct {
	Code          string `json:"code"`           // 合约
	Scope         string `json:"scope"`          // 交易对
	Table         string `json:"table"`          // 表名
	TableKey      string `json:"table_key"`      // 键名称，字符串，可选
	Json          bool   `json:"json"`           // 是否返回JSON格式的结果，布尔型，默认值：true
	LowerBound    string `json:"lower_bound"`    // 结果数据应当满足的下界值，字符串
	UpperBound    string `json:"upper_bound"`    // 结果数据应当满足的上界值，字符串
	Limit         int    `json:"limit"`          // 返回数量上限 默认10条
	IndexPosition int    `json:"index_position"` // 使用的索引序号,例如，主键索引为1，次级索引为2，字符串，默认值:1
	KeyType       int    `json:"key_type"`       // 索引键类型，例如uint64_t或name，字符串
	EncodeType    string `json:"encode_type"`    // 编码类型字符串默认:dec
}

type TransactionLifecycleReq struct {
	CommonReq
	Data struct {
		Id string `json:"id"` // tx-hash
	} `json:"data"`
}

type UnListenReq struct {
	Type string `json:"type"`
	Data struct {
		ReqId string `json:"req_id"` // 请求id
	} `json:"data"`
}
