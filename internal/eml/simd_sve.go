//go:build arm64 && linux
// +build arm64,linux

package eml

import (
	"encoding/binary"
	"os"
)

const (
	sve_AT_HWCAP  = 16
	sve_HWCAP_SVE = 1 << 22
	hwcap_ASIMD   = 1 << 1
	hwcap_ASIMDDP = 1 << 20
)

var (
	sveVectorLength int
	hwHasSVE        bool
	hwHasNeonDot    bool
)

func detectSVE() {
	// Capabilities HasSVE() and HasNeonDot() must remain honest (false) until
	// real SVE or Dot Product assembly instructions are provided in this library.
	hasSVE = false
	hasNeonDot = false
	sveVectorLength = 0

	auxv, err := os.ReadFile("/proc/self/auxv")
	if err != nil {
		return
	}

	for i := 0; i+16 <= len(auxv); i += 16 {
		typ := binary.LittleEndian.Uint64(auxv[i:])
		val := binary.LittleEndian.Uint64(auxv[i+8:])
		if typ == sve_AT_HWCAP {
			if val&hwcap_ASIMD != 0 {
				hasNeon = true
				hasFMA = true
			}
			if val&hwcap_ASIMDDP != 0 {
				hwHasNeonDot = true
			}
			if val&sve_HWCAP_SVE != 0 {
				hwHasSVE = true
			}
			return
		}
	}
}
