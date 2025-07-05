package drawer

import (
	"QRCodeGenerator/generator"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func saveImage(image *image.RGBA, saveLocation string) {
	myfile, err := os.Create(saveLocation)
	if err != nil {
		panic(err)
	}
	defer myfile.Close()
	png.Encode(myfile, image)
}

func DrawQRCode(QRArray [][]uint8, QRversion generator.QRCodeInfo, locationToSave string) {
	cellSize := 10
	quietArea := 100
	imageSize := (QRversion.Size * cellSize) + quietArea

	backgroundColor := color.RGBA{255, 255, 255, 255} // white
	QRImage := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))
	draw.Draw(QRImage, QRImage.Bounds(), &image.Uniform{backgroundColor}, image.ZP, draw.Src)

	for i := range QRArray {
		for jPosition, j := range QRArray[i] {
			var colorCell color.RGBA

			switch j {
			case BLACK_COLOR:
				colorCell = color.RGBA{0, 0, 0, 255}
			case WHITE_COLOR:
				colorCell = color.RGBA{255, 255, 255, 255}
			case GREEN_COLOR:
				colorCell = color.RGBA{0, 255, 0, 255}
			case BLUE_COLOR:
				colorCell = color.RGBA{0, 0, 255, 255}
			default:
				colorCell = color.RGBA{255, 0, 0, 255} // to debug basically
			}

			cell := image.Rect(quietArea/2+cellSize*jPosition, quietArea/2+cellSize*i, quietArea/2+cellSize*(jPosition+1), quietArea/2+cellSize*(i+1))
			draw.Draw(QRImage, cell, &image.Uniform{colorCell}, image.ZP, draw.Src)
		}
	}
	saveImage(QRImage, locationToSave)
}
