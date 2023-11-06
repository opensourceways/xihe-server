package utils

import "unsafe"

func ClearByteArrayMemory(bs ...[]byte) {
	f := func(b []byte) {
		for i := range b {
			b[i] = 0
		}
	}

	for i := range bs {
		f(bs[i])
	}
}

func ClearStringMemory(s ...string) {
	for i := range s {
		clearStringMemory(s[i])
	}
}

func clearStringMemory(s string) {
	if len(s) <= 1 {
		return
	}

	bs := *(*[]byte)(unsafe.Pointer(&s))
	for i := 0; i < len(bs); i++ {
		bs[i] = 0
	}
}
