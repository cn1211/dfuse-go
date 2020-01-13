package dfuse

// request type
const (
	GetActionTraces = "get_action_traces"
	GetTransaction  = "get_transaction"
	GetTableRows    = "get_table_rows"
	GetHeadInfo     = "get_head_info"
)

// response type
const (
	ActionTraces         = "action_trace"
	TransactionLifecycle = "transaction_lifecycle"
	TableSnapshot        = "table_snapshot"
	TableDelta           = "table_delta"
	HeadInfo             = "head_info"
)

// common request type
const (
	UnListen = "unlisten"
	Pong     = "pong"
)

// common response type
const (
	UnListened = "unlistened"
	Listening  = "listening"
	Progress   = "progress"
	Error      = "error"
	Ping       = "ping"
)
