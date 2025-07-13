package answer

const (
	EcodeOK                        = "DTWM.0000" // 结果正常
	EcodeError                     = "DTWM.0101" // 结果错误
	EcodeQueryError                = "DTWM.0102" // 查询参数错误
	EcodeInvalidRequestParamsError = "DTWM.0103" // 请求参数校验失败。
	EcodeInvalidRequestError       = "DTWM.0104" // 请求体错误

	EcodeInvalidTokenError      = "DTWM.0175" // 无效，错误的 token
	EcodeNoActionError          = "DTWM.0177" // action 不存在
	EcodePolicyNotAuthorized    = "DTWM.0178" // 策略未授权此操作
	EcodeDeleteResourceConflict = "DTWM.0198" // 409 删除资源时冲突

	EcodeUpstreamResponseError = "DTWM.0301" // 上游返回处理错误
)
