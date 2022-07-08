package controller

const (
	errorBadRequestBody  = "bad_request_body"
	errorBadRequestParam = "bad_request_param"
	errorSystemError     = "system_error"
)

// responseData is the response data to client
type responseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newResponse(code, msg string, data interface{}) responseData {
	return responseData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func newResponseData(data interface{}) responseData {
	return responseData{
		Data: data,
	}
}

func newResponseError(code string, err error) responseData {
	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseMsg(code, msg string) responseData {
	return responseData{
		Code: code,
		Msg:  msg,
	}
}
