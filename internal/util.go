package internal

import (
	"reflect"
	"unsafe"
)

func IntToByteArray(i int32) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Len = 8
	sh.Cap = 8
	sh.Data = uintptr(unsafe.Pointer(&i))

	return b[:]
}

func ByteArrayToInt(b []byte) int32 {
	return *(*int32)(unsafe.Pointer(&b[0]))
}

func ItoBSlice(i int) []byte {
	data := *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
	return data[:]
}
