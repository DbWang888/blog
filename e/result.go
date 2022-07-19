package e

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func GetSucessResult(data interface{}) Result {
	return Result{
		Code: SUCCESS,
		Msg:  GetMsg(SUCCESS),
		Data: data,
	}
}

func GetErrResult(code int, err error) Result {
	return Result{
		Code: code,
		Msg:  GetMsg(code),
		Data: err.Error(),
	}
}
