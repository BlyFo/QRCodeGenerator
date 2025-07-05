package utils

func GetMax(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func Pow16(x, y uint16) uint16 {
	if y == 0 {
		return 1
	}
	return x * Pow16(x, y-1)
}

func DeepCopy2D[T any](input [][]T) [][]T {
	output := make([][]T, len(input))
	for i, row := range input {
		output[i] = make([]T, len(row))
		copy(output[i], row)
	}
	return output
}

func BoolToInt(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func ByteToBoolArray(b byte) []bool {
	var boolArray []bool
	for j := 7; j >= 0; j-- {
		boolArray = append(boolArray, (b&(1<<j)) != 0)
	}
	return boolArray
}

func Byte16ToBoolArray(b uint16) []bool {
	boolArray := make([]bool, 16)
	for j := 0; j < 16; j++ {
		boolArray[j] = (b & (1 << (15 - j))) != 0
	}
	return boolArray
}

func BoolArrayToByte(b []bool) []uint8 {
	byteArray := make([]uint8, len(b)/8)
	for i := 0; i < len(b)/8; i++ {
		offset := i * 8
		var number uint8 = 0
		number += BoolToInt(b[offset]) * 128
		number += BoolToInt(b[offset+1]) * 64
		number += BoolToInt(b[offset+2]) * 32
		number += BoolToInt(b[offset+3]) * 16
		number += BoolToInt(b[offset+4]) * 8
		number += BoolToInt(b[offset+5]) * 4
		number += BoolToInt(b[offset+6]) * 2
		number += BoolToInt(b[offset+7]) * 1
		byteArray[i] = number
	}
	return byteArray
}

func getLogAndAntiLogTables() ([]uint8, []uint8) {
	logTable := make([]uint8, 256)
	antiLogTable := make([]uint8, 256)

	for exponent, value := 1, 1; exponent < 256; exponent++ {
		if value > 127 {
			value = (value << 1) ^ 285
		} else {
			value = value << 1
		}
		logTable[value] = uint8(exponent) % 255
		antiLogTable[exponent%255] = uint8(value)
	}
	return logTable, antiLogTable
}

var LOG, ANTI_LOG = getLogAndAntiLogTables()

func gFmul(a uint8, b uint8) uint8 {
	if a == 0 || b == 0 {
		return 0
	}
	// hay que convertrlos porque se suman y terminaria dando overflow y arruinando el resultado
	firstTerm := uint16(LOG[a])
	secondTerm := uint16(LOG[b])
	result := (firstTerm + secondTerm) % 255
	return ANTI_LOG[result]
}

func gFdiv(a uint8, b uint8) uint8 {
	// hay que convertrlos porque se suman y terminaria dando overflow y arruinando el resultado
	// 255*254 da poco menos que el max de uint16
	firstTerm := uint16(LOG[a])
	secondTerm := uint16(LOG[b]) * 254
	result := (firstTerm + secondTerm) % 255
	return ANTI_LOG[result]
}

func polyMul(poly1 []uint8, poly2 []uint8) []uint8 {
	// This is going to be the product polynomial, that we pre-allocate.
	// We know it's going to be `poly1.length + poly2.length - 1` long.
	coeffs := make([]uint8, len(poly1)+len(poly2)-1)

	// Instead of executing all the steps in the example, we can jump to
	// computing the coefficients of the result
	for index := range coeffs {
		// tlvez aca ya que javascript use double y no uint8, lo que podria explicar porque da mal para valores altos
		// or maybe not since we do  ^= and that can't return a number bigger than the inputs
		coeff := uint8(0)
		for p1index := 0; p1index <= index; p1index++ {
			p2index := index - p1index
			// We *should* do better here, as `p1index` and `p2index` could
			// be out of range, but `mul` defined above will handle that case.
			// Just beware of that when implementing in other languages.
			if len(poly1) <= p1index || len(poly2) <= p2index {
				coeff ^= 0
			} else {
				coeff ^= gFmul(poly1[p1index], poly2[p2index])
			}
		}
		coeffs[index] = coeff
	}
	return coeffs
}

func polyRest(dividend []uint8, divisor []uint8) []uint8 {
	quotientLength := len(dividend) - len(divisor) + 1
	// Let's just say that the dividend is the rest right away
	rest := dividend
	for range quotientLength {
		// If the first term is 0, we can just skip this iteration
		if rest[0] == 0 {
			rest = rest[1:]
			continue
		}

		factor := []uint8{gFdiv(rest[0], divisor[0])}
		subtr := make([]uint8, len(rest))
		copy(subtr, polyMul(divisor, factor))
		newReast := make([]uint8, len(rest))
		for index, value := range rest {
			newReast[index] = value ^ subtr[index]
		}

		rest = newReast[1:]
	}
	return rest
}

func getGeneratorPoly(degree int) []uint8 {
	lastPoly := []uint8{1}
	for index := range degree {
		lastPoly = polyMul(lastPoly, []uint8{1, ANTI_LOG[index]})
	}
	return lastPoly
}

func GetErrorCorrectionCodegowrds(data []uint8, codewords int) []uint8 {
	degree := codewords

	messagePoly := make([]uint8, len(data)+codewords)
	copy(messagePoly, data)
	return polyRest(messagePoly, getGeneratorPoly(degree))

}

func IsNumericString(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func IsAphaNumeric(s string) bool {
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == ' ':
		case r == '$', r == '%', r == '*', r == '+', r == '-', r == '.', r == '/', r == ':':
		default:
			return false
		}
	}
	return true
}
