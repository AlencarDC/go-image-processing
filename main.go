package photochopp

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello")
	reader, err := os.Open("images/Gramado_22k.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	fmt.Println("Bounds: ", img.Bounds())
	fmt.Println("Width: ", width, "Height: ", height)

	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	fmt.Println(rgba.Stride, "len:", len(rgba.Pix))

	customImg := NewImage(rgba)
	cq := ColorQuantization{NumberOfDesiredColors: 8}
	cq.Apply(customImg)
	// Flip vertically
	// lineBuff := make([]uint8, 4*width)
	// for y := 0; y < height/2; y++ {
	// 	// fmt.Printf("Upperline -> coord: (%v, %v) até (%v, %v) offsets [%v : %v]\n", 0, y, width, y, rgba.PixOffset(0, y), rgba.PixOffset(width, y)-1)
	// 	// fmt.Printf("Lowerline -> coord: (%v, %v) até (%v, %v) offsets [%v : %v]\n", 0, height-y-1, width, height-y-1, rgba.PixOffset(0, height-y-1), rgba.PixOffset(width, height-y-1)-1)
	// 	upperLine := rgba.Pix[rgba.PixOffset(0, y) : rgba.PixOffset(width, y)-1]
	// 	lowerLine := rgba.Pix[rgba.PixOffset(0, height-y-1) : rgba.PixOffset(width, height-y-1)-1]
	// 	copy(lineBuff, upperLine)
	// 	copy(upperLine, lowerLine)
	// 	copy(lowerLine, lineBuff)
	// }

	// Flip horizontally
	// pixelBuff := make([]uint8, 4)
	// for x := 0; x < width/2; x++ {
	// 	for y := 0; y < height; y++ {
	// 		leftPixel := rgba.Pix[rgba.PixOffset(x, y) : rgba.PixOffset(x, y)+4]
	// 		rightPixel := rgba.Pix[rgba.PixOffset(width-x-1, y) : rgba.PixOffset(width-x-1, y)+4]
	// 		copy(pixelBuff, leftPixel)
	// 		copy(leftPixel, rightPixel)
	// 		copy(rightPixel, pixelBuff)
	// 	}
	// }

	// Luminance (RGBA to Gray Scale)
	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		offset := rgba.PixOffset(x, y)
	// 		r := &rgba.Pix[offset+0]
	// 		g := &rgba.Pix[offset+1]
	// 		b := &rgba.Pix[offset+2]

	// 		luminance := 0.299*float64(*r) + 0.587*float64(*g) + 0.114*float64(*b)

	// 		*r = uint8(luminance)
	// 		*g = uint8(luminance)
	// 		*b = uint8(luminance)
	// 	}
	// }

	// Histogram RGB
	// var histogramR [256]uint8
	// var histogramG [256]uint8
	// var histogramB [256]uint8
	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		offset := rgba.PixOffset(x, y)
	// 		histogramR[rgba.Pix[offset+0]] += 1
	// 		histogramG[rgba.Pix[offset+1]] += 1
	// 		histogramB[rgba.Pix[offset+2]] += 1
	// 	}
	// }

	// minColor := 0
	// maxColor := 255
	// intensitySize := maxColor - minColor + 1
	// desiredColorSize := 8

	// binSize := intensitySize / desiredColorSize
	// fmt.Println("Bin size:", binSize, minColor, maxColor, intensitySize, desiredColorSize)
	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		offset := rgba.PixOffset(x, y)
	// 		r := &rgba.Pix[offset+0]
	// 		g := &rgba.Pix[offset+1]
	// 		b := &rgba.Pix[offset+2]

	// 		newBinBegin := (int(*r) / binSize) * binSize
	// 		newColor := uint8(newBinBegin + binSize/2)
	// 		*r = newColor

	// 		newBinBegin = (int(*g) / binSize) * binSize
	// 		newColor = uint8(newBinBegin + binSize/2)
	// 		*g = newColor

	// 		newBinBegin = (int(*b) / binSize) * binSize
	// 		newColor = uint8(newBinBegin + binSize/2)
	// 		*b = newColor
	// 	}
	// }

	// fmt.Println(histogramR)
	// fmt.Println(histogramG)
	// fmt.Println(histogramB)

	outputImage, err := os.Create("./outputs/output.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputImage.Close()

	err = jpeg.Encode(outputImage, rgba, nil)
	if err != nil {
		log.Fatal(err)
	}
}
