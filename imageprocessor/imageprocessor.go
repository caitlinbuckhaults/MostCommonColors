package imageprocessor

import (
	"fmt"
	"github.com/XuanMaoSecLab/kmeans"
	"image"
	"image/color"
	"net"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	mu        sync.Mutex
	transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	client = &http.Client{Transport: transport}
)

func ProcessImageSimpleish(url string) ([]color.Color, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//decode
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return extractDominantColors(img), nil
}

func ProcessImageOptimized(url string) error {

	//download the img
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//decode the img
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}

	//extract the dominant colors
	colors := extractDominantColors(img)

	//writeout
	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile("output.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(url)
	for _, c := range colors {
		_, err = file.WriteString(",")
		_, err = file.WriteString(fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B))
	}
	_, err = file.WriteString("\n")
	if err != nil {
		return err

	}
	return nil
}

//Simplest approach to extracting the dominant colors.
//Counts the number of pixels of each color in the image, and
//return the 3 most common colors.
func extractDominantColors(img image.Image) []color.Color {
	//we're hardcoded to returning the top three colors for now, but this can
	//be updated to return as many as we want instead or refactored to respond
	//to an external input if we want instead (GUI? Config File? ETC)
	var countOfColorsToReturn = 3

	//create a map for each color's count
	colorCount := make(map[color.Color]int)
	//loop through each pixel in the image and increment the count for the appropriate color
	for _, c := range image.RGBAColorModel.Convert(img).([]color.RGBA) {
		colorCount[c]++
	}

	//sort our color counts
	var colors []color.Color
	for c := range colorCount {
		colors = append(colors, c)
	}
	//I could handcraft a more efficient sorting algorithm here, however
	//the compiler under the hood is likely pattern matching this code and
	//optimizing it better than I could. I'll leave it as is until I need
	//to change it for a specific goal.
	sort.Slice(colors, func(i, j int) bool {
		return colorCount[colors[i]] > colorCount[colors[j]]
	})

	//return the top colors, capped by our return quantity choice.
	if len(colors) > countOfColorsToReturn {
		return colors[:countOfColorsToReturn]
	}
	return colors
}

func extractDominantColorsKmeans(img image.Image) []color.Color {
	//convert the image to a slice of kmeans.Vectors
	var vectors []kmeans.Vector
	for _, c := range image.RGBA.ColorModel().Convert(img).([]color.Color) {
		vectors = append(vectors, kmeans.Vector{float64(c.R), float64(c.G), floats64(c.B)})
	}

	// group the pixels into clusters
	clusters := kmeans.KMeans(vectors, 3)

	//return the centroids of the largest clusters as the dominant colors
	var colors []color.Color
	for _, c := range clusters {
		colors = append(colors, color.RGBA[uint8(c.Centroid[0]), uint8(c.Centroid[1]), uint8(c.Centroid[2]), 255})
}
return colors
	}

