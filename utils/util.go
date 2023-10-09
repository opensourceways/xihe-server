package utils

import (
	"crypto/rand"
	"math/big"
	"os"
	"time"
	"unicode/utf8"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

const (
	layout     = "2006-01-02"
	timeLayout = "2006-01-02 15:04:05"
)

func LoadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func Now() int64 {
	return time.Now().Unix()
}

func ToDate(n int64) string {
	if n == 0 {
		n = Now()
	}

	return time.Unix(n, 0).Format(layout)
}

func Date() string {
	return time.Now().Format(layout)
}

func DateAndTime(n int64) (string, string) {
	if n <= 0 {
		return "", ""
	}

	t := time.Unix(n, 0)

	return t.Format(layout), t.Format(timeLayout)
}

func ToUnixTime(v string) (time.Time, error) {
	t, err := time.Parse(layout, v)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func ExpiryReduceSecond(expiry int64) time.Time {
	return time.Now().Add(time.Duration(expiry-10) * time.Second)
}

func Expiry(expiry int64) int64 {
	return time.Now().Add(time.Second * time.Duration(expiry)).Unix()
}

func IsExpiry(expiry int64) bool {
	if expiry <= 0 {
		return false
	}

	return time.Now().Unix() > expiry
}

func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}

func GenRandoms(max, total int) []int {
	i := 0
	m := make(map[int]struct{})
	r := make([]int, total)
	for {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(max+1)))
		if err != nil {
			logrus.Debugf("Error generating random number: %s", err)
			return nil
		}
		n := int(randInt.Int64())

		if _, ok := m[n]; !ok {
			m[n] = struct{}{}
			r[i] = n
			if i++; i == total {
				break
			}
		}
	}

	return r
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int) int {
	return a * (b / GCD(a, b))
}

func RetryThreeTimes(f func() error) {
	if err := f(); err == nil {
		return
	}

	for i := 1; i < 3; i++ {
		time.Sleep(time.Second)

		if err := f(); err == nil {
			return
		}
	}
}
