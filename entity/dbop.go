package entity

// database operation
type DBOp struct {
	Op          string `json:"op,omitempty"`      // 操作名称(REM:删除 UPD:更新 INS:插入)
	ActionIndex int    `json:"action_idx"`        //
	Account     string `json:"account,omitempty"` // 操作这个数据库的合约账户
	Table       string `json:"table,omitempty"`   // 表名
	Scope       string `json:"scope,omitempty"`   // scope
	Key         string `json:"key,omitempty"`     // 表的主键
	Old         *DBRow `json:"old,omitempty"`     // 旧的行数据(删除、更新时触发)
	New         *DBRow `json:"new,omitempty"`     // 新的行数据(插入、更新时触发)
}

// database row data
type DBRow struct {
	Payer string      `json:"payer,omitempty"` // 执行者
	Hex   string      `json:"hex,omitempty"`   // 二进制数据的十六进制编码的字符串
	JSON  interface{} `json:"json,omitempty"`  // 数据对象
	Error string      `json:"error,omitempty"` // 错误信息
}
