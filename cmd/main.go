package main

import (
	"github.com/nfnt/resize"
	"github.com/jax-b/ssd1306FilePrep"
)

func main() {
	inFile := flag.String("inputFile", "", "input image file. (Required)")
	outFile := flag.String("outputFile", "", "output .b file. (Required)")
	displaySizeX := flag.Int("displaySizeX", 128, "sets the with of the screen (defaults to 128)")
	displaySizeY := flag.Int("displaySizeX", 64, "sets the Height of the screen (defaults to 64)")
	threshold := flag.Int("threshold", 128, "sets the threshold value for the conversion (defaults to 128, from 0-255)")
	flag.Parse()

	if *inFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *outFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	image := ssd1306FilePrep.openImage(*inFile)
	resizedImage := resize.Resize(uint(*displaySizeX), uint(*displaySizeY), image, resize.NearestNeighbor)
	test := ssd1306FilePrep.toBWByteSlice(convertBW(resizedImage, uint(*threshold), uint(*threshold))
	ssd1306FilePrep.writeBWByte(test, *outFile)
}