package arena

import "unsafe"

// Memclr zeroes the memory behind a typed pointer as raw bytes.
func Memclr[T any](ptr *T) {
	size := unsafe.Sizeof(*ptr)
	if size > 0 {
		clear(unsafe.Slice((*byte)(unsafe.Pointer(ptr)), size))
	}
}

// MemclrSlice zeroes the memory behind a slice pointer as raw bytes.
func MemclrSlice[T any](s []T) {
	var zero T
	elemSize := unsafe.Sizeof(zero)
	if elemSize > 0 && len(s) > 0 {
		totalSize := elemSize * uintptr(len(s))
		clear(unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(s))), totalSize))
	}
}
