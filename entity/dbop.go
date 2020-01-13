package entity

// database operation
type DBOp struct {
	Op          string `json:"op,omitempty"`      // 操作名称(REM:删除 UPD:更新 INS:插入)
	ActionIndex int    `json:"action_idx"`        //
	Account     string `json:"account,omitempty"` // 操作这个数据库的合约账户
	Table       string `json:"table,omitempty"`   // 表名
	Scope       string `json:"scope,omitempty"`   // scope
	Key         string `json:"key,omitempty"`     // 表的主键
	Old         *DBRow `json:"old,omitempty"`     //
	New         *DBRow `json:"new,omitempty"`     //
}

// database row data
type DBRow struct {
	Payer string      `json:"payer,omitempty"`
	Hex   string      `json:"hex,omitempty"`
	JSON  interface{} `json:"json,omitempty"`
	Error string      `json:"error,omitempty"`
}
