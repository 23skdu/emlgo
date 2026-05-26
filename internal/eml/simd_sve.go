//go:build arm64 && linux
// +build arm64,linux

package eml

import (
	"encoding/binary"
	"os"
	"syscall"
)

const (
	sve_AT_HWCAP    = 16
	sve_HWCAP_SVE   = 1 << 22
	sve_PR_GET_VL   = 51
	sve_VL_LEN_MASK = 0xffff
	sve_DEFAULT_VL  = 32
)

func detectSVE() {
	hasSVE = false
	sveVectorLength = 0

	auxv, err := os.ReadFile("/proc/self/auxv")
	if err != nil {
		return
	}

	for i := 0; i+16 <= len(auxv); i += 16 {
		typ := binary.LittleEndian.Uint64(auxv[i:])
		val := binary.LittleEndian.Uint64(auxv[i+8:])
		if typ == sve_AT_HWCAP {
			if val&sve_HWCAP_SVE != 0 {
				hasSVE = true
				sveVectorLength = detectSVEVL()
			}
			return
		}
	}
}

func detectSVEVL() int {
	vl, _, err := syscall.Syscall6(
		syscall.SYS_PRCTL, sve_PR_GET_VL, 0, 0, 0, 0, 0,
	)
	if err == 0 && vl != 0 {
		return int(vl) & sve_VL_LEN_MASK
	}

	data, err2 := os.ReadFile("/proc/sys/abi/sve_default_vector_length")
	if err2 == nil {
		var v int
		for _, b := range data {
			if b >= '0' && b <= '9' {
				v = v*10 + int(b-'0')
			} else if v > 0 {
				break
			}
		}
		if v > 0 {
			return v
		}
	}

	return sve_DEFAULT_VL
}
