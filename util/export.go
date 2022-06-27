package util

import (
	"bytes"
	"strconv"
)

const (
	//CodeStatusError for server error
	CodeStatusServerError int = 500
	//CodeStatusNormal for client error
	CodeStatusClientError int = 400
	//CodeStatusNormal normal statu
	CodeStatusNormal int = 200
)

//JsonData export to clent
type JsonData struct {
	Code   int         `json:"code"`
	Title  interface{} `json:"title"`
	Attach interface{} `json:"attach,omitempty"`
	Data   interface{} `json:"data"`
	Other  interface{} `json:"other,omitempty"`
}

//ExportData ExportData
func ExportData(code int, title interface{}, data ...interface{}) *JsonData {
	var resultData JsonData
	resultData.Code = code
	resultData.Title = title
	resultData.Data = data[0]

	if len(data) > 1 {
		resultData.Attach = data[1]
		if len(data) > 2 {
			resultData.Other = data[2]
		}

	}
	if code > 300 {
		// if GetConfig().AppModel == "release" {
		// 	Log.Errorf("code:%d,title:%s,data:%v", resultData.Code, resultData.Title, resultData.Data)
		// 	resultData.Title = "Error Information"
		// 	resultData.Data = ""
		// }
		if err, ok := data[0].(error); ok {
			Log.Warnln(err.Error())
		} else {
			Log.Warnln(data)
		}
	}
	return &resultData
}

//StringsToJSON StringsToJSON
func StringsToJSON(str string) string {
	var jsons bytes.Buffer
	for _, r := range str {
		rint := int(r)
		if rint < 128 {
			jsons.WriteRune(r)
		} else {
			jsons.WriteString("\\u")
			jsons.WriteString(strconv.FormatInt(int64(rint), 16))
		}
	}
	return jsons.String()
}
