package utils

import (
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	"unicode/utf8"

	"sigs.k8s.io/yaml"
)

func LoadFromYaml(path string, cfg interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := []byte(os.ExpandEnv(string(b)))

	return yaml.Unmarshal(content, cfg)
}

func Now() int64 {
	return time.Now().Unix()
}

func ToDate(n int64) string {
	if n == 0 {
		n = Now()
	}

	return time.Unix(n, 0).Format("2006-01-02")
}

func Date() string {
	return time.Now().Format("2006-01-02")
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
	// set seed
	rand.Seed(time.Now().UnixNano())

	i := 0
	m := make(map[int]struct{})
	r := make([]int, total)
	for {
		n := rand.Intn(max) + 1

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
