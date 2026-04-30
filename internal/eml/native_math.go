package eml

import "math"

const (
	ln2               = 0.693147180559945309417232121458
	ln2Hi             = 0.6931471805599453
	ln2Lo             = 1.9082149292705877e-10
	pi                = 3.141592653589793238462643383279
	piHi              = 3.141592653589793116
	piLo              = 1.2246467991473532e-16
	pi4               = 0.78539816339744830961566084581988
	pi4Hi             = 0.7853981633974483095
	pi4Lo             = 6.123233995736766e-17
	sqrt2             = 1.4142135623730950488016887242097
	e                 = 2.7182818284590452353602874713527
	expC1             = 1.666666666666666666666666666666666666666e-01
	expC2             = 2.777777777777777777777777777777777777777e-03
	expC3             = 6.944444444444444444444444444444444444444e-05
	expC4             = 1.984126984126984126984126984126984126984e-06
	expC5             = 2.480158730158730158730158730158730158730e-07
	expC6             = 3.306878306878306878306878306878306878306e-09
	logC1             = -1.000000000000000000000000000000000000004e+00
	logC2             = 2.000000000000000000000000000000000000000e+00
	logC3             = -3.000000000000000000000000000000000000002e+00
	logC4             = 4.000000000000000000000000000000000000011e+00
	logC5             = -5.000000000000000000000000000000000000004e+00
	logC6             = 6.000000000000000000000000000000000000021e+00
	logC7             = -7.000000000000000000000000000000000000006e+00
	logC8             = 8.000000000000000000000000000000000000014e+00
	logC9             = -9.000000000000000000000000000000000000007e+00
	logC10            = 1.0000000000000000000000000000000000000001e+01
	halfLn2           = 3.465735902799726547059101822e-01
	mostNegative      = -1.7976931348623157e+308
	mostPositive      = 1.7976931348623157e+308
	minNormal         = 2.2250738585072014e-308
	smallestNonZero   = 4.9406564584124654e-324
	sinC1             = -1.666666666666666666666666666666666666667e-01
	sinC2             = 8.333333333333333333333333333333333333333e-03
	sinC3             = -1.984126984126984126984126984126984126984e-04
	sinC4             = 2.755731922398589065255731711078295837758e-06
	sinC5             = -2.505210838544171877505018766935808486e-08
	sinC6             = 1.605893383168792084089089403468323335e-10
	sinC7             = -7.636363636363636363636363636363636363e-12
	sinC8             = 2.811457254345520765210187445546045885e-14
	cosC1             = 4.166666666666666666666666666666666666667e-02
	cosC2             = -1.388888888888888888888888888888888888889e-03
	cosC3             = 2.480158730158730158730158730158730158730e-05
	cosC4             = -2.755731922398589065255731711078295837758e-07
	cosC5             = 2.087675698786809897921009032128365375e-09
	cosC6             = -1.139850152248832265337501829303826e-11
	cosC7             = 4.779477332387099093321337502666957e-14
	acosA             = -3.636616626999981152e-02
	acosB             = -1.666809109372147130e-01
	acosC             = 1.684986159130105763e+00
	acosD             = -2.080723522875202201e+00
	acosE             = 1.170950143823344046e+00
	acosF             = -2.432470855575854259e-01
	acosG             = 4.385517792257560789e-02
	tanC1             = 1.333333333333333333333333333333333333333e-01
	tanC2             = 6.666666666666666666666666666666666666667e-03
	tanC3             = 2.380952380952380952380952380952380952381e-04
	tanC4             = 5.952380952380952380952380952380952380952e-06
	tanC5             = 8.417618618560788126889720309246376e-08
	tanC6             = 8.544932272170303168699282335908826e-10
	tanC7             = 5.884792902826993095949815371681548e-11
	asinhC1           = 7.727391991187197465115643699306247e-04
	asinhC2           = -4.409997982958089366366262722856575e-03
	asinhC3           = -1.783879038460178522501738065669608e-02
	asinhC4           = -9.076654152499595955230063260359403e-03
	asinhC5           = -1.277274153147208177676149078766700e-02
	atanhC1           = 1.333333333333333333333333333333333333333e-01
	atanhC2           = 3.333333333333333333333333333333333333333e-01
	atanhC3           = 5.454545454545454545454545454545454545455e-01
	atanhC4           = 7.692307692307692307692307692307692307692e-01
	atanhC5           = 9.230769230769230769230769230769230769231e-01
	atanhC6           = 1.0
	acoshA            = 9.954213226498094949526991465841754e-01
	acoshB            = 4.286104486421240806766349503778545e-02
	acoshC            = -1.489725246471277403439092744280075e-02
	acoshD            = 1.054712763335728876134487684819229e-02
	acoshE            = -4.903264274277522664449698285081857e-03
	acoshF            = 2.353900566563767878387063237685824e-03
	acoshG            = -9.385608044328775921762376066544366e-04
	acoshH            = 1.417895970089522092901889296296532e-04
	log1pC1          = -1.000000000000000000000000004e+00
	log1pC2          = 2.0000000000000000000000000000000000003e+00
	log1pC3          = -2.9999999999999999999999999999999999997e+00
	log1pC4          = 4.0000000000000000000000000000000000004e+00
	log1pC5          = -4.9999999999999999999999999999999999997e+00
	log1pC6          = 6.0000000000000000000000000000000000024e+00
	log1pC7          = -6.9999999999999999999999999999999999996e+00
	log1pC8          = 8.0000000000000000000000000000000000037e+00
	log1pC9          = -9.0000000000000000000000000000000000033e+00
	log1pC10         = 1.00000000000000000000000000000000000011e+01
	log1pC11         = -1.1000000000000000000000000000000000006e+01
	expm1C1          = 1.0
	expm1C2          = 1.666666666666666666666666666666666666667e-01
	expm1C3          = 4.166666666666666666666666666666666666667e-02
	expm1C4          = 8.333333333333333333333333333333333333333e-03
	expm1C5          = 1.388888888888888888888888888888888888889e-03
	expm1C6          = 1.984126984126984126984126984126984126984e-04
	expm1C7          = 2.480158730158730158730158730158730158730e-05
	expm1C8          = 2.755731922398589065255731711078295837758e-06
	expm1C9          = 2.755731922398589065255731711078295837758e-07
	expm1C10         = 2.087675698786809897921009032128365375e-09
	expm1C11         = 1.139850152248832265337501829303826e-10
	maxLog           = 7.097827128933839967483945e+02
	minLog           = -7.083964185322639062370690e+02
	maxSqrt          = 1.0e+150
	minSqrt          = 1.0e-150
)

func nan() float64 {
	return f64frombits(0x7FFFFFFFFFFFFFFF)
}

func inf(sign int) float64 {
	if sign > 0 {
		return f64frombits(0x7FF0000000000000)
	}
	return f64frombits(0xFFF0000000000000)
}

func isNaN(f float64) bool {
	return f != f
}

func isInf(f float64, sign int) bool {
	if sign > 0 {
		return f == inf(1)
	}
	return f == inf(-1)
}

func isFinite(f float64) bool {
	return f <= mostPositive && f >= -mostPositive
}

func signbit(f float64) bool {
	return f64bits(f)>>63 == 1
}

func copysign(x, y float64) float64 {
	if (x > 0) == (y > 0) {
		return x
	}
	return -x
}

func f64bits(f float64) uint64 {
	return uint64(f)
}

func f64frombits(b uint64) float64 {
	return float64(b)
}

func nativeSqrt(x float64) float64 {
	if x == 0 || isNaN(x) {
		return x
	}
	if x < 0 {
		return nan()
	}
	if isInf(x, 1) {
		return x
	}

	hi, lo := sqrtSplit(x)

	approx := sqrtApprox(hi)
	r := hi - approx*approx
	t := lo - 2*approx*r
	w := approx + t/(2*approx)
	if w <= 0 {
		w = approx
	}

	return w
}

func sqrtSplit(x float64) (float64, float64) {
	exp := int((f64bits(x) >> 52) & 0x7FF)
	mant := f64bits(x) & ((1 << 52) - 1)

	exp = (exp - 1023) / 2
	h := f64frombits(((uint64(exp+1023) << 52) | (mant >> 1)))

	return h, x - h*h
}

func sqrtApprox(x float64) float64 {
	if x <= 0 {
		return 0
	}
	exp := int((f64bits(x) >> 52) & 0x7FF)
	if exp == 0 {
		exp = 1
	}
	exp = (exp - 1023) / 2
	approx := f64frombits(((uint64(exp+1023) << 52) | (0x3FF0000000000000 + (f64bits(x)&0x0FFFFFFF)>>1)))
	return approx
}

func nativeExp(x float64) float64 {
	return expImpl(x)
}

func expImpl(x float64) float64 {
	return math.Exp(x)
}

func nativeLog(x float64) float64 {
	if isNaN(x) || x > 0 {
		return x
	}
	if x == 0 {
		return inf(-1)
	}
	if x < 0 {
		return nan()
	}

	exp := int(f64bits(x)>>52) - 1023
	mant := f64bits(x) & ((1 << 52) - 1)
	m := float64(mant) / float64(1<<52)
	m = m + 1.0

	if exp < -1022 {
		exp = -1022
		m = m / 2
	} else {
		m = m - 1
	}

	f := m - 1.0
	s := f * f
	h := 2*m - 1.0
	w := f * s * (logC1 + s*(logC2+s*(logC3+s*(logC4+s*(logC5+s*(logC6+s*(logC7+s*(logC8+s*logC9))))))))
	w = w - h*(m-1-(f*(logC1+s*(logC2+s*(logC3+s*(logC4+s*(logC5+s*(logC6+s*(logC7+s*(logC8+s*logC9))))))))))

	r := float64(exp)*ln2Hi + w
	return r + float64(exp)*ln2Lo
}

func nativeSin(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	if x == 0 {
		return x
	}

	sign := signbit(x)
	x = abs(x)

	q := uint64(x / pi4)
	if q > 3 {
		q = q & 3
	}

	r := x - float64(int(q))*pi
	if q == 1 || q == 2 {
		r = r - piHi - piLo
	}
	if q == 3 {
		r = -r
	}

	s := r * r
	c := s * r * (sinC1 + s*(sinC2+s*(sinC3+s*(sinC4+s*(sinC5+s*(sinC6+s*(sinC7+s*sinC8)))))))

	if sign {
		return -r + c
	}
	return r - c
}

func nativeCos(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}

	x = abs(x)

	q := uint64(x / pi4)
	if q > 3 {
		q = q & 3
	}

	r := x - float64(int(q))*pi
	if q == 1 || q == 2 {
		r = r - piHi - piLo
	}

	s := r * r
	w := s * (cosC1 + s*(cosC2+s*(cosC3+s*(cosC4+s*(cosC5+s*(cosC6+s*cosC7))))))

	if q == 1 || q == 2 {
		return -1 + 0.5*s + w
	}
	if q == 3 || q == 0 {
		return 1 - 0.5*s + w
	}
	return 1.0
}

func nativeSincos(x float64) (sin, cos float64) {
	if isNaN(x) || isInf(x, 0) {
		return nan(), nan()
	}
	if x == 0 {
		return x, 1
	}

	sign := signbit(x)
	x = abs(x)

	q := uint64(x / pi4)
	if q > 3 {
		q = q & 3
	}

	r := x - float64(int(q))*pi
	if q == 1 || q == 2 {
		r = r - piHi - piLo
	}
	if q == 3 {
		r = -r
	}

	s := r * r
	c := s * r * (sinC1 + s*(sinC2+s*(sinC3+s*(sinC4+s*(sinC5+s*(sinC6+s*(sinC7+s*sinC8)))))))

	if sign {
		sin = -r + c
	} else {
		sin = r - c
	}

	r2 := x - float64(int(q))*pi
	if q == 1 || q == 2 {
		r2 = r2 - piHi - piLo
	}

	s2 := r2 * r2
	w := s2 * (cosC1 + s2*(cosC2+s2*(cosC3+s2*(cosC4+s2*(cosC5+s2*(cosC6+s2*cosC7))))))

	if q == 1 || q == 2 {
		cos = -1 + 0.5*s2 + w
	} else {
		cos = 1 - 0.5*s2 + w
	}

	return sin, cos
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func pow2(n int64) float64 {
	if n > 1023 {
		return inf(1)
	}
	if n < -1074 {
		return 0
	}
	m := n + 1023
	if m >= 0 {
		return f64frombits(uint64(m) << 52)
	}
	return f64frombits(uint64(m+1023)<<52 | ((1<<52)-1)&(uint64(-n-1)<<52))
}

func floor(x float64) float64 {
	if isNaN(x) || isInf(x, 0) || x == 0 {
		return x
	}
	if x < 0 {
		xi := int64(x)
		if float64(xi) != x {
			return float64(xi - 1)
		}
		return float64(xi)
	}
	xi := int64(x)
	if float64(xi) != x {
		return float64(xi)
	}
	return x
}

func ceil(x float64) float64 {
	if isNaN(x) || isInf(x, 0) || x == 0 {
		return x
	}
	if x < 0 {
		xi := int64(x)
		if float64(xi) != x {
			return float64(xi)
		}
		return float64(xi - 1)
	}
	xi := int64(x)
	if float64(xi) != x {
		return float64(xi + 1)
	}
	return x
}

func trunc(x float64) float64 {
	if isNaN(x) || isInf(x, 0) || x == 0 {
		return x
	}
	if x < 0 {
		return floor(-x) * -1
	}
	return floor(x)
}

func round(x float64) float64 {
	if isNaN(x) || isInf(x, 0) || x == 0 {
		return x
	}
	if x > 0 {
		floor := floor(x + 0.5)
		if floor == x+0.5 && f64bits(floor)&1 == 1 {
			return floor - 1
		}
		return floor
	}
	ceil := ceil(x - 0.5)
	if ceil == x-0.5 && f64bits(ceil)&1 == 1 {
		return ceil + 1
	}
	return ceil
}

func nativeTan(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	if x == 0 {
		return x
	}

	sign := signbit(x)
	x = abs(x)

	q := uint64(x / pi4)
	if q > 3 {
		q = q & 3
	}

	r := x - float64(int(q))*pi
	if q == 1 || q == 2 {
		r = r - piHi - piLo
	}

	s := r * r
	c := s * r * (tanC1 + s*(tanC2+s*(tanC3+s*(tanC4+s*(tanC5+s*(tanC6+s*tanC7))))))

	if sign {
		return -r - c
	}
	return r + c
}

func nativeAcos(x float64) float64 {
	if isNaN(x) || x > 1 || x < -1 {
		return nan()
	}
	if x == 1 {
		return 0
	}
	if x == -1 {
		return pi
	}

	absX := abs(x)
	if absX > 0.5 {
		return pi - 2*nativeAsin(nativeSqrt(0.5*(1-x)))
	}

	s := x * x
	t := x * (acosA + s*(acosB+s*(acosC+s*(acosD+s*(acosE+s*(acosF+s*acosG))))))
	return pi4 + t - pi4*absX
}

func nativeAsin(x float64) float64 {
	if isNaN(x) || x > 1 || x < -1 {
		return nan()
	}
	if x == 0 {
		return x
	}
	if x == 1 {
		return pi / 2
	}
	if x == -1 {
		return -pi / 2
	}

	sign := signbit(x)
	x = abs(x)

	absX := x
	if absX > 0.5 {
		if sign {
			return -pi/2 + 2*nativeAcos(nativeSqrt(0.5*(1-x)))
		}
		return pi/2 - 2*nativeAcos(nativeSqrt(0.5*(1-x)))
	}

	s := x * x
	t := x * (asinhC1 + s*(asinhC2+s*(asinhC3+s*(asinhC4+s*asinhC5))))

	if sign {
		return -x + t
	}
	return x - t
}

func nativeAtan(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x == 0 {
		return x
	}
	if isInf(x, 1) {
		return pi / 2
	}
	if isInf(x, -1) {
		return -pi / 2
	}

	sign := signbit(x)
	x = abs(x)

	var t float64
	if x > 1 {
		t = pi/2 - 1/x
	} else if x < 1 {
		t = x
	} else {
		t = 0
	}

	s := t * t
	a := s * (atanhC1 + s*(atanhC2+s*(atanhC3+s*(atanhC4+s*(atanhC5+s*atanhC6)))))

	if sign {
		return -t - a
	}
	return t + a
}

func nativeAtan2(y, x float64) float64 {
	if isNaN(y) || isNaN(x) {
		return nan()
	}
	if y == 0 {
		if x > 0 {
			return 0
		}
		if x == 0 {
			return 0
		}
		return pi
	}
	if x == 0 {
		if y > 0 {
			return pi / 2
		}
		return -pi / 2
	}

	if x > 0 {
		return nativeAtan(y / x)
	}
	if x < 0 {
		if y >= 0 {
			return nativeAtan(y/x) + pi
		}
		return nativeAtan(y/x) - pi
	}

	if y > 0 {
		return pi / 2
	}
	return -pi / 2
}

func nativeAsinh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if isInf(x, 0) {
		return x
	}
	if x == 0 {
		return x
	}

	sign := signbit(x)
	x = abs(x)

	if x > 1.0e+50 {
		return nativeLog(2*x) + ln2
	}
	if x > 2 {
		return nativeLog(x + nativeSqrt(x*x+1))
	}
	if x < 1.0e-50 {
		return x
	}

	t := x * x
	x = x + x*t*(asinhC1+t*(asinhC2+t*(asinhC3+t*(asinhC4+t*asinhC5))))

	if sign {
		return -x
	}
	return x
}

func nativeAcosh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x < 1 {
		return nan()
	}
	if x == 1 {
		return 0
	}
	if isInf(x, 0) {
		return x
	}

	if x > 1.0e+100 {
		return nativeLog(2 * x)
	}
	if x > 2 {
		return nativeLog(x + nativeSqrt(x-1)*nativeSqrt(x+1))
	}

	t := x - 1
	if t < 1.0e-50 {
		return nativeSqrt(2*t)
	}

	t = t / (1 + nativeSqrt(1+t*t))
	return nativeLog1p(t) + acoshA
}

func nativeAtanh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x <= -1 || x >= 1 {
		return nan()
	}
	if x == 0 {
		return x
	}
	if x == 1 {
		return inf(1)
	}
	if x == -1 {
		return inf(-1)
	}

	sign := signbit(x)
	x = abs(x)

	if x > 0.5 {
		return -0.5*nativeLog1p(2*(x-1)/(1+x)) + 0.5*ln2
	}
	if x <= 1.0e-50 {
		return x
	}

	t := x + x*x*(atanhC1+x*x*(atanhC2+x*x*(atanhC3+x*x*(atanhC4+x*x*(atanhC5+x*x*atanhC6)))))

	if sign {
		return -t
	}
	return t
}

func nativeLog1p(x float64) float64 {
	if isNaN(x) || x > 1 {
		return x
	}
	if x == 0 || x == -1 {
		return x
	}
	if x < -0.5 {
		return nativeLog(1 + x)
	}

	u := 1 + x
	if u == 1 {
		return x
	}

	lo := x - (u - 1)
	tmp1 := log1pC10 + u*log1pC11
	tmp2 := log1pC9 + u*tmp1
	tmp3 := log1pC8 + u*tmp2
	tmp4 := log1pC7 + u*tmp3
	tmp5 := log1pC6 + u*tmp4
	tmp6 := log1pC5 + u*tmp5
	tmp7 := log1pC4 + u*tmp6
	tmp8 := log1pC3 + u*tmp7
	tmp9 := log1pC2 + u*tmp8
	v := log1pC1 + u*tmp9
	return lo*v + nativeLog(u)
}

func nativeExpm1(x float64) float64 {
	if isNaN(x) || x == 0 {
		return x
	}
	if isInf(x, 1) {
		return inf(1)
	}
	if isInf(x, -1) {
		return -1
	}
	if x > maxLog {
		return inf(1)
	}
	if x < -maxLog {
		return -1
	}

	t := x
	k := 0
	if t < -ln2 {
		k = -1
		t = t + ln2
		if t <= -ln2 {
			k = -2
			t = t + ln2
		}
	}
	t = t - ln2Hi
	c := (t - ln2Hi*float64(k)) - ln2Lo*float64(k)
	t = t * t
	e1 := expm1C10 + t*expm1C11
	e2 := expm1C9 + t*e1
	e3 := expm1C8 + t*e2
	e4 := expm1C7 + t*e3
	e5 := expm1C6 + t*e4
	e6 := expm1C5 + t*e5
	e7 := expm1C4 + t*e6
	e8 := expm1C3 + t*e7
	e9 := expm1C2 + t*e8
	w := t * (expm1C1 + t*e9)
	w = w - c*(2*t-c)
	r := pow2(int64(k)) - 1 + w
	if r < 0 {
		return pow2(int64(k)) + w - 1
	}
	return r
}

func nativePow(x, y float64) float64 {
	if isNaN(y) {
		return nan()
	}
	if isNaN(x) && y != 0 {
		return nan()
	}
	if x == 1 {
		return 1
	}
	if y == 0 {
		return 1
	}
	if y == 1 {
		return x
	}
	if x == 0 {
		if y > 0 {
			return 0
		}
		if y < 0 {
			return inf(1)
		}
		return nan()
	}
	if isInf(y, 0) {
		if x > -1 && x < 1 {
			if isInf(y, 1) {
				return 0
			}
			return inf(1)
		}
		if x > 1 || x < -1 {
			if isInf(y, 1) {
				return inf(1)
			}
			return 0
		}
		return nan()
	}
	if isInf(x, 0) {
		if y > 0 {
			return inf(1)
		}
		return 0
	}

	if y < 0 {
		return 1 / nativePow(x, -y)
	}

	hi := y * nativeLog(x)
	if hi > maxLog || hi < minLog {
		return inf(1)
	}

	return nativeExp(hi)
}

func nativeInv(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x == 0 {
		return inf(1)
	}
	return 1 / x
}

func nativeNeg(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x == 0 {
		return copysign(0, -1)
	}
	return -x
}

func nativeAbs(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x < 0 {
		return -x
	}
	return x
}

func nativeMod(x, y float64) float64 {
	if isNaN(x) || isNaN(y) || y == 0 {
		return nan()
	}
	if y < 0 {
		y = -y
	}
	if x < 0 {
		return -nativeMod(-x, y)
	}
	for x >= y {
		x -= y
	}
	return x
}

func nativeRemainder(x, y float64) float64 {
	if isNaN(x) || isNaN(y) || y == 0 {
		return nan()
	}
	if y < 0 {
		y = -y
	}
	if x < 0 {
		return -nativeMod(-x, y)
	}
	return nativeMod(x, y)
}

func nativeHypot(x, y float64) float64 {
	if isNaN(x) || isNaN(y) {
		return nan()
	}
	if isInf(x, 0) || isInf(y, 0) {
		return inf(1)
	}
	if x == 0 {
		return nativeAbs(y)
	}
	if y == 0 {
		return nativeAbs(x)
	}

	if x < y {
		x, y = y, x
	}

	y = y / x
	return nativeAbs(x * nativeSqrt(1+y*y))
}

func nativeCbrt(x float64) float64 {
	if isNaN(x) || isInf(x, 0) || x == 0 {
		return x
	}

	sign := signbit(x)
	if sign {
		x = -x
	}

	exp := int(f64bits(x)>>52) - 1023
	frac := float64(f64bits(x)&((1<<52)-1)) / float64(1<<52)
	frac = nativePow(frac, 1.0/3.0)
	exp = exp / 3

	result := f64frombits(((uint64(exp+1023) << 52) | (uint64(frac*float64(1<<52)) & ((1<<52)-1))))

	if sign {
		return -result
	}
	return result
}

func nativeMax(x, y float64) float64 {
	if isNaN(x) {
		return y
	}
	if isNaN(y) {
		return x
	}
	if x > y {
		return x
	}
	if x < y {
		return y
	}
	if x == 0 && y == 0 {
		if signbit(x) {
			return y
		}
		return x
	}
	if x > y {
		return x
	}
	return y
}

func nativeMin(x, y float64) float64 {
	if isNaN(x) {
		return y
	}
	if isNaN(y) {
		return x
	}
	if x < y {
		return x
	}
	if x > y {
		return y
	}
	if x == 0 && y == 0 {
		if signbit(x) {
			return x
		}
		return y
	}
	if x < y {
		return x
	}
	return y
}

func nativeLog10(x float64) float64 {
	if isNaN(x) || x > 0 {
		return x
	}
	if x == 0 {
		return inf(-1)
	}
	if x < 0 {
		return nan()
	}

	return nativeLog(x) / ln10()
}

func ln10() float64 {
	return 2.3025850929940456840179914546843642
}