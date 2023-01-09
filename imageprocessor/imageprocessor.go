package imageprocessor

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"math"
	"math/rand"
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

func ProcessImageSimple(url string) ([]color.Color, error) {
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
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Problem downloading the image at ", url, " : ", err)
		return err
	}
	defer resp.Body.Close()

	//decode the img
	img, err := decodeJpeg(resp.Body)
	if err != nil {
		fmt.Println("Problem decoding the image at ", url, " : ", err)
		return err
	}

	//extract the dominant colors
	colors := extractDominantColorsKmeans(img)
	fmt.Println("Dominant colors successfully extracted")

	//writeout
	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile("output.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Can't Open output.csv: ", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(url)
	for _, c := range colors {
		_, err = file.WriteString(",")
		r, g, b, _ := c.RGBA()
		fmt.Println("Color: ", r, g, b)
		_, err = file.WriteString(fmt.Sprintf("#%02x%02x%02x", r, g, b))
	}
	_, err = file.WriteString("\n")
	if err != nil {
		fmt.Println("Problem writing out to the file: ", err)
		return err

	}
	return nil
}

//Simplest approach to extracting the dominant colors.
//Counts the number of pixels of each color in the image, and
//return the 3 most common colors.
func extractDominantColors(img image.Image) []color.Color {

	//create a map for each color's count
	colorCount := make(map[color.Color]int)
	//loop through each pixel in the image and increment the count for the appropriate color
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			colorCount[img.At(x, y)]++
		}
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

	// Return the 3 most prevalent colors
	return []color.Color{colors[0], colors[1], colors[2]}
}

func extractDominantColorsKmeans(img image.Image) []color.Color {
	//initialize the centroids with random pixels
	pixels := getPixels(img)
	centroids := []color.Color{pixels[rand.Intn(len(pixels))], pixels[rand.Intn(len(pixels))], pixels[rand.Intn(len(pixels))]}

	//iterate until the centroids converge
	for {
		//assign each pixel to its nearest centroid
		clusters := [][]color.Color{{}, {}, {}}

		for _, p := range pixels {
			minimumDistance := math.MaxFloat64
			minimumIndex := 0
			for i, c := range centroids {
				distance := distance(p, c)
				if distance < minimumDistance {
					minimumDistance = distance
					minimumIndex = i
				}
			}
			clusters[minimumIndex] = append(clusters[minimumIndex], p)
		}
		// calculate the new centroids
		newCentroids := []color.Color{averageColor(clusters[0]), averageColor(clusters[1]), averageColor(clusters[2])}
		//have the centroids converged yet?
		if centroidsEqual(centroids, newCentroids) {
			return centroids
		}
		//if not, update the centroids
		centroids = newCentroids
	}
}

//return a slice of color.RGBA values for the pixels in a given image
func getPixels(img image.Image) []color.Color {
	bounds := img.Bounds()
	var pixels []color.Color
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixels = append(pixels, img.At(x, y))
		}
	}
	return pixels
}

//get the euclidean distance between two colors
func distance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return math.Sqrt(float64((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2)))
}

//return the average of a slice of colors.
func averageColor(colors []color.Color) color.Color {
	var r, g, b float64
	for _, c := range colors {
		r1, g1, b1, _ := c.RGBA()
		r += float64(r1)
		g += float64(g1)
		b += float64(b1)
	}
	n := float64(len(colors))
	return color.RGBA{uint8(r / n), uint8(g / n), uint8(b / n), 255}
}

//are the two centroids slices the same?
func centroidsEqual(c1, c2 []color.Color) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i := range c1 {
		if c1[i] != c2[i] {
			return false
		}
	}
	return true
}

func decodeJpeg(reader io.Reader) (image.Image, error) {

	img, err := jpeg.Decode(reader)
	if err != nil {
		fmt.Println("Problem jpeg decoding the file: ", err)
		return nil, err
	}
	return img, nil
}
