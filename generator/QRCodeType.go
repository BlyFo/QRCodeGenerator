package generator

type EncodingMode uint8

const (
	EncodingMode_Numeric EncodingMode = 1
	EncodingMode_Alpha   EncodingMode = 2
	EncodingMode_Byte    EncodingMode = 4
	EncodingMode_Kanji   EncodingMode = 8
	EncodingMode_ECI     EncodingMode = 7
)

type ErrorLevel uint16

const (
	ErrorLevel_L ErrorLevel = 1 // 7%
	ErrorLevel_M ErrorLevel = 0 // 15%
	ErrorLevel_Q ErrorLevel = 3 // 25%
	ErrorLevel_H ErrorLevel = 2 // 30%
)

type MaskPattern uint16

const (
	MaskPattern_0 MaskPattern = iota
	MaskPattern_1
	MaskPattern_2
	MaskPattern_3
	MaskPattern_4
	MaskPattern_5
	MaskPattern_6
	MaskPattern_7
)

type QRCapacity struct {
	Numeric      int
	AlphaNumeric int
	Binary       int
	Kanji        int
}

type QRCodeInfo struct {
	Version               int
	Size                  int
	ErrorLevel            ErrorLevel
	MaskPatern            MaskPattern
	InfoToEncode          string
	EncodingMode          EncodingMode
	BitsPerData           uint8
	AlignSquareCordenates []int
	MaxNumberOfBits       int
	CodeWords             ERCodeWords
}

type ERCodeWords struct {
	Total                  int
	ECCWPerBlock           int
	BlocksGroup1           int
	DataCodeWordsPerGroup1 int
	BlocksGroup2           int
	DataCodeWordsPerGroup2 int
}

const MAX_SUPPORTED_VERSION = 7

var QRVersionInfo = map[int]map[ErrorLevel]QRCapacity{
	1: {
		ErrorLevel_L: QRCapacity{Numeric: 41, AlphaNumeric: 25, Binary: 17, Kanji: 10},
		ErrorLevel_M: QRCapacity{Numeric: 34, AlphaNumeric: 20, Binary: 14, Kanji: 8},
		ErrorLevel_Q: QRCapacity{Numeric: 27, AlphaNumeric: 16, Binary: 11, Kanji: 7},
		ErrorLevel_H: QRCapacity{Numeric: 17, AlphaNumeric: 10, Binary: 7, Kanji: 4},
	},
	2: {
		ErrorLevel_L: QRCapacity{Numeric: 77, AlphaNumeric: 47, Binary: 32, Kanji: 20},
		ErrorLevel_M: QRCapacity{Numeric: 63, AlphaNumeric: 38, Binary: 26, Kanji: 16},
		ErrorLevel_Q: QRCapacity{Numeric: 48, AlphaNumeric: 29, Binary: 20, Kanji: 12},
		ErrorLevel_H: QRCapacity{Numeric: 34, AlphaNumeric: 20, Binary: 14, Kanji: 8},
	},
	3: {
		ErrorLevel_L: QRCapacity{Numeric: 127, AlphaNumeric: 77, Binary: 53, Kanji: 32},
		ErrorLevel_M: QRCapacity{Numeric: 101, AlphaNumeric: 61, Binary: 42, Kanji: 26},
		ErrorLevel_Q: QRCapacity{Numeric: 77, AlphaNumeric: 47, Binary: 32, Kanji: 20},
		ErrorLevel_H: QRCapacity{Numeric: 58, AlphaNumeric: 35, Binary: 24, Kanji: 15},
	},
	4: {
		ErrorLevel_L: QRCapacity{Numeric: 187, AlphaNumeric: 114, Binary: 78, Kanji: 48},
		ErrorLevel_M: QRCapacity{Numeric: 149, AlphaNumeric: 90, Binary: 62, Kanji: 38},
		ErrorLevel_Q: QRCapacity{Numeric: 111, AlphaNumeric: 67, Binary: 46, Kanji: 28},
		ErrorLevel_H: QRCapacity{Numeric: 82, AlphaNumeric: 50, Binary: 34, Kanji: 21},
	},
	5: {
		ErrorLevel_L: QRCapacity{Numeric: 255, AlphaNumeric: 154, Binary: 106, Kanji: 65},
		ErrorLevel_M: QRCapacity{Numeric: 202, AlphaNumeric: 122, Binary: 84, Kanji: 52},
		ErrorLevel_Q: QRCapacity{Numeric: 144, AlphaNumeric: 87, Binary: 60, Kanji: 37},
		ErrorLevel_H: QRCapacity{Numeric: 106, AlphaNumeric: 64, Binary: 44, Kanji: 27},
	},
	6: {
		ErrorLevel_L: QRCapacity{Numeric: 322, AlphaNumeric: 195, Binary: 134, Kanji: 82},
		ErrorLevel_M: QRCapacity{Numeric: 255, AlphaNumeric: 154, Binary: 106, Kanji: 65},
		ErrorLevel_Q: QRCapacity{Numeric: 178, AlphaNumeric: 108, Binary: 74, Kanji: 45},
		ErrorLevel_H: QRCapacity{Numeric: 139, AlphaNumeric: 84, Binary: 58, Kanji: 36},
	},
	7: {
		ErrorLevel_L: QRCapacity{Numeric: 370, AlphaNumeric: 224, Binary: 154, Kanji: 95},
		ErrorLevel_M: QRCapacity{Numeric: 293, AlphaNumeric: 178, Binary: 122, Kanji: 75},
		ErrorLevel_Q: QRCapacity{Numeric: 207, AlphaNumeric: 125, Binary: 86, Kanji: 53},
		ErrorLevel_H: QRCapacity{Numeric: 154, AlphaNumeric: 93, Binary: 64, Kanji: 39},
	},
}

var MaskPatternByErrorLevel = map[ErrorLevel]map[MaskPattern]uint16{
	ErrorLevel_L: {
		MaskPattern_0: 0b111011111000100,
		MaskPattern_1: 0b111001011110011,
		MaskPattern_2: 0b111110110101010,
		MaskPattern_3: 0b111100010011101,
		MaskPattern_4: 0b110011000101111,
		MaskPattern_5: 0b110001100011000,
		MaskPattern_6: 0b110110001000001,
		MaskPattern_7: 0b110100101110110,
	},
	ErrorLevel_M: {
		MaskPattern_0: 0b101010000010010,
		MaskPattern_1: 0b101000100100101,
		MaskPattern_2: 0b101111001111100,
		MaskPattern_3: 0b101101101001011,
		MaskPattern_4: 0b100010111111001,
		MaskPattern_5: 0b100000011001110,
		MaskPattern_6: 0b100111110010111,
		MaskPattern_7: 0b100101010100000,
	},
	ErrorLevel_Q: {
		MaskPattern_0: 0b011010101011111,
		MaskPattern_1: 0b011000001101000,
		MaskPattern_2: 0b011111100110001,
		MaskPattern_3: 0b011101000000110,
		MaskPattern_4: 0b010010010110100,
		MaskPattern_5: 0b010000110000011,
		MaskPattern_6: 0b010111011011010,
		MaskPattern_7: 0b010101111101101,
	},
	ErrorLevel_H: {
		MaskPattern_0: 0b001011010001001,
		MaskPattern_1: 0b001001110111110,
		MaskPattern_2: 0b001110011100111,
		MaskPattern_3: 0b001100111010000,
		MaskPattern_4: 0b000011101100010,
		MaskPattern_5: 0b000001001010101,
		MaskPattern_6: 0b000110100001100,
		MaskPattern_7: 0b000100000111011,
	},
}

var MaskFunctions = map[MaskPattern]func(x int, y int) bool{
	MaskPattern_0: func(x int, y int) bool {
		return (x+y)%2 == 0
	},
	MaskPattern_1: func(x int, y int) bool {
		return x%2 == 0
	},
	MaskPattern_2: func(x int, y int) bool {
		return y%3 == 0
	},
	MaskPattern_3: func(x int, y int) bool {
		return (x+y)%3 == 0
	},
	MaskPattern_4: func(x int, y int) bool {
		return (y/2+x/3)%2 == 0
	},
	MaskPattern_5: func(x int, y int) bool {
		return ((x*y)%2)+((x*y)%3) == 0
	},
	MaskPattern_6: func(x int, y int) bool {
		return (((x*y)%2)+((x*y)%3))%2 == 0
	},
	MaskPattern_7: func(x int, y int) bool {
		return (((x+y)%2)+((x+y)%3))%2 == 0
	},
}

var QRAlignSquareCordinates = map[int][]int{
	1:  {},
	2:  {6, 18},
	3:  {6, 22},
	4:  {6, 26},
	5:  {6, 30},
	6:  {6, 34},
	7:  {6, 22, 38},
	8:  {6, 24, 42},
	9:  {6, 26, 46},
	10: {6, 28, 50},
	11: {6, 30, 50},
	12: {6, 32, 58},
	13: {6, 34, 62},
	14: {6, 26, 46, 66},
	15: {6, 26, 48, 70},
	16: {6, 26, 50, 74},
	17: {6, 30, 54, 78},
	18: {6, 30, 56, 82},
	19: {6, 30, 58, 86},
	20: {6, 34, 62, 90},
	21: {6, 28, 50, 72, 94},
	22: {6, 26, 50, 74, 98},
	23: {6, 30, 54, 78, 102},
	24: {6, 28, 54, 80, 106},
	25: {6, 32, 58, 84, 110},
	26: {6, 30, 58, 86, 114},
	27: {6, 34, 62, 90, 118},
	28: {6, 26, 50, 74, 98, 122},
	29: {6, 30, 54, 78, 102, 126},
	30: {6, 26, 52, 78, 104, 130},
	31: {6, 30, 56, 82, 108, 134},
	32: {6, 34, 60, 86, 112, 138},
	33: {6, 30, 58, 86, 114, 142},
	34: {6, 34, 62, 90, 118, 146},
	35: {6, 30, 54, 78, 102, 126, 150},
	36: {6, 24, 50, 76, 102, 128, 154},
	37: {6, 28, 54, 80, 106, 132, 158},
	38: {6, 32, 58, 84, 110, 136, 162},
	39: {6, 26, 54, 82, 110, 138, 166},
	40: {6, 30, 58, 86, 114, 142, 170},
}

var ErrorCorrectionCodeWords = map[int]map[ErrorLevel]ERCodeWords{
	1: {
		ErrorLevel_L: ERCodeWords{Total: 19, ECCWPerBlock: 7, BlocksGroup1: 1, DataCodeWordsPerGroup1: 19, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_M: ERCodeWords{Total: 16, ECCWPerBlock: 10, BlocksGroup1: 1, DataCodeWordsPerGroup1: 16, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_Q: ERCodeWords{Total: 13, ECCWPerBlock: 13, BlocksGroup1: 1, DataCodeWordsPerGroup1: 13, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_H: ERCodeWords{Total: 9, ECCWPerBlock: 17, BlocksGroup1: 1, DataCodeWordsPerGroup1: 9, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
	},
	2: {
		ErrorLevel_L: ERCodeWords{Total: 34, ECCWPerBlock: 10, BlocksGroup1: 1, DataCodeWordsPerGroup1: 34, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_M: ERCodeWords{Total: 28, ECCWPerBlock: 16, BlocksGroup1: 1, DataCodeWordsPerGroup1: 28, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_Q: ERCodeWords{Total: 22, ECCWPerBlock: 22, BlocksGroup1: 1, DataCodeWordsPerGroup1: 22, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_H: ERCodeWords{Total: 16, ECCWPerBlock: 28, BlocksGroup1: 1, DataCodeWordsPerGroup1: 16, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
	},
	3: {
		ErrorLevel_L: ERCodeWords{Total: 55, ECCWPerBlock: 15, BlocksGroup1: 1, DataCodeWordsPerGroup1: 55, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_M: ERCodeWords{Total: 44, ECCWPerBlock: 26, BlocksGroup1: 1, DataCodeWordsPerGroup1: 44, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_Q: ERCodeWords{Total: 34, ECCWPerBlock: 18, BlocksGroup1: 1, DataCodeWordsPerGroup1: 17, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_H: ERCodeWords{Total: 26, ECCWPerBlock: 22, BlocksGroup1: 1, DataCodeWordsPerGroup1: 13, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
	},
	4: {
		ErrorLevel_L: ERCodeWords{Total: 80, ECCWPerBlock: 20, BlocksGroup1: 1, DataCodeWordsPerGroup1: 80, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_M: ERCodeWords{Total: 64, ECCWPerBlock: 18, BlocksGroup1: 2, DataCodeWordsPerGroup1: 32, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_Q: ERCodeWords{Total: 48, ECCWPerBlock: 26, BlocksGroup1: 2, DataCodeWordsPerGroup1: 24, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_H: ERCodeWords{Total: 36, ECCWPerBlock: 16, BlocksGroup1: 1, DataCodeWordsPerGroup1: 9, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
	},
	5: {
		ErrorLevel_L: ERCodeWords{Total: 108, ECCWPerBlock: 26, BlocksGroup1: 1, DataCodeWordsPerGroup1: 108, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_M: ERCodeWords{Total: 86, ECCWPerBlock: 24, BlocksGroup1: 2, DataCodeWordsPerGroup1: 43, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_Q: ERCodeWords{Total: 62, ECCWPerBlock: 18, BlocksGroup1: 2, DataCodeWordsPerGroup1: 15, BlocksGroup2: 2, DataCodeWordsPerGroup2: 16},
		ErrorLevel_H: ERCodeWords{Total: 46, ECCWPerBlock: 22, BlocksGroup1: 2, DataCodeWordsPerGroup1: 11, BlocksGroup2: 2, DataCodeWordsPerGroup2: 12},
	},
	6: {
		ErrorLevel_L: ERCodeWords{Total: 136, ECCWPerBlock: 18, BlocksGroup1: 2, DataCodeWordsPerGroup1: 68, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_M: ERCodeWords{Total: 108, ECCWPerBlock: 16, BlocksGroup1: 4, DataCodeWordsPerGroup1: 27, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_Q: ERCodeWords{Total: 76, ECCWPerBlock: 24, BlocksGroup1: 4, DataCodeWordsPerGroup1: 19, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_H: ERCodeWords{Total: 60, ECCWPerBlock: 28, BlocksGroup1: 4, DataCodeWordsPerGroup1: 15, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
	},
	7: {
		ErrorLevel_L: ERCodeWords{Total: 156, ECCWPerBlock: 20, BlocksGroup1: 2, DataCodeWordsPerGroup1: 78, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_M: ERCodeWords{Total: 124, ECCWPerBlock: 18, BlocksGroup1: 4, DataCodeWordsPerGroup1: 31, BlocksGroup2: 0, DataCodeWordsPerGroup2: 0},
		ErrorLevel_Q: ERCodeWords{Total: 88, ECCWPerBlock: 18, BlocksGroup1: 2, DataCodeWordsPerGroup1: 14, BlocksGroup2: 4, DataCodeWordsPerGroup2: 15},
		ErrorLevel_H: ERCodeWords{Total: 66, ECCWPerBlock: 26, BlocksGroup1: 2, DataCodeWordsPerGroup1: 13, BlocksGroup2: 1, DataCodeWordsPerGroup2: 14},
	},
}

var ReminderBits = map[int]int{
	1:  0,
	2:  7,
	3:  7,
	4:  7,
	5:  7,
	6:  7,
	7:  0,
	8:  0,
	9:  0,
	10: 0,
	11: 0,
	12: 0,
	13: 0,
}

var VersionInformationString = map[int][]bool{
	0: {},
	1: {},
	2: {},
	3: {},
	4: {},
	5: {},
	6: {},
	7: {false, false, false, true, true, true, true, true, false, false, true, false, false, true, false, true, false, false},
	8: {false, false, true, false, false, false, false, true, false, true, true, false, true, true, true, true, false, false},
}

func GetErrorLevels() []ErrorLevel {
	return []ErrorLevel{ErrorLevel_L, ErrorLevel_M, ErrorLevel_Q, ErrorLevel_H}
}

func GetMaskPatterns() []MaskPattern {
	return []MaskPattern{MaskPattern_0, MaskPattern_1, MaskPattern_2, MaskPattern_3, MaskPattern_4, MaskPattern_5, MaskPattern_6, MaskPattern_7}
}

func GetCharacterCountIndicator(version int, encodingMode EncodingMode) uint8 {
	if version <= 9 {
		switch encodingMode {
		case EncodingMode_Numeric:
			return 10
		case EncodingMode_Alpha:
			return 9
		case EncodingMode_Byte:
			return 8
		case EncodingMode_Kanji:
			return 8
		}
	}
	if version <= 26 {
		switch encodingMode {
		case EncodingMode_Numeric:
			return 12
		case EncodingMode_Alpha:
			return 11
		case EncodingMode_Byte:
			return 16
		case EncodingMode_Kanji:
			return 10
		}
	}
	if version <= 40 {
		switch encodingMode {
		case EncodingMode_Numeric:
			return 14
		case EncodingMode_Alpha:
			return 13
		case EncodingMode_Byte:
			return 16
		case EncodingMode_Kanji:
			return 12
		}
	}
	//should never happen
	return 0
}

var AlphaEncodeDict = map[string]uint16{
	"0": 0,
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"A": 10,
	"B": 11,
	"C": 12,
	"D": 13,
	"E": 14,
	"F": 15,
	"G": 16,
	"H": 17,
	"I": 18,
	"J": 19,
	"K": 20,
	"L": 21,
	"M": 22,
	"N": 23,
	"O": 24,
	"P": 25,
	"Q": 26,
	"R": 27,
	"S": 28,
	"T": 29,
	"U": 30,
	"V": 31,
	"W": 32,
	"X": 33,
	"Y": 34,
	"Z": 35,
	" ": 36,
	"$": 37,
	"%": 38,
	"*": 39,
	"+": 40,
	"-": 41,
	".": 42,
	"/": 43,
	":": 44,
}

// helper functions
func (b ErrorLevel) String() string {
	switch b {
	case ErrorLevel_H:
		return "H (High 30%)"
	case ErrorLevel_L:
		return "L (Low 7%)"
	case ErrorLevel_Q:
		return "Q (Quartile 25%)"
	case ErrorLevel_M:
		return "M (Medium 15%)"
	}
	return "Error" //should never happen
}

func (b EncodingMode) String() string {
	switch b {
	case EncodingMode_Alpha:
		return "Alpha"
	case EncodingMode_Byte:
		return "Byte"
	case EncodingMode_ECI:
		return "ECI"
	case EncodingMode_Kanji:
		return "Kanji"
	case EncodingMode_Numeric:
		return "Numeric"
	}
	return "Error" //should never happen
}
