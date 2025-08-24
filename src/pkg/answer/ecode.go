package answer

const (
	EcodeOK                        = "Ledger.0000" // 结果正常
	EcodeError                     = "Ledger.0101" // 结果错误
	EcodeQueryError                = "Ledger.0102" // 查询参数错误
	EcodeInvalidRequestParamsError = "Ledger.0103" // 请求参数校验失败。
	EcodeInvalidRequestError       = "Ledger.0104" // 请求体错误

	EcodeInvalidTokenError      = "Ledger.0175" // 无效，错误的 token
	EcodeNoActionError          = "Ledger.0177" // action 不存在
	EcodePolicyNotAuthorized    = "Ledger.0178" // 策略未授权此操作
	EcodeDeleteResourceConflict = "Ledger.0198" // 409 删除资源时冲突

	EcodeUpstreamResponseError = "Ledger.0301" // 上游返回处理错误
)
