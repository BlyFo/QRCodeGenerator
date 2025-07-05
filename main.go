package main

import (
	"QRCodeGenerator/drawer"
	"QRCodeGenerator/generator"
	logger "QRCodeGenerator/logger"
	"QRCodeGenerator/utils"
	"errors"
	"math"
	"strconv"
	"strings"
)

type QRCodeInfo = generator.QRCodeInfo
type ErrorLevel = generator.ErrorLevel
type QRCapacity = generator.QRCapacity
type MaskPattern = generator.MaskPattern
type ERCodeWords = generator.ERCodeWords

const (
	QR_CODE_STEP_ENCODE_MODE uint8 = iota
	QR_CODE_STEP_CHARACTER_COUNT
	QR_CODE_STEP_ENCODE_DATA
	QR_CODE_STEP_ERROR_CORRECTION
	QR_CODE_STEP_MASK
)

func getCapacityByEncodeMode(capacity QRCapacity, encodedMode generator.EncodingMode) int {
	switch encodedMode {
	case generator.EncodingMode_Byte:
		return capacity.Binary
	case generator.EncodingMode_Alpha:
		return capacity.AlphaNumeric
	case generator.EncodingMode_Numeric:
		return capacity.Numeric
	}
	panic("Error not found Encoded mode")
}

func getEncodeMode(data string) generator.EncodingMode {
	if utils.IsNumericString(data) {
		return generator.EncodingMode_Numeric
	}
	if utils.IsAphaNumeric(data) {
		return generator.EncodingMode_Alpha
	}
	return generator.EncodingMode_Byte
}

func getQRInfoByData(stringToEncode string) (QRCodeInfo, error) {

	QRVersionInfo := generator.QRVersionInfo
	QRAlignSquareCordinates := generator.QRAlignSquareCordinates
	encodedMode := getEncodeMode(stringToEncode)

	errorLevels := generator.GetErrorLevels()
	errorLevelsSize := len(errorLevels)
	stringToEncodeSize := len(stringToEncode)

	for version := 1; version < generator.MAX_SUPPORTED_VERSION+1; version++ {
		capacity := QRVersionInfo[version]

		for errorLevelIndex := errorLevelsSize - 1; errorLevelIndex >= 0; errorLevelIndex-- {

			errorLevel := errorLevels[errorLevelIndex]
			errorLevelCapacity := getCapacityByEncodeMode(capacity[errorLevel], encodedMode)
			QRCodeWords := generator.ErrorCorrectionCodeWords[version][errorLevel]

			if errorLevelCapacity >= stringToEncodeSize {
				QRInfo := QRCodeInfo{
					Version:               version,
					Size:                  4*version + 17,
					ErrorLevel:            errorLevel,
					MaskPatern:            generator.MaskPattern_2, // esto se pisara mas adelante
					InfoToEncode:          stringToEncode,
					EncodingMode:          encodedMode,
					BitsPerData:           generator.GetCharacterCountIndicator(version, encodedMode),
					CodeWords:             QRCodeWords,
					AlignSquareCordenates: QRAlignSquareCordinates[version],
					MaxNumberOfBits:       errorLevelCapacity,
				}
				return QRInfo, nil
			}
		}
	}

	return QRCodeInfo{}, errors.New("no se encontro version compatible")
}

func drawSquarePattern(QRArray [][]uint8, x int, y int, radius int) {
	x0 := x - radius
	y0 := y - radius
	x1 := x + radius
	y1 := y + radius

	for i := x0; i <= x1; i++ {
		for j := y0; j <= y1; j++ {

			//midle white layer
			if (j >= y0+1 && j <= y1-1 && (i == x0+1 || i == x1-1)) || ((j == y0+1 || j == y1-1) && i >= x0+1 && i <= x1-1) {
				QRArray[i][j] = drawer.WHITE_COLOR
				continue
			}

			QRArray[i][j] = drawer.BLACK_COLOR
		}
	}
}

func addPositionSquare(QRArray [][]uint8, size int) {
	drawSquarePattern(QRArray, 3, 3, 3)      //upper left
	drawSquarePattern(QRArray, 3, size-4, 3) //upper right
	drawSquarePattern(QRArray, size-4, 3, 3) //lower left

	//white borders (separators)
	color := drawer.WHITE_COLOR
	for i := 0; i < 8; i++ {
		//vertical borders
		QRArray[i][7] = color
		QRArray[i][size-8] = color
		QRArray[size-8+i][7] = color

		//horizontal borders
		QRArray[7][i] = color
		QRArray[7][size-8+i] = color
		QRArray[size-8][i] = color
	}

}

func addAlignSquares(QRArray [][]uint8, alignSquareCordenates []int) {
	for _, i := range alignSquareCordenates {
		for _, j := range alignSquareCordenates {
			if QRArray[i][j] == 0 {
				drawSquarePattern(QRArray, i, j, 2)
			}
		}
	}
}

func addTiming(QRArray [][]uint8) {
	size := len(QRArray)
	color1 := drawer.BLACK_COLOR
	color2 := drawer.BLACK_COLOR
	for i := 0; i < size; i++ {
		if (QRArray)[6][i] == 0 {
			(QRArray)[6][i] = uint8(color1)
			color1 = (color1 % 2) + 1
		}
		if (QRArray)[i][6] == 0 {
			(QRArray)[i][6] = uint8(color2)
			color2 = (color2 % 2) + 1
		}
	}
}

func addFormatVersion(QRArray [][]uint8, errorLevel ErrorLevel, maskPatern MaskPattern, size int) {

	informationString := generator.MaskPatternByErrorLevel

	formatString := informationString[errorLevel][maskPatern]
	binaryFormatString := utils.Byte16ToBoolArray(formatString)
	binaryFormatString = binaryFormatString[1:] //only the last 15 are used

	//left to right
	for i, j := 0, 0; i < len(QRArray); i, j = i+1, j+1 {
		if i == 6 {
			i++
		}
		if j == 7 {
			i = len(QRArray) - 8
		}
		if !binaryFormatString[j] {
			(QRArray)[8][i] = drawer.WHITE_COLOR
		} else {
			(QRArray)[8][i] = drawer.BLACK_COLOR
		}
	}

	//bottom to top
	for i, j := 0, 0; i < len(QRArray); i, j = i+1, j+1 {
		if j == 7 {
			i = len(QRArray) - 9
		}
		if j == 9 {
			i++
		}
		if !binaryFormatString[j] {
			(QRArray)[size-i-1][8] = drawer.WHITE_COLOR
		} else {
			(QRArray)[size-i-1][8] = drawer.BLACK_COLOR
		}
	}
}

func writeOrder(value int, index int) int {
	if index%2 == 0 {
		return value - 1
	}
	return value + 1
}

func generateQRTemplate(QRVersionInfo QRCodeInfo) [][]uint8 {
	QRArray := make([][]uint8, QRVersionInfo.Size)
	for i := 0; i < QRVersionInfo.Size; i++ {
		QRArray[i] = make([]uint8, QRVersionInfo.Size)
	}

	addPositionSquare(QRArray, QRVersionInfo.Size)
	addAlignSquares(QRArray, QRVersionInfo.AlignSquareCordenates)
	addTiming(QRArray)
	addFormatVersion(QRArray, QRVersionInfo.ErrorLevel, QRVersionInfo.MaskPatern, QRVersionInfo.Size)
	//add black square
	QRArray[QRVersionInfo.Size-8][8] = drawer.BLACK_COLOR

	return QRArray
}

func getEncodeMode_Binary(QRVersionInfo QRCodeInfo) []bool {
	return utils.ByteToBoolArray(byte(QRVersionInfo.EncodingMode))[4:]
}

func getCharacterCount_Binary(QRVersionInfo QRCodeInfo) []bool {
	dataToEncode := QRVersionInfo.InfoToEncode
	bitsForCharacterCount := generator.GetCharacterCountIndicator(QRVersionInfo.Version, QRVersionInfo.EncodingMode)
	// this should use GetCharacterCountIndicator isntead of always asuming 8
	return utils.Byte16ToBoolArray(uint16(len(dataToEncode)))[16-bitsForCharacterCount:]
}

func getString_Encoded_Byte(QRVersionInfo QRCodeInfo) []bool {
	var encodedData []bool
	dataToEncode := QRVersionInfo.InfoToEncode
	maxBitesforQr := QRVersionInfo.MaxNumberOfBits * 8
	paddingBites := maxBitesforQr - len(dataToEncode)

	for i := 0; i < len(dataToEncode); i++ {
		// Convertir cada caracter a su valor byte (8 bits)
		charToByte := byte(dataToEncode[i])
		encodedData = append(encodedData, utils.ByteToBoolArray(charToByte)...)
	}

	//add pading bytes
	for i := 0; i < 4 && i < paddingBites; i++ {
		encodedData = append(encodedData, false)
	}

	return encodedData
}

func getString_Encoded_Numeric(QRVersionInfo QRCodeInfo) []bool {
	encodedString := make([]bool, 0, (len(QRVersionInfo.InfoToEncode)/3)*10)
	numericString := strings.TrimLeft(QRVersionInfo.InfoToEncode, "0")
	extraNumbers := ""

	if len(numericString)%3 != 0 {
		extraNumbers = numericString[len(numericString)-(len(numericString)%3):]
		numericString = numericString[:len(numericString)-len(extraNumbers)]
	}

	for i := range len(numericString) / 3 {
		numToEncode, _ := strconv.ParseUint(numericString[i*3:(i+1)*3], 10, 16)
		encodedNum := utils.Byte16ToBoolArray(uint16(numToEncode))[6:]
		encodedString = append(encodedString, encodedNum...)
	}

	if len(extraNumbers) > 0 {
		numToEncode, _ := strconv.ParseUint(extraNumbers, 10, 16)
		encodedNum := utils.Byte16ToBoolArray(uint16(numToEncode))

		if len(extraNumbers) == 1 {
			encodedString = append(encodedString, encodedNum[12:]...)
		} else if len(extraNumbers) == 2 {
			encodedString = append(encodedString, encodedNum[9:]...)
		}
	}
	return encodedString
}

func getString_Encoded_Alpha(QRVersionInfo QRCodeInfo) []bool {
	var encodedData []bool
	alphaDic := generator.AlphaEncodeDict
	dataToEncode := QRVersionInfo.InfoToEncode
	oddCharacter := ""

	if len(dataToEncode)%2 == 1 {
		oddCharacter = string(dataToEncode[len(dataToEncode)-1])
		dataToEncode = dataToEncode[:len(dataToEncode)-1]
	}

	for i := range len(dataToEncode) / 2 {
		encodedNumber := (45 * alphaDic[string(dataToEncode[i*2])]) + alphaDic[string(dataToEncode[(i*2)+1])]
		numberInBites := utils.Byte16ToBoolArray(encodedNumber)[5:]
		encodedData = append(encodedData, numberInBites...)
	}

	if len(oddCharacter) > 0 {
		encodedCharacter := utils.Byte16ToBoolArray(alphaDic[oddCharacter])[10:]
		encodedData = append(encodedData, encodedCharacter...)
	}

	return encodedData
}

func getString_Encoded(QRVersionInfo QRCodeInfo) []bool {
	switch QRVersionInfo.EncodingMode {
	case generator.EncodingMode_Byte:
		return getString_Encoded_Byte(QRVersionInfo)
	case generator.EncodingMode_Alpha:
		return getString_Encoded_Alpha(QRVersionInfo)
	case generator.EncodingMode_Numeric:
		return getString_Encoded_Numeric(QRVersionInfo)
	case generator.EncodingMode_Kanji:
		return []bool{}
	}
	panic("Encode Mode not Found: " + QRVersionInfo.EncodingMode.String())
}

func addDataToQRCode(QRArray [][]uint8, QRVersionInfo QRCodeInfo, data []bool) [][]uint8 {
	QRArrayCopy := utils.DeepCopy2D(QRArray)
	index := 0
	row := 0
	j := QRVersionInfo.Size
	for i := (QRVersionInfo.Size - 1); i >= 0; i -= 2 { //rigth to left
		j = writeOrder(j, row)
		// ignore the left vertical timming
		if i == 6 {
			i--
		}
		for j >= 0 && j < QRVersionInfo.Size { //down to up or up to down
			for k := 0; k < 2; k++ {
				if index >= len(data) {
					// some versions have empty bites at the end between 7 and 0,
					// for example versions 2 to 6 have 7 empty bites
					if QRArrayCopy[j][i-k] != 0 {
						// not re write used cells
						continue
					}
					QRArrayCopy[j][i-k] = drawer.RED_COLOR // debug
					continue
				}

				if i == 0 && k > 0 {
					//last column right to left
					continue
				}

				if QRArrayCopy[j][i-k] != 0 {
					// not re write used cells
					continue
				}
				if data[index] {
					QRArrayCopy[j][i-k] = drawer.BLACK_COLOR
				} else {
					QRArrayCopy[j][i-k] = drawer.WHITE_COLOR
				}

				index++
			}
			j = writeOrder(j, row)
		}
		row++
	}
	return QRArrayCopy
}

func getPadingBits_Binary(QRVersionInfo QRCodeInfo, dataLenght int) []bool {
	var paddingBits []bool
	totalSpace := QRVersionInfo.CodeWords.Total * 8

	//make the data multiple of 8
	if dataLenght%8 != 0 {
		newBitsLenght := 8 - (dataLenght % 8)
		newBits := make([]bool, newBitsLenght)
		paddingBits = append(paddingBits, newBits...)
	}

	spaceLeft := totalSpace - (dataLenght + len(paddingBits))
	if spaceLeft == 0 {
		return paddingBits
	}

	//se podria hardcodear para no tener que hacer append
	constantBitsForPading := make([]bool, 0, 16)
	constantBitsForPading = append(constantBitsForPading, utils.ByteToBoolArray(0b11101100)...)
	constantBitsForPading = append(constantBitsForPading, utils.ByteToBoolArray(0b00010001)...)

	for i := 0; i < spaceLeft; i++ {
		paddingBits = append(paddingBits, constantBitsForPading[i%16])
	}

	return paddingBits
}

func getCodeWords_Encoded(QRVersionInfo QRCodeInfo, data []bool) [][]bool {
	infoToEncode := utils.BoolArrayToByte(data)
	codeWordsInfo := QRVersionInfo.CodeWords
	codeBlocksCount := codeWordsInfo.BlocksGroup1 + codeWordsInfo.BlocksGroup2
	codeBlocksArrays := make([][]uint8, 0, codeBlocksCount)

	//first block
	for i := range codeWordsInfo.BlocksGroup1 {
		dataForFirstBlock := infoToEncode[i*codeWordsInfo.DataCodeWordsPerGroup1 : (i+1)*codeWordsInfo.DataCodeWordsPerGroup1]
		codewordsFirstBlock := utils.GetErrorCorrectionCodegowrds(dataForFirstBlock, codeWordsInfo.ECCWPerBlock)
		codeBlocksArrays = append(codeBlocksArrays, codewordsFirstBlock)
	}

	//second block
	for i := range codeWordsInfo.BlocksGroup2 {
		offset := codeWordsInfo.BlocksGroup1 * codeWordsInfo.DataCodeWordsPerGroup1
		dataForSecondBlock := infoToEncode[i*codeWordsInfo.DataCodeWordsPerGroup2+offset : (i+1)*codeWordsInfo.DataCodeWordsPerGroup2+offset]
		codewordsSecondBlock := utils.GetErrorCorrectionCodegowrds(dataForSecondBlock, codeWordsInfo.ECCWPerBlock)
		codeBlocksArrays = append(codeBlocksArrays, codewordsSecondBlock)
	}

	encodedCodeWrods := make([][]bool, codeBlocksCount)
	for i := range codeBlocksCount {
		for j := range codeWordsInfo.ECCWPerBlock {
			encodedCodeWrods[i] = append(encodedCodeWrods[i], utils.ByteToBoolArray(codeBlocksArrays[i][j])...)
		}
	}

	return encodedCodeWrods
}

func getStructuredFinalMessage(QRVersionInfo QRCodeInfo, dataCodeWords []bool, ErrorCorrectioncodeWords [][]bool) []bool {
	if len(ErrorCorrectioncodeWords) == 0 {
		return dataCodeWords
	}
	//TODO optimize the code by creating the data code words groups using only 1 for loop (it won't inpact performance)
	//and find out why I'm using the +2 in the make below this
	finalMessage := make([]bool, 0, len(dataCodeWords)+len(ErrorCorrectioncodeWords)+2) // porque rayos le puse +2 aca ?

	//Format Data Code Words
	dataCodeWordsInGroups := make([][]bool, QRVersionInfo.CodeWords.BlocksGroup2+QRVersionInfo.CodeWords.BlocksGroup1)

	blockGroup1 := QRVersionInfo.CodeWords.BlocksGroup1
	DCWsPerGroup1 := QRVersionInfo.CodeWords.DataCodeWordsPerGroup1 * 8

	blockGroup2 := QRVersionInfo.CodeWords.BlocksGroup2
	DCWsPerGroup2 := QRVersionInfo.CodeWords.DataCodeWordsPerGroup2 * 8

	// get Data Codewords for first Group
	for i := range blockGroup1 {
		dataCodeWordsInGroups[i] = dataCodeWords[i*DCWsPerGroup1 : DCWsPerGroup1*(i+1)]
	}

	// get Data Codewords for Second Group
	DCWInFirstGroup := QRVersionInfo.CodeWords.BlocksGroup1 * DCWsPerGroup1
	for i := range blockGroup2 {
		dataCodeWordsInGroups[i+DCWsPerGroup1] = dataCodeWords[i*DCWsPerGroup2+DCWInFirstGroup : (i+1)*DCWsPerGroup2+DCWInFirstGroup]
	}

	//intervale Data Code Words
	bigestCodeWrdsLenght := utils.GetMax(QRVersionInfo.CodeWords.DataCodeWordsPerGroup1, QRVersionInfo.CodeWords.DataCodeWordsPerGroup2)
	for i := range bigestCodeWrdsLenght {
		for j := range len(dataCodeWordsInGroups) {
			//groups have differente amount of data code words
			if i < len(dataCodeWordsInGroups[j]) {
				finalMessage = append(finalMessage, dataCodeWordsInGroups[j][i*8:(i+1)*8]...)
			}
		}
	}

	//Format Error Correction CodeWords
	for i := range QRVersionInfo.CodeWords.ECCWPerBlock {
		for j := range len(ErrorCorrectioncodeWords) {
			finalMessage = append(finalMessage, ErrorCorrectioncodeWords[j][i*8:(i+1)*8]...)
		}
	}

	return finalMessage
}

func applyMask(maskpatern MaskPattern, errorLevel ErrorLevel, QRTemplate [][]uint8, QRFinal [][]uint8) [][]uint8 {
	maskPatternFunction := generator.MaskFunctions[maskpatern]

	for i := range len(QRTemplate) {
		for j := range len(QRTemplate[0]) {
			if QRTemplate[i][j] != 0 {
				continue
			}
			if maskPatternFunction(i, j) {
				if QRFinal[i][j] == drawer.BLACK_COLOR {
					QRFinal[i][j] = drawer.WHITE_COLOR
				} else {
					QRFinal[i][j] = drawer.BLACK_COLOR
				}
			}
		}
	}
	addFormatVersion(QRFinal, errorLevel, maskpatern, len(QRTemplate))
	return QRFinal
}

func getBestMaskPattern(QRVersionInfo QRCodeInfo, QRTemplate [][]uint8, QRFinal [][]uint8) MaskPattern {
	maskPaterns := generator.GetMaskPatterns()
	var bestMask MaskPattern
	var lowestScore uint = 4294967295 //max uint
	logger.Info("Finding best mask pattern")

	for maskIndex := range len(maskPaterns) {
		QRFinalCopy := utils.DeepCopy2D(QRFinal)
		//overwrite the temporal mask infomration
		addFormatVersion(QRFinalCopy, QRVersionInfo.ErrorLevel, maskPaterns[maskIndex], QRVersionInfo.Size)
		QrArrayWithMask := applyMask(maskPaterns[maskIndex], QRVersionInfo.ErrorLevel, QRTemplate, QRFinalCopy)

		penaltyPoints1 := uint(0)
		penaltyPoints2 := uint(0)
		penaltyPoints3 := uint(0)
		penaltyPoints4 := uint(0)

		continuousBlocksH := uint(0)
		continuousBlocksV := uint(0)
		currentColorH := QrArrayWithMask[0][0]
		currentColorV := QrArrayWithMask[0][0]

		PaternPenalty3 := []uint8{drawer.BLACK_COLOR, drawer.WHITE_COLOR, drawer.BLACK_COLOR, drawer.BLACK_COLOR, drawer.BLACK_COLOR, drawer.WHITE_COLOR, drawer.BLACK_COLOR, drawer.WHITE_COLOR, drawer.WHITE_COLOR, drawer.WHITE_COLOR, drawer.WHITE_COLOR}
		PaternInvertedPenalty3 := []uint8{drawer.WHITE_COLOR, drawer.WHITE_COLOR, drawer.WHITE_COLOR, drawer.WHITE_COLOR, drawer.BLACK_COLOR, drawer.WHITE_COLOR, drawer.BLACK_COLOR, drawer.BLACK_COLOR, drawer.BLACK_COLOR, drawer.WHITE_COLOR, drawer.BLACK_COLOR}

		blackScuares := 0
		whiteSqueares := 0

		for i := range QRVersionInfo.Size {
			for j := range QRVersionInfo.Size {

				//penalty rule 1
				if QrArrayWithMask[i][j] == currentColorH {
					continuousBlocksH++
				} else {
					if continuousBlocksH >= 3 {
						penaltyPoints1 += continuousBlocksH - 2
					}
					continuousBlocksH = 1
					currentColorH = QrArrayWithMask[i][j]
				}

				if QrArrayWithMask[j][i] == currentColorV {
					continuousBlocksV++
				} else {
					if continuousBlocksV >= 3 {
						penaltyPoints1 += continuousBlocksV - 2
					}
					continuousBlocksV = 1
					currentColorV = QrArrayWithMask[j][i]
				}

				//penalty rule 2
				if j < QRVersionInfo.Size-1 && i < QRVersionInfo.Size-1 {
					if QrArrayWithMask[i][j] == QrArrayWithMask[i+1][j+1] && QrArrayWithMask[i][j] == QrArrayWithMask[i+1][j] && QrArrayWithMask[i][j] == QrArrayWithMask[i][j+1] {
						penaltyPoints2 += 3
					}
				}

				//penalty rule 3
				if !(i > QRVersionInfo.Size-len(PaternPenalty3) && j > QRVersionInfo.Size-len(PaternPenalty3)) {
					if i < QRVersionInfo.Size-len(PaternPenalty3) {
						patternChecks := true
						invertedPatternChecks := true
						for k := 0; k < len(PaternPenalty3); k++ {
							patternChecks = patternChecks && PaternPenalty3[k] == QrArrayWithMask[i+k][j]
							invertedPatternChecks = invertedPatternChecks && PaternInvertedPenalty3[k] == QrArrayWithMask[i+k][j]
						}
						if patternChecks || invertedPatternChecks {
							penaltyPoints3 += 40
						}
					}
					if j < QRVersionInfo.Size-len(PaternPenalty3) {
						patternChecks := true
						invertedPatternChecks := true
						for k := 0; k < len(PaternPenalty3); k++ {
							patternChecks = patternChecks && PaternPenalty3[k] == QrArrayWithMask[i][j+k]
							invertedPatternChecks = invertedPatternChecks && PaternInvertedPenalty3[k] == QrArrayWithMask[i][j+k]
						}
						if patternChecks || invertedPatternChecks {
							penaltyPoints3 += 40
						}
					}
				}

				//penalty rule 4 part 1
				if QrArrayWithMask[i][j] == drawer.BLACK_COLOR {
					blackScuares++
				} else {
					whiteSqueares++
				}
			}
		}

		//penalty score 4 part 2
		blackPercentaje := math.Floor((float64(blackScuares) / float64(whiteSqueares)) * 100)
		mod5 := float64(int(blackPercentaje) % 5)
		firstNumber := math.Abs(blackPercentaje-mod5-50) / 5
		secondNumber := math.Abs(blackPercentaje+(5-mod5)-50) / 5
		if firstNumber < secondNumber {
			penaltyPoints4 = uint(firstNumber) * 10
		} else {
			penaltyPoints4 = uint(secondNumber) * 10
		}

		maskScore := penaltyPoints1 + penaltyPoints2 + penaltyPoints3 + penaltyPoints4
		logger.Info("-- mask: ", maskPaterns[maskIndex], " value: ", maskScore)
		if maskScore < lowestScore {
			bestMask = maskPaterns[maskIndex]
			lowestScore = maskScore
		}

	}
	return bestMask
}

func generateQR(QRVersionInfo QRCodeInfo, QRCode_final_step uint8) [][]uint8 {

	QRArrayBase := generateQRTemplate(QRVersionInfo)

	var data []bool //aca tendria que hacer un make ya que se cuando espacio tiene el qrcode en base a las words
	data = append(data, getEncodeMode_Binary(QRVersionInfo)...)

	if QRCode_final_step >= QR_CODE_STEP_CHARACTER_COUNT {
		data = append(data, getCharacterCount_Binary(QRVersionInfo)...)
	}

	if QRCode_final_step >= QR_CODE_STEP_ENCODE_DATA {
		data = append(data, getString_Encoded(QRVersionInfo)...)
		data = append(data, getPadingBits_Binary(QRVersionInfo, len(data))...)
		logger.Info("✓ Data encoded.")
	}

	var ERcodewords [][]bool

	if QRCode_final_step >= QR_CODE_STEP_ERROR_CORRECTION {
		ERcodewords = getCodeWords_Encoded(QRVersionInfo, data)
		logger.Info("✓ Created Error Correction Codewords.")
	} else {
		ERcodewords = [][]bool{}
	}

	encodedMessage := getStructuredFinalMessage(QRVersionInfo, data, ERcodewords)
	QRArrayWithData := addDataToQRCode(QRArrayBase, QRVersionInfo, encodedMessage)

	if QRCode_final_step < QR_CODE_STEP_MASK {
		return QRArrayWithData
	}

	logger.Info("✓ Added code words to QR code.")

	QRVersionInfo.MaskPatern = getBestMaskPattern(QRVersionInfo, QRArrayBase, QRArrayWithData)
	logger.Info("✓ Got best mask pattern: ", QRVersionInfo.MaskPatern)

	QRArrayWithMask := applyMask(QRVersionInfo.MaskPatern, QRVersionInfo.ErrorLevel, QRArrayBase, QRArrayWithData)
	return QRArrayWithMask
}

func main() {
	//TODO add the option to chose the level of error correction

	stringToEncode := "123456"

	imageName := "QRCode"
	saveLocation := "C:\\Users\\marce\\Documents\\Git\\QRCodeGenerator\\" + imageName + ".png"

	logger.Info("Generating QR code for data: ", stringToEncode)
	QRversion, err := getQRInfoByData(stringToEncode)
	if err != nil {
		logger.Error("Error obtaining info for QR Code, Error: ", err)
		return
	}
	logger.Info("Using Version: ", QRversion.Version, ", Size: ", QRversion.Size, ", Error Correction: ", QRversion.ErrorLevel, ", Encoding Mode: ", QRversion.EncodingMode)

	QRArray := generateQR(QRversion, QR_CODE_STEP_MASK)
	logger.Info("Finished encoding data")

	logger.Info("Generating Img")
	drawer.DrawQRCode(QRArray, QRversion, saveLocation)

	logger.Info("Finished generating QR code, saved in: ", saveLocation)
}
