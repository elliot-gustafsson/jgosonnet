package interner

import (
	"math/bits"
	"unsafe"
)

const (
	rapidSeed uint64 = 0xf4d8ef1732c4ed6e

	rapidSecret0 uint64 = 0x2d358dccaa6c78a5
	rapidSecret1 uint64 = 0x8bb84b93962eacc9
	rapidSecret2 uint64 = 0x4b33a62ed433d4a3
)

// wymum operation: 64x64 -> 128-bit multiply and XOR fold
func rapidMix(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}

func read64(p *byte, offset int) uint64 {
	return *(*uint64)(unsafe.Add(unsafe.Pointer(p), offset))
}

func read32(p *byte, offset int) uint32 {
	return *(*uint32)(unsafe.Add(unsafe.Pointer(p), offset))
}

func read8(p *byte, offset int) uint8 {
	return *(*uint8)(unsafe.Add(unsafe.Pointer(p), offset))
}

// HashString implements the rapidhash algorithm.
func HashString(s string) uint64 {
	l := len(s)
	seed := rapidSeed

	// init seed with mixing and length
	seed ^= rapidMix(seed^rapidSecret0, rapidSecret1) ^ uint64(l)

	if l == 0 {
		return rapidMix(seed, rapidSecret2)
	}

	p := unsafe.StringData(s)
	var a, b uint64

	// 16 byte or smaller strings
	if l <= 16 {
		// split strings in left and right then mix them
		if l >= 8 {
			a = read64(p, 0)
			b = read64(p, l-8)
		} else if l >= 4 {
			a = uint64(read32(p, 0))
			b = uint64(read32(p, l-4))
		} else {
			a = uint64(read8(p, 0)) << 16
			a |= uint64(read8(p, l>>1)) << 8
			a |= uint64(read8(p, l-1))
		}
		return rapidMix(seed^a, rapidSecret2^b)
	}

	i := 0
	if l > 96 {
		// reads 96 bytes at a time and splits into 6 chunks continously and mixes n with n+3 to 3 results, then xor those.
		see1 := seed
		see2 := seed
		for ; i < l-96; i += 96 {
			seed = rapidMix(read64(p, i)^rapidSecret0, read64(p, i+8)^seed)
			see1 = rapidMix(read64(p, i+16)^rapidSecret1, read64(p, i+24)^see1)
			see2 = rapidMix(read64(p, i+32)^rapidSecret2, read64(p, i+40)^see2)
			seed = rapidMix(read64(p, i+48)^rapidSecret0, read64(p, i+56)^seed)
			see1 = rapidMix(read64(p, i+64)^rapidSecret1, read64(p, i+72)^see1)
			see2 = rapidMix(read64(p, i+80)^rapidSecret2, read64(p, i+88)^see2)
		}
		seed ^= see1 ^ see2
	}

	for ; i < l-16; i += 16 {
		// consume 16 bytes at a time, first 8 bytes mixed with next 8
		seed = rapidMix(read64(p, i)^rapidSecret2, read64(p, i+8)^seed^rapidSecret1)
	}

	// handle the remaning bytes
	a = read64(p, l-16)
	b = read64(p, l-8)

	return rapidMix(seed^a, rapidSecret2^b)
}
