package ssd1306fileprep

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
)

// OpenImage opens a image
func OpenImage(filename string) (loadedImage image.Image) {
	// Read image from file that already exists
	existingImageFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	_, imageType, err := image.Decode(existingImageFile)
	if err != nil {
		log.Fatal(err)
		log.Println(imageType)
	}

	existingImageFile.Seek(0, 0)

	switch imageType {
	case "png":
		loadedImage, err = png.Decode(existingImageFile)
	case "jpeg":
		loadedImage, err = jpeg.Decode(existingImageFile)
	default:
		err = errors.New("FileType Not Detected")
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	return loadedImage
}

// ConvertBW converts a image to black and white
func ConvertBW(imgToConvert image.Image, threshold uint8) (convertedImage *image.Gray) {
	w, h := imgToConvert.Bounds().Max.X, imgToConvert.Bounds().Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	grayWhite := color.Gray{uint8(255)}
	grayBlack := color.Gray{uint8(0)}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			imageColor := imgToConvert.At(x, y)
			rr, gg, bb, _ := imageColor.RGBA()
			r := math.Pow(float64(rr), 2.2)
			g := math.Pow(float64(gg), 2.2)
			b := math.Pow(float64(bb), 2.2)
			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
			Y := uint16(m + 0.5)
			if uint8(Y>>8) >= threshold {
				grayScale.Set(x, y, grayWhite)
			} else {
				grayScale.Set(x, y, grayBlack)
			}

		}
	}
	return grayScale
}

// ToBWByteSlice converts a gray image to a black and white byte slice structured like the ssd1306 graphics ram
func ToBWByteSlice(imgToConvert *image.Gray, threshold uint8) (outputBytes [][]byte) {
	w, h := imgToConvert.Bounds().Max.X, imgToConvert.Bounds().Max.Y
	outputBytes = make([][]byte, h/8, h/8)
	for i := 0; i < len(outputBytes); i++ {
		outputBytes[i] = make([]byte, w)
	}
	for currentY := 0; currentY < h/8; currentY++ {
		var CurrentPage []byte = make([]byte, w, w)
		for currentX := 0; currentX < w; currentX++ {
			var constructedVLine byte = 0
			if imgToConvert.GrayAt(currentX, 0+(8*currentY)).Y >= threshold {
				constructedVLine++
			}
			if imgToConvert.GrayAt(currentX, 1+(8*currentY)).Y >= threshold {
				constructedVLine += 2
			}
			if imgToConvert.GrayAt(currentX, 2+(8*currentY)).Y >= threshold {
				constructedVLine += 4
			}
			if imgToConvert.GrayAt(currentX, 3+(8*currentY)).Y >= threshold {
				constructedVLine += 8
			}
			if imgToConvert.GrayAt(currentX, 4+(8*currentY)).Y >= threshold {
				constructedVLine += 16
			}
			if imgToConvert.GrayAt(currentX, 5+(8*currentY)).Y >= threshold {
				constructedVLine += 32
			}
			if imgToConvert.GrayAt(currentX, 6+(8*currentY)).Y >= threshold {
				constructedVLine += 64
			}
			if imgToConvert.GrayAt(currentX, 7+(8*currentY)).Y >= threshold {
				constructedVLine += 128
			}
			CurrentPage[currentX] = constructedVLine
		}
		outputBytes[currentY] = CurrentPage
	}
	return outputBytes
}

// WriteImage wires a image to a file
func WriteImage(imgToWrite image.Image, fileName string) {
	newfile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer newfile.Close()
	png.Encode(newfile, imgToWrite)
}

// WriteBWByte Writes a byteslice to file
func WriteBWByte(ByteSlice [][]byte, fileName string) {
	newfile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer newfile.Close()

	for i := 0; i < len(ByteSlice); i++ {
		newfile.Write(ByteSlice[i])
		// newfile.WriteString("\n")
		newfile.Sync()
	}
}
