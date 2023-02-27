package constans

// 错误码
const (
	EMPTY            = ""
	OK               = 0
	Success          = 0
	ErrNormal        = 1
	ErrTokenNotFound = 2
	ErrNoProp        = 3
	//ErrOpTooFast 某请求操作过快
	ErrOpTooFast = 4
	// 禁止
	ErrForbidden       = 403
	ErrNoAuth          = "您没有权限"
	ErrTxtTooFast      = "操作太快啦"
	ErrTxtTooFastShow  = "操作太快啦，请稍后重试"
	ErrTxtParamErr     = "参数有误"
	ErrTxtAccountErr   = "账号异常"
	ErrTxtCreate       = "创建数据失败"
	ErrServerFix       = "维护中，请稍后重试"
	ErrSMSFailed       = "短信发送失败，请稍后再试"
	ErrSMSInvalid      = "验证码已过期，请重新发送"
	ErrSMSCodeErr      = "验证码有误"
	ErrOldOrderExpired = "订单已过期"
	ErrOrderNotExist   = "订单不存在"
	ErrCancelError     = "订单取消失败"
	ErrVerifyTime      = "请检查设备系统时间"
	// 繁忙
	ErrBusyCode = 502
	ErrBusyTxt  = "服务器繁忙，请稍后重试"
	ServerPanic = "系统错误0x01"
)
