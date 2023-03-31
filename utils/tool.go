package utils

import "math/bits"

// BitLength 只适用于2的幂,输入4返回2
func BitLength(n uint32) uint32 {
	if n == 0 {
		return 0
	}
	return uint32(bits.Len(uint(n) - 1))
}

func IsPowerOfTwo(n uint32) bool {
	return n&(n-1) == 0
}

func NextPowerOf2(val uint32) uint32 {
	if IsPowerOfTwo(val) {
		if val == 0 {
			return 1
		} else {
			return val
		}
	} else {
		val |= val >> 1
		val |= val >> 2
		val |= val >> 4
		val |= val >> 8
		val |= val >> 16
		return val + 1
	}

}
