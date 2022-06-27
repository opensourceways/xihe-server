package util

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	"unicode"

	validator "github.com/go-playground/validator/v10"
	"k8s.io/client-go/rest"
)

var GlobK8sConfig *rest.Config
var RequiredValidator validator.Validate

func Catchs() {
	if err := recover(); err != nil {
		Log.Error("The program is abnormal, err: ", err)
	}
}

const DATE_FORMAT = "2006-01-02 15:04:05"
const DATE_T_FORMAT = "2006-01-02T15:04:05"
const DATE_T_Z_FORMAT = "2006-01-02T15:04:05Z"
const DT_FORMAT = "2006-01-02"

func GetCurDate() string {
	return time.Now().In(CnTime).Format(DT_FORMAT)
}

func GetCurTime() string {
	return time.Now().In(CnTime).Format(DATE_FORMAT)
}

func createDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dir, 0777)
		}
	}
	return err
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func TimeConverStr(ts, oldLayout, newLayout string) string {
	if ts == "" || oldLayout == "" || newLayout == "" {
		return ""
	}
	timeStr := ts
	if timeStr != "" && len(timeStr) > 19 {
		timeStr = timeStr[:19]
	}
	unixTime := int64(0)
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(oldLayout, timeStr, loc)
	if err == nil {
		unixTime = theTime.Unix() + 8*3600
	} else {
		Log.Error(err)
		return ""
	}
	timx := time.Unix(unixTime, 0).Format(newLayout)
	return timx
}

func TimeTConverStr(ts string) string {
	if len(ts) > 19 {
		ts = ts[:19]
	}
	return TimeConverStr(ts, DATE_T_FORMAT, DATE_FORMAT)
}

func TimeStrToInt(ts, layout string) int64 {
	if ts == "" {
		return 0
	}
	if layout == "" {
		layout = DATE_FORMAT
	}
	timeStr := ts
	if timeStr != "" && len(timeStr) > 19 {
		timeStr = timeStr[:19]
	}
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(layout, timeStr, loc)
	if err == nil {
		unixTime := theTime.Unix()
		return unixTime
	} else {
		Log.Error(err)
	}
	return 0
}

// Time string to timestamp
func PraseTimeInt(stringTime string) int64 {
	return TimeStrToInt(stringTime, DATE_FORMAT)
}

func PraseTimeTint(tsStr string) int64 {
	return TimeStrToInt(tsStr, DATE_T_FORMAT)
}

func LocalTimeToUTC(strTime string) time.Time {
	local, _ := time.ParseInLocation(DATE_FORMAT, strTime, time.Local)
	return local
}

func GetTZHTime(hours time.Duration) string {
	now := time.Now().In(CnTime)
	h, _ := time.ParseDuration("-1h")
	dateTime := now.Add(hours * h).Format(DATE_T_Z_FORMAT)
	fmt.Println("dateTime: ", dateTime)
	return dateTime
}

func DelFile(fileList []string) {
	if len(fileList) > 0 {
		for _, filex := range fileList {
			if FileExists(filex) {
				err := os.Remove(filex)
				if err != nil {
					Log.Error(err)
				}
			}
		}
	}
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func GetRandomString(l int) string {
	str := "abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().Local().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func IsLetter(chars rune) bool {
	return unicode.IsLetter(chars)
}
