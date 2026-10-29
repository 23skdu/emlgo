//go:build amd64 && !purego

package eml

func amd64ExpSIMD(x, result []float64) {
	n := len(x)
	simdLen := 0
	if hasAVX512 {
		simdLen = (n / 8) * 8
		expAVX512(x[:simdLen], result[:simdLen])
	} else if hasAVX2 {
		simdLen = (n / 4) * 4
		expAVX2(x[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = nativeExp(x[i])
	}
}

func amd64LogSIMD(x, result []float64) {
	n := len(x)
	simdLen := 0
	if hasAVX512 {
		simdLen = (n / 8) * 8
		logAVX512(x[:simdLen], result[:simdLen])
	} else if hasAVX2 {
		simdLen = (n / 4) * 4
		logAVX2(x[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = logScalar(x[i])
	}
}

func amd64SinSIMD(x, result []float64) {
	n := len(x)
	simdLen := 0
	if hasAVX512 {
		simdLen = (n / 8) * 8
		sinAVX512(x[:simdLen], result[:simdLen])
	} else if hasAVX2 {
		simdLen = (n / 4) * 4
		sinAVX2(x[:simdLen], result[:simdLen])
	}
	for i := 0; i < simdLen; i++ {
		if result[i] != result[i] {
			result[i] = nativeSin(x[i])
		}
	}
	for i := simdLen; i < n; i++ {
		result[i] = nativeSin(x[i])
	}
}

func amd64CosSIMD(x, result []float64) {
	n := len(x)
	simdLen := 0
	if hasAVX512 {
		simdLen = (n / 8) * 8
		cosAVX512(x[:simdLen], result[:simdLen])
	} else if hasAVX2 {
		simdLen = (n / 4) * 4
		cosAVX2(x[:simdLen], result[:simdLen])
	}
	for i := 0; i < simdLen; i++ {
		if result[i] != result[i] {
			result[i] = nativeCos(x[i])
		}
	}
	for i := simdLen; i < n; i++ {
		result[i] = nativeCos(x[i])
	}
}

func amd64TanSIMD(x, result []float64) {
	n := len(x)
	simdLen := 0
	if hasAVX512 {
		simdLen = (n / 8) * 8
		tanAVX512(x[:simdLen], result[:simdLen])
	} else if hasAVX2 {
		simdLen = (n / 4) * 4
		tanAVX2(x[:simdLen], result[:simdLen])
	}
	for i := 0; i < simdLen; i++ {
		if result[i] != result[i] {
			result[i] = nativeTan(x[i])
		}
	}
	for i := simdLen; i < n; i++ {
		result[i] = nativeTan(x[i])
	}
}

func dispatchExpSIMDTo(x, result []float64) {
	ForEachChunk(len(x), func(start, end int) {
		amd64ExpSIMD(x[start:end], result[start:end])
	})
}

func dispatchLogSIMDTo(x, result []float64) {
	ForEachChunk(len(x), func(start, end int) {
		amd64LogSIMD(x[start:end], result[start:end])
	})
}

func dispatchSinSIMDTo(x, result []float64) {
	ForEachChunk(len(x), func(start, end int) {
		amd64SinSIMD(x[start:end], result[start:end])
	})
}

func dispatchCosSIMDTo(x, result []float64) {
	ForEachChunk(len(x), func(start, end int) {
		amd64CosSIMD(x[start:end], result[start:end])
	})
}

func dispatchTanSIMDTo(x, result []float64) {
	ForEachChunk(len(x), func(start, end int) {
		amd64TanSIMD(x[start:end], result[start:end])
	})
}
