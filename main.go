package main

import (
	"MostCommonColors/fileManager"
	"MostCommonColors/imageprocessor"
	"fmt"
	"image"
	"sync"
)

var wg sync.WaitGroup

func main() {
	mostCommonColors := []DecodedImage{}

	//import the set of URLs to process
	processedURLs, err := fileManager.ImportURLs("input.txt")
	if err != nil {
		fmt.Println("Error : ", err)
	}

	//download the images and decode them into image filetype files
	for u, _ := range processedURLs {
		i, err := fileManager.DownloadAndDecodeImage(u)
		if err != nil {
			fmt.Println("Error Decoding the Image :", err)
			continue
		}

		decodedImage := DecodedImage{
			url:              u,
			img:              i,
			mostCommonColors: "",
		}
		mostCommonColors = append(mostCommonColors, decodedImage)
	}
	//extract the dominant colors
	for _, img := range mostCommonColors {
		colors := imageprocessor.ExtractDominantColors(img)
		fmt.Println("Colors ", colors)
		for _, c := range colors {
			_, err = file.WriteString(",")
			r, g, b, _ := c.RGBA()
			_, err = file.WriteString(fmt.Sprintf("#%02x%02x%02x", r, g, b))
			fmt.Println("DONE! ", url, r, g, b)
		}
		_, err = file.WriteString("\n")
		if err != nil {
			errorChan <- fmt.Errorf("problem writing out to the file: %v", err)
			fmt.Println("problem writing out to the file: ", err)
			return
		}
	}

	//write the data to the output file

}

type DecodedImage struct {
	url              string
	img              *image.RGBA
	mostCommonColors string
}
