package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/exp/shiny/widget"
)

func main() {
	file := flag.String("filename", "", "filename to open. (Required)")
	outFile := flag.String("outFile", "", "If you want to save the image insted of viewing it set this option must have a .png extention")
	displaySizeX := flag.Int("displaySizeX", 128, "sets the with of the screen (defaults to 128)")
	displaySizeY := flag.Int("displaySizeY", 64, "sets the Height of the screen (defaults to 64)")
	flag.Parse()

	if *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	f, err := os.Open(*file)
	if err != nil {
		os.Stderr.WriteString("File not found\n")
		os.Exit(-1)
	}
	defer f.Close()
	byteImage, err := openimage(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	decodedImage := parseImage(byteImage, *displaySizeX, *displaySizeY)

	if *outFile != "" {
		f, err = os.Create(*outFile)
		defer f.Close()
		err = png.Encode(f, decodedImage)
	} else {
		driver.Main(func(s screen.Screen) {
			image := widget.NewImage(decodedImage, decodedImage.Bounds())
			w := widget.NewSheet(image)
			if err := widget.RunWindow(s, w, &widget.RunWindowOptions{
				NewWindowOptions: screen.NewWindowOptions{
					Width:  *displaySizeX * 2,
					Height: *displaySizeY,
					Title:  "SSD .b Image Viewr",
				},
			}); err != nil {
				log.Fatal(err)
			}
		})
	}
}

func openimage(file *os.File) ([]byte, error) {
	filestats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	len := filestats.Size()
	bytesIn := make([]byte, len)
	_, err = file.Read(bytesIn)
	if err != nil {
		return nil, err
	}
	return bytesIn, nil
}

func parseImage(byteImage []byte, XSize int, YSize int) *image.Gray {
	picture := image.NewGray(image.Rect(0, 0, XSize, YSize+1))
	for currentY := 0; currentY < picture.Bounds().Max.X; currentY++ {
		for currentX := 0; currentX < picture.Bounds().Max.Y; currentX++ {
			picture.Set(currentY, currentX, color.Black)
		}
	}
	slicePosition := ((YSize * XSize) / 8) - 1

	for currentY := ((YSize / 8) - 1); currentY >= 0; currentY-- {
		for currentX := XSize - 1; currentX >= 0; currentX-- {
			currentVByte := uint8(byteImage[slicePosition])
			if currentVByte >= 128 {
				picture.Set(currentX, 7+(currentY*8), image.White)
				currentVByte = currentVByte - 128
			}

			if currentVByte >= 64 {
				picture.Set(currentX, 6+(currentY*8), image.White)
				currentVByte = currentVByte - 64
			}

			if currentVByte >= 32 {
				picture.Set(currentX, 5+(currentY*8), image.White)
				currentVByte = currentVByte - 32
			}

			if currentVByte >= 16 {
				picture.Set(currentX, 4+(currentY*8), image.White)
				currentVByte = currentVByte - 16
			}

			if currentVByte >= 8 {
				picture.Set(currentX, 3+(currentY*8), image.White)
				currentVByte = currentVByte - 8
			}

			if currentVByte >= 4 {
				picture.Set(currentX, 2+(currentY*8), image.White)
				currentVByte = currentVByte - 4
			}

			if currentVByte >= 2 {
				picture.Set(currentX, 1+(currentY*8), image.White)
				currentVByte = currentVByte - 2
			}

			if currentVByte == 1 {
				picture.Set(currentX, 0+(currentY*8), image.White)
			}

			slicePosition--
		}
	}
	return picture
}
