package controller

const (
	errorBadRequestBody = "bad_request_body"
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
