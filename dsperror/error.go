package dsperror

/*
错误代码总汇
*/
const (
	ErrorUnauthorized = 10000001 //错误代码规则 最右边三位是错误代码，右边数第四，五，六位是模块代号，剩下二位为项目代码
	ErrorSyncIndex    = 10000002 //同步索引错误
	ErrorStopBid      = 10000003 //停止竞价失败
	ErrorUnknown      = 11111111 //未知错误
)

var errorText = map[int]string{
	ErrorUnauthorized: "Unauthorized",
	ErrorSyncIndex:    "NoneIndex",
	ErrorStopBid:      "NotStopBid",
	ErrorUnknown:      "UnknownError",
}

func ErrorText(code int) string {
	return errorText[code]
}
